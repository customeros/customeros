package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"golang.org/x/net/context"
	"time"
)

type TicketRepository interface {
	GetMatchedTicketId(ctx context.Context, tenant string, ticket entity.TicketData) (string, error)
	MergeTicket(ctx context.Context, tenant string, syncDate time.Time, user entity.TicketData) error
	MergeTagForTicket(ctx context.Context, tenant, ticketId, tagName, externalSystem string) error
	LinkTicketWithCollaboratorUserByExternalId(ctx context.Context, tenant, ticketId, userExternalId, externalSystem string) error
	LinkTicketWithFollowerUserByExternalId(ctx context.Context, tenant, ticketId, userExternalId, externalSystem string) error
	LinkTicketWithSubmitterUserOrContactByExternalId(ctx context.Context, tenant, ticketId, externalId, externalSystem string) error
	LinkTicketWithRequesterUserOrContactByExternalId(ctx context.Context, tenant, ticketId, externalId, externalSystem string) error
	LinkTicketWithAssigneeUserByExternalId(ctx context.Context, tenant, ticketId, externalId, externalSystem string) error
}

type ticketRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTicketRepository(driver *neo4j.DriverWithContext) TicketRepository {
	return &ticketRepository{
		driver: driver,
	}
}

func (r *ticketRepository) GetMatchedTicketId(ctx context.Context, tenant string, user entity.TicketData) (string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:TICKET_BELONGS_TO_TENANT]-(tt:Ticket)-[:IS_LINKED_WITH {externalId:$ticketExternalId}]->(e)
				WITH tt WHERE tt is not null
				return tt.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":           tenant,
				"externalSystem":   user.ExternalSystem,
				"ticketExternalId": user.ExternalId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	ticketIDs := dbRecords.([]*db.Record)
	if len(ticketIDs) > 0 {
		return ticketIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *ticketRepository) MergeTicket(ctx context.Context, tenant string, syncDate time.Time, ticket entity.TicketData) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem}) " +
		" MERGE (tt:Ticket {id:$ticketId})-[:TICKET_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET tt.createdAt=$createdAt, " +
		"				tt.updatedAt=$updatedAt, " +
		"               tt.source=$source, " +
		"				tt.sourceOfTruth=$sourceOfTruth, " +
		"				tt.appSource=$appSource, " +
		"				tt.subject=$subject, " +
		"				tt.status=$status, " +
		"				tt.priority=$priority, " +
		"				tt.description=$description, " +
		"               tt:%s" +
		" ON MATCH SET " +
		"				tt.subject = CASE WHEN tt.sourceOfTruth=$sourceOfTruth OR tt.subject is null OR tt.subject = '' THEN $subject ELSE tt.subject END, " +
		"				tt.description = CASE WHEN tt.sourceOfTruth=$sourceOfTruth OR tt.description is null OR tt.description = '' THEN $description ELSE tt.description END, " +
		"				tt.status = CASE WHEN tt.sourceOfTruth=$sourceOfTruth OR tt.status is null OR tt.status = '' THEN $status ELSE tt.status END, " +
		"				tt.priority = CASE WHEN tt.sourceOfTruth=$sourceOfTruth OR tt.priority is null OR tt.priority = '' THEN $priority ELSE tt.priority END, " +
		"				tt.updatedAt=$now " +
		" WITH tt, ext " +
		" MERGE (tt)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate, r.externalUrl=$externalUrl " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" WITH tt " +
		" FOREACH (x in CASE WHEN tt.sourceOfTruth <> $sourceOfTruth THEN [tt] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateTicket {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.subject=$subject, alt.status=$status, alt.priority=$priority, alt.description=$description" +
		") RETURN tt.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Ticket_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"ticketId":       ticket.Id,
				"externalSystem": ticket.ExternalSystem,
				"externalId":     ticket.ExternalId,
				"externalUrl":    ticket.ExternalUrl,
				"syncDate":       syncDate,
				"createdAt":      ticket.CreatedAt,
				"updatedAt":      ticket.UpdatedAt,
				"source":         ticket.ExternalSystem,
				"sourceOfTruth":  ticket.ExternalSystem,
				"appSource":      ticket.ExternalSystem,
				"subject":        ticket.Subject,
				"description":    ticket.Description,
				"status":         ticket.Status,
				"priority":       ticket.Priority,
				"now":            time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *ticketRepository) MergeTagForTicket(ctx context.Context, tenant, ticketId, tagName, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant})<-[:TICKET_BELONGS_TO_TENANT]-(tt:Ticket {id:$ticketId}) " +
		" MERGE (tag:Tag {name:$tagName})-[:TAG_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET tag.id=randomUUID(), " +
		"				tag.createdAt=$now, " +
		"				tag.updatedAt=$now, " +
		"				tag.source=$source," +
		"				tag.appSource=$source," +
		"				tag:%s  " +
		" WITH DISTINCT tt, tag " +
		" MERGE (tt)-[r:TAGGED]->(tag) " +
		"	ON CREATE SET r.taggedAt=$now " +
		" return r"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Tag_"+tenant),
			map[string]interface{}{
				"tenant":   tenant,
				"ticketId": ticketId,
				"tagName":  tagName,
				"source":   externalSystem,
				"now":      time.Now().UTC(),
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *ticketRepository) LinkTicketWithCollaboratorUserByExternalId(ctx context.Context, tenant, ticketId, userExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (tt:Ticket {id:$ticketId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:FOLLOWS]->(tt)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"ticketId":       ticketId,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *ticketRepository) LinkTicketWithFollowerUserByExternalId(ctx context.Context, tenant, ticketId, userExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (tt:Ticket {id:$ticketId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:FOLLOWS]->(tt)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"ticketId":       ticketId,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *ticketRepository) LinkTicketWithSubmitterUserOrContactByExternalId(ctx context.Context, tenant, ticketId, externalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$externalId}]-(n)
					WHERE (n:User OR n:Contact)
				MATCH (tt:Ticket {id:$ticketId})-[:IS_LINKED_WITH]->(e)
				MERGE (n)-[:SUBMITTED]->(tt)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"ticketId":       ticketId,
				"externalId":     externalId,
			})
		return nil, err
	})
	return err
}

func (r *ticketRepository) LinkTicketWithRequesterUserOrContactByExternalId(ctx context.Context, tenant, ticketId, externalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$externalId}]-(n)
					WHERE (n:User OR n:Contact)
				MATCH (tt:Ticket {id:$ticketId})-[:IS_LINKED_WITH]->(e)
				MERGE (n)-[:REQUESTED]->(tt)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"ticketId":       ticketId,
				"externalId":     externalId,
			})
		return nil, err
	})
	return err
}

func (r *ticketRepository) LinkTicketWithAssigneeUserByExternalId(ctx context.Context, tenant, ticketId, externalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$externalId}]-(n:User)
				MATCH (tt:Ticket {id:$ticketId})-[:IS_LINKED_WITH]->(e)
				MERGE (n)-[:IS_ASSIGNED_TO]->(tt)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"ticketId":       ticketId,
				"externalId":     externalId,
			})
		return nil, err
	})
	return err
}
