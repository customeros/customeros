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
	GetMatchedTicketId(ctx context.Context, tenant string, user entity.TicketData) (string, error)
	MergeTicket(ctx context.Context, tenant string, syncDate time.Time, user entity.TicketData) error
	LinkTicketWithCollaboratorUserByExternalId(ctx context.Context, tenant, ticketId, userExternalId, externalSystem string) error

	NoteLinkWithContactByExternalId(ctx context.Context, tenant, noteId, contactExternalId, externalSystem string) error
	NoteLinkWithOrganizationByExternalId(ctx context.Context, tenant, noteId, organizationExternalId, externalSystem string) error
	NoteLinkWithUserByExternalOwnerId(ctx context.Context, tenant, noteId, userExternalOwnerId, externalSystem string) error
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
		"               tt:%s" +
		" ON MATCH SET " +
		"				tt.subject = CASE WHEN tt.sourceOfTruth=$sourceOfTruth OR tt.subject is null OR tt.subject = '' THEN $subject ELSE tt.subject END, " +
		"				tt.updatedAt=$now " +
		" WITH tt, ext " +
		" MERGE (tt)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate, r.externalUrl=$externalUrl " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" WITH tt " +
		" FOREACH (x in CASE WHEN tt.sourceOfTruth <> $sourceOfTruth THEN [tt] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateTicket {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource " +
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

func (r *ticketRepository) LinkTicketWithCollaboratorUserByExternalId(ctx context.Context, tenant, ticketId, userExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (tt:Ticket {id:$ticketId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:COLLABORATES_ON]->(tt)
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

func (r *ticketRepository) NoteLinkWithContactByExternalId(ctx context.Context, tenant, noteId, contactExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$contactExternalId}]-(c:Contact)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (c)-[:NOTED]->(n)
				`,
			map[string]interface{}{
				"tenant":            tenant,
				"externalSystem":    externalSystem,
				"noteId":            noteId,
				"contactExternalId": contactExternalId,
			})
		return nil, err
	})
	return err
}

func (r *ticketRepository) NoteLinkWithOrganizationByExternalId(ctx context.Context, tenant, noteId, organizationExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$organizationExternalId}]-(org:Organization)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (org)-[:NOTED]->(n)
				`,
			map[string]interface{}{
				"tenant":                 tenant,
				"externalSystem":         externalSystem,
				"noteId":                 noteId,
				"organizationExternalId": organizationExternalId,
			})
		return nil, err
	})
	return err
}

func (r *ticketRepository) NoteLinkWithUserByExternalOwnerId(ctx context.Context, tenant, noteId, userExternalOwnerId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalOwnerId:$userExternalOwnerId}]-(u:User)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:CREATED]->(n)
				`,
			map[string]interface{}{
				"tenant":              tenant,
				"externalSystem":      externalSystem,
				"noteId":              noteId,
				"userExternalOwnerId": userExternalOwnerId,
			})
		return nil, err
	})
	return err
}
