package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"time"
)

type NoteDbNodeWithParentId struct {
	Node     *dbtype.Node
	ParentId string
}

type NoteDbNodesWithTotalCount struct {
	Nodes []*NoteDbNodeWithParentId
	Count int64
}

type NoteRepository interface {
	GetPaginatedNotesForContact(session neo4j.Session, tenant, contactId string, skip, limit int) (*NoteDbNodesWithTotalCount, error)
	CreateNoteForContact(session neo4j.Session, tenant, contactId string, entity entity.NoteEntity) (*dbtype.Node, error)
	CreateNoteForOrganization(session neo4j.Session, tenant, organization string, entity entity.NoteEntity) (*dbtype.Node, error)
	UpdateNote(session neo4j.Session, tenant string, entity entity.NoteEntity) (*dbtype.Node, error)
	Delete(session neo4j.Session, tenant, noteId string) error
}

type noteRepository struct {
	driver *neo4j.Driver
}

func NewNoteRepository(driver *neo4j.Driver) NoteRepository {
	return &noteRepository{
		driver: driver,
	}
}

func (r *noteRepository) GetPaginatedNotesForContact(session neo4j.Session, tenant, contactId string, skip, limit int) (*NoteDbNodesWithTotalCount, error) {
	result := new(NoteDbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
											(c)-[:NOTED]->(n:Note)
											RETURN count(n) as count`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		result.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(fmt.Sprintf(
			"MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
				" (c)-[:NOTED]->(n:Note)"+
				" RETURN n, c.id "+
				" SKIP $skip LIMIT $limit"),
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"skip":      skip,
				"limit":     limit,
			})
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		noteDBNodeWithParentId := new(NoteDbNodeWithParentId)
		noteDBNodeWithParentId.Node = utils.NodePtr(v.Values[0].(neo4j.Node))
		noteDBNodeWithParentId.ParentId = v.Values[1].(string)
		result.Nodes = append(result.Nodes, noteDBNodeWithParentId)
	}
	return result, nil
}

func (r *noteRepository) UpdateNote(session neo4j.Session, tenant string, entity entity.NoteEntity) (*dbtype.Node, error) {
	query := "MATCH (n:%s {id:$noteId}) " +
		" SET 	n.html=$html, " +
		"		n.sourceOfTruth=$sourceOfTruth, " +
		"		n.updatedAt=$now " +
		" RETURN n"
	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(fmt.Sprintf(query, "Note_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"noteId":        entity.Id,
				"html":          entity.Html,
				"sourceOfTruth": entity.SourceOfTruth,
				"now":           time.Now().UTC(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(txResult, err)
	})
	if err != nil {
		return nil, err
	}
	return queryResult.(*dbtype.Node), nil
}

func (r *noteRepository) CreateNoteForContact(session neo4j.Session, tenant, contactId string, entity entity.NoteEntity) (*dbtype.Node, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (c)-[:NOTED]->(n:Note {id:randomUUID()}) " +
		" ON CREATE SET n.html=$html, " +
		"				n.createdAt=$createdAt, " +
		"				n.updatedAt=$createdAt, " +
		"				n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"				n:%s " +
		" RETURN n"

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "Note_"+tenant),
			map[string]any{
				"tenant":        tenant,
				"contactId":     contactId,
				"html":          entity.Html,
				"createdAt":     time.Now().UTC(),
				"source":        entity.Source,
				"sourceOfTruth": entity.SourceOfTruth,
				"appSource":     entity.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *noteRepository) CreateNoteForOrganization(session neo4j.Session, tenant, organizationId string, entity entity.NoteEntity) (*dbtype.Node, error) {
	query := "MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (org)-[:NOTED]->(n:Note {id:randomUUID()}) " +
		" ON CREATE SET n.html=$html, " +
		"				n.createdAt=$createdAt, " +
		"				n.updatedAt=$createdAt, " +
		"				n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"				n:%s " +
		" RETURN n"

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "Note_"+tenant),
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
				"html":           entity.Html,
				"createdAt":      time.Now().UTC(),
				"source":         entity.Source,
				"sourceOfTruth":  entity.SourceOfTruth,
				"appSource":      entity.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *noteRepository) Delete(session neo4j.Session, tenant, noteId string) error {
	query := "MATCH (n:%s {id:$noteId}) DETACH DELETE n"
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, "Note_"+tenant),
			map[string]interface{}{
				"tenant": tenant,
				"noteId": noteId,
			})
		return nil, err
	})
	return err
}
