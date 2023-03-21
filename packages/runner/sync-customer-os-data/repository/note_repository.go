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

type NoteRepository interface {
	GetMatchedNoteId(ctx context.Context, tenant string, note entity.NoteData) (string, error)
	MergeNote(ctx context.Context, tenant string, syncDate time.Time, note entity.NoteData) error
	NoteLinkWithContactByExternalId(ctx context.Context, tenant, noteId, contactExternalId, externalSystem string) error
	NoteLinkWithOrganizationByExternalId(ctx context.Context, tenant, noteId, organizationExternalId, externalSystem string) error
	NoteLinkWithTicketByExternalId(ctx context.Context, tenant, noteId, organizationExternalId, externalSystem string) error
	NoteLinkWithCreatorUserByExternalId(ctx context.Context, tenant, noteId, userExternalId, externalSystem string) error
	NoteLinkWithCreatorUserByExternalOwnerId(ctx context.Context, tenant, noteId, userExternalOwnerId, externalSystem string) error
	NoteLinkWithCreatorUserOrContactByExternalId(ctx context.Context, tenant, noteId, creatorExternalId, externalSystem string) error
}

type noteRepository struct {
	driver *neo4j.DriverWithContext
}

func NewNoteRepository(driver *neo4j.DriverWithContext) NoteRepository {
	return &noteRepository{
		driver: driver,
	}
}

func (r *noteRepository) GetMatchedNoteId(ctx context.Context, tenant string, note entity.NoteData) (string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (e)<-[:IS_LINKED_WITH {externalId:$noteExternalId}]-(n:Note)
				WITH n WHERE n is not null
				return n.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": note.ExternalSystem,
				"noteExternalId": note.ExternalId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	noteIDs := dbRecords.([]*db.Record)
	if len(noteIDs) > 0 {
		return noteIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *noteRepository) MergeNote(ctx context.Context, tenant string, syncDate time.Time, note entity.NoteData) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new Note if it does not exist
	// If Note exists, and sourceOfTruth is acceptable then update Note.
	//   otherwise create/update AlternateNote for incoming source, with a new relationship 'ALTERNATE'
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem}) " +
		" MERGE (n:Note {id:$noteId}) " +
		" ON CREATE SET " +
		"				n.createdAt=$createdAt, " +
		"				n.updatedAt=$createdAt, " +
		"              	n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"              	n.html=$html, " +
		"              	n.text=$text, " +
		"				n:%s " +
		" ON MATCH SET 	n.html = CASE WHEN n.sourceOfTruth=$sourceOfTruth OR n.html is null or n.html = '' THEN $html ELSE n.html END, " +
		"             	n.text = CASE WHEN n.sourceOfTruth=$sourceOfTruth OR n.text is null or n.text = '' THEN $text ELSE n.text END, " +
		"				n.updatedAt = $now " +
		" WITH n, ext " +
		" MERGE (n)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.syncDate=$syncDate " +
		" ON MATCH SET r.syncDate=$syncDate " +
		" WITH n " +
		" FOREACH (x in CASE WHEN n.sourceOfTruth <> $sourceOfTruth THEN [n] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateNote {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.html=$html, alt.text=$text " +
		" ) " +
		" RETURN n.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Note_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"noteId":         note.Id,
				"source":         note.ExternalSystem,
				"sourceOfTruth":  note.ExternalSystem,
				"appSource":      note.ExternalSystem,
				"externalSystem": note.ExternalSystem,
				"externalId":     note.ExternalId,
				"syncDate":       syncDate,
				"html":           note.Html,
				"text":           note.Text,
				"createdAt":      note.CreatedAt,
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

func (r *noteRepository) NoteLinkWithContactByExternalId(ctx context.Context, tenant, noteId, contactExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$contactExternalId}]-(c:Contact)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (c)-[:NOTED]->(n)
				SET n:TimelineEvent, n:TimelineEvent_%s`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
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

func (r *noteRepository) NoteLinkWithOrganizationByExternalId(ctx context.Context, tenant, noteId, organizationExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$organizationExternalId}]-(org:Organization)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (org)-[:NOTED]->(n)
				SET n:TimelineEvent, n:TimelineEvent_%s`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
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

func (r *noteRepository) NoteLinkWithTicketByExternalId(ctx context.Context, tenant, noteId, ticketExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$ticketExternalId}]-(tt:Ticket)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (tt)-[:NOTED]->(n)
				`,
			map[string]interface{}{
				"tenant":           tenant,
				"externalSystem":   externalSystem,
				"noteId":           noteId,
				"ticketExternalId": ticketExternalId,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) NoteLinkWithCreatorUserByExternalId(ctx context.Context, tenant, noteId, userExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$userExternalId}]-(u:User)
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (u)-[:CREATED]->(n)
				`,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": externalSystem,
				"noteId":         noteId,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) NoteLinkWithCreatorUserByExternalOwnerId(ctx context.Context, tenant, noteId, userExternalOwnerId, externalSystem string) error {
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

func (r *noteRepository) NoteLinkWithCreatorUserOrContactByExternalId(ctx context.Context, tenant, noteId, creatorExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
				MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})<-[:IS_LINKED_WITH {externalId:$creatorExternalId}]-(c)
				WHERE c:User OR c:Contact
				MATCH (n:Note {id:$noteId})-[:IS_LINKED_WITH]->(e)
				MERGE (c)-[:CREATED]->(n)
				`,
			map[string]interface{}{
				"tenant":            tenant,
				"externalSystem":    externalSystem,
				"noteId":            noteId,
				"creatorExternalId": creatorExternalId,
			})
		return nil, err
	})
	return err
}
