package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type NoteRepository interface {
	MergeNote(tenant string, syncDate time.Time, note entity.NoteData) (string, error)
	NoteLinkWithContactByExternalId(tenant, noteId, contactExternalId, externalSystem string) error
	NoteLinkWithUserByExternalId(tenant, noteId, userExternalId, externalSystem string) error
}

type noteRepository struct {
	driver *neo4j.Driver
}

func NewNoteRepository(driver *neo4j.Driver) NoteRepository {
	return &noteRepository{
		driver: driver,
	}
}

func (r *noteRepository) MergeNote(tenant string, syncDate time.Time, note entity.NoteData) (string, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MERGE (n:Note)-[r:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET n.id=randomUUID(), n.createdAt=$createdAt, " +
		"               n.source=$source, n.sourceOfTruth=$sourceOfTruth, n.appSource=$appSource, " +
		"               r.syncDate=$syncDate, n.html=$html, n:%s " +
		" ON MATCH SET r.syncDate=$syncDate, n.html=$html " +
		" RETURN n.id"

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "Note_"+tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"source":         note.ExternalSystem,
				"sourceOfTruth":  note.ExternalSystem,
				"appSource":      note.ExternalSystem,
				"externalSystem": note.ExternalSystem,
				"externalId":     note.ExternalId,
				"syncDate":       syncDate,
				"html":           note.Html,
				"createdAt":      note.CreatedAt,
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return "", err
	}
	return dbRecord.(string), nil
}

func (r *noteRepository) NoteLinkWithContactByExternalId(tenant, noteId, contactExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
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

func (r *noteRepository) NoteLinkWithUserByExternalId(tenant, noteId, userExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
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
