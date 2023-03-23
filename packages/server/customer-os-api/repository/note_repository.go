package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
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
	GetPaginatedNotesForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, skip, limit int) (*NoteDbNodesWithTotalCount, error)
	GetTimeRangeNotesForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, start, end time.Time) ([]*neo4j.Node, error)
	GetPaginatedNotesForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, skip, limit int) (*NoteDbNodesWithTotalCount, error)
	CreateNoteForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, entity entity.NoteEntity) (*dbtype.Node, error)
	CreateNoteForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organization string, entity entity.NoteEntity) (*dbtype.Node, error)

	UpdateNote(ctx context.Context, session neo4j.SessionWithContext, tenant string, entity entity.NoteEntity) (*dbtype.Node, error)

	Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, noteId string) error
	SetNoteCreator(ctx context.Context, session neo4j.SessionWithContext, tenant, userId, noteId string) error

	GetNotedEntitiesForNotes(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
}

type noteRepository struct {
	driver *neo4j.DriverWithContext
}

func NewNoteRepository(driver *neo4j.DriverWithContext) NoteRepository {
	return &noteRepository{
		driver: driver,
	}
}

func (r *noteRepository) GetPaginatedNotesForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, skip, limit int) (*NoteDbNodesWithTotalCount, error) {
	result := new(NoteDbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
											(c)-[:NOTED]->(n:Note)
											RETURN count(n) as count`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		result.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
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
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
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

func (r *noteRepository) GetTimeRangeNotesForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, start, end time.Time) ([]*neo4j.Node, error) {

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, fmt.Sprintf(
			"MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
				" (c)-[:NOTED]->(n:Note)"+
				" WHERE n.createdAt > $start AND n.createdAt < $end"+
				" RETURN n, c.id "),
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"start":     start.UTC(),
				"end":       end.UTC(),
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	result := make([]*neo4j.Node, len(dbRecords.([]*neo4j.Record)))

	for i, v := range dbRecords.([]*neo4j.Record) {
		result[i] = utils.NodePtr(v.Values[0].(neo4j.Node))
	}
	return result, nil
}

func (r *noteRepository) GetPaginatedNotesForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, skip, limit int) (*NoteDbNodesWithTotalCount, error) {
	result := new(NoteDbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
											(org)-[:NOTED]->(n:Note)
											RETURN count(n) as count`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		result.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			"MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
				" (org)-[:NOTED]->(n:Note)"+
				" RETURN n, org.id "+
				" SKIP $skip LIMIT $limit"),
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
				"skip":           skip,
				"limit":          limit,
			})
		return queryResult.Collect(ctx)
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

func (r *noteRepository) UpdateNote(ctx context.Context, session neo4j.SessionWithContext, tenant string, entity entity.NoteEntity) (*dbtype.Node, error) {
	query := "MATCH (n:%s {id:$noteId}) " +
		" SET 	n.html=$html, " +
		"		n.sourceOfTruth=$sourceOfTruth, " +
		"		n.updatedAt=$now " +
		" RETURN n"
	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		txResult, err := tx.Run(ctx, fmt.Sprintf(query, "Note_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"noteId":        entity.Id,
				"html":          entity.Html,
				"sourceOfTruth": entity.SourceOfTruth,
				"now":           utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, txResult, err)
	})
	if err != nil {
		return nil, err
	}
	return queryResult.(*dbtype.Node), nil
}

func (r *noteRepository) CreateNoteForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, entity entity.NoteEntity) (*dbtype.Node, error) {
	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (c)-[:NOTED]->(n:Note {id:randomUUID()}) " +
		" ON CREATE SET n.html=$html, " +
		"				n.createdAt=$now, " +
		"				n.updatedAt=$now, " +
		"				n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"				n:Note_%s," +
		"				n:TimelineEvent," +
		"				n:TimelineEvent_%s " +
		" RETURN n"

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]any{
				"tenant":        tenant,
				"contactId":     contactId,
				"html":          entity.Html,
				"now":           utils.Now(),
				"source":        entity.Source,
				"sourceOfTruth": entity.SourceOfTruth,
				"appSource":     entity.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *noteRepository) CreateNoteForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, entity entity.NoteEntity) (*dbtype.Node, error) {
	query := "MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (org)-[:NOTED]->(n:Note {id:randomUUID()}) " +
		" ON CREATE SET n.html=$html, " +
		"				n.createdAt=$now, " +
		"				n.updatedAt=$now, " +
		"				n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"				n:Note_%s," +
		"				n:TimelineEvent," +
		"				n:TimelineEvent_%s " +
		" RETURN n"

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
				"html":           entity.Html,
				"now":            utils.Now(),
				"source":         entity.Source,
				"sourceOfTruth":  entity.SourceOfTruth,
				"appSource":      entity.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *noteRepository) Delete(ctx context.Context, session neo4j.SessionWithContext, tenant, noteId string) error {
	query := "MATCH (n:%s {id:$noteId}) DETACH DELETE n"
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, "Note_"+tenant),
			map[string]interface{}{
				"tenant": tenant,
				"noteId": noteId,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) SetNoteCreator(ctx context.Context, session neo4j.SessionWithContext, tenant, userId, noteId string) error {
	query := "MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), " +
		" (n:Note {id:$noteId})" +
		"  MERGE (u)-[:CREATED]->(n) "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant": tenant,
				"userId": userId,
				"noteId": noteId,
			})
		return nil, err
	})
	return err
}

func (r *noteRepository) GetNotedEntitiesForNotes(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (n:Note_%s)<-[rel:NOTED]-(e) " +
		" WHERE n.id IN $ids " +
		" RETURN e, n.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"ids": ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}
