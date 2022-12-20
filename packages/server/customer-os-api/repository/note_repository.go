package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
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
