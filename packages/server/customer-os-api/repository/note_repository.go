package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	GetPaginatedNotesForContact(ctx context.Context, tenant, contactId string, skip, limit int) (*NoteDbNodesWithTotalCount, error)
	GetTimeRangeNotesForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, start, end time.Time) ([]*neo4j.Node, error)
	GetPaginatedNotesForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string, skip, limit int) (*NoteDbNodesWithTotalCount, error)
	GetNotesForMeetings(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetMentionedByNotesForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeAndId, error)

	CreateNoteForContact(ctx context.Context, tenant, contactId string, entity entity.NoteEntity) (*dbtype.Node, error)
	CreateNoteForOrganization(ctx context.Context, tenant, organization string, entity entity.NoteEntity) (*dbtype.Node, error)
	CreateNoteForMeeting(ctx context.Context, tenant, meeting string, entity *entity.NoteEntity) (*dbtype.Node, error)
	CreateNoteForMeetingTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, meeting string, entity *entity.NoteEntity) (*dbtype.Node, error)

	UpdateNote(ctx context.Context, session neo4j.SessionWithContext, tenant string, entity entity.NoteEntity) (*dbtype.Node, error)

	Delete(ctx context.Context, tenant, noteId string) error
	SetNoteCreator(ctx context.Context, tenant, userId, noteId string) error

	GetNotedEntitiesForNotes(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetMentionedEntitiesForNotes(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
}

type noteRepository struct {
	driver *neo4j.DriverWithContext
}

func NewNoteRepository(driver *neo4j.DriverWithContext) NoteRepository {
	return &noteRepository{
		driver: driver,
	}
}

func (r *noteRepository) GetPaginatedNotesForContact(ctx context.Context, tenant, contactId string, skip, limit int) (*NoteDbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.GetPaginatedNotesForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	result := new(NoteDbNodesWithTotalCount)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.GetTimeRangeNotesForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.GetPaginatedNotesForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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

func (r *noteRepository) GetNotesForMeetings(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.GetNotesForMeetings")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(`MATCH (m:Meeting_%s), 
											(m)-[:NOTED]->(n:Note_%s)  WHERE m.id IN $ids 
											RETURN n, m.id`, tenant, tenant),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}

	})

	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), nil
}

func (r *noteRepository) UpdateNote(ctx context.Context, session neo4j.SessionWithContext, tenant string, entity entity.NoteEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.UpdateNote")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (n:%s {id:$noteId}) " +
		" SET 	n.content=$content, " +
		"		n.contentType=$contentType, " +
		"		n.sourceOfTruth=$sourceOfTruth, " +
		"		n.updatedAt=$now " +
		" RETURN n"
	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		txResult, err := tx.Run(ctx, fmt.Sprintf(query, "Note_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"noteId":        entity.Id,
				"content":       entity.Content,
				"contentType":   entity.ContentType,
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

func (r *noteRepository) CreateNoteForContact(ctx context.Context, tenant, contactId string, entity entity.NoteEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.CreateNoteForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (c)-[:NOTED]->(n:Note {id:randomUUID()}) " +
		" ON CREATE SET n.content=$content, " +
		"				n.contentType=$contentType, " +
		"				n.createdAt=$now, " +
		"				n.updatedAt=$now, " +
		"				n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"				n:Note_%s," +
		"				n:TimelineEvent," +
		"				n:TimelineEvent_%s " +
		" RETURN n"

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]any{
				"tenant":        tenant,
				"contactId":     contactId,
				"content":       entity.Content,
				"contentType":   entity.ContentType,
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

func (r *noteRepository) CreateNoteForOrganization(ctx context.Context, tenant, organizationId string, entity entity.NoteEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.CreateNoteForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (org)-[:NOTED]->(n:Note {id:randomUUID()}) " +
		" ON CREATE SET n.content=$content, " +
		"				n.contentType=$contentType, " +
		"				n.createdAt=$now, " +
		"				n.updatedAt=$now, " +
		"				n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"				n:Note_%s," +
		"				n:TimelineEvent," +
		"				n:TimelineEvent_%s " +
		" RETURN n"

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
				"content":        entity.Content,
				"contentType":    entity.ContentType,
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

func (r *noteRepository) CreateNoteForMeeting(ctx context.Context, tenant, meetingId string, entity *entity.NoteEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.CreateNoteForMeeting")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	params, query := r.createMeetingQueryAndParams(tenant, meetingId, entity)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *noteRepository) CreateNoteForMeetingTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, meetingId string, entity *entity.NoteEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.CreateNoteForMeetingTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	params, query := r.createMeetingQueryAndParams(tenant, meetingId, entity)
	result, err := tx.Run(ctx, query, params)
	if err != nil {
		return nil, err
	}

	return utils.ExtractSingleRecordFirstValueAsNode(ctx, result, err)
}

func (r *noteRepository) Delete(ctx context.Context, tenant, noteId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (n:%s {id:$noteId}) DETACH DELETE n"

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

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

func (r *noteRepository) SetNoteCreator(ctx context.Context, tenant, userId, noteId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.SetNoteCreator")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), " +
		" (n:Note {id:$noteId})" +
		"  MERGE (u)-[:CREATED]->(n) "

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.GetNotedEntitiesForNotes")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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

func (r *noteRepository) GetMentionedEntitiesForNotes(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.GetMentionedEntitiesForNotes")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (n:Note_%s)-[:MENTIONED]->(:Tag)<-[:TAGGED]-(e:Issue_%s)
			WHERE n.id IN $ids
			RETURN e, n.id`, tenant, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
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

func (r *noteRepository) createMeetingQueryAndParams(tenant string, meetingId string, entity *entity.NoteEntity) (map[string]any, string) {
	query := "MATCH (m:Meeting_%s {id:$meetingId}) " +
		" MERGE (m)-[:NOTED]->(n:Note {id:randomUUID()}) " +
		" ON CREATE SET n.content=$content, " +
		"				n.contentType=$contentType, " +
		"				n.createdAt=$now, " +
		"				n.updatedAt=$now, " +
		"				n.source=$source, " +
		"				n.sourceOfTruth=$sourceOfTruth, " +
		"				n.appSource=$appSource, " +
		"				n:Note_%s," +
		"				n:TimelineEvent," +
		"				n:TimelineEvent_%s " +
		" RETURN n"
	params := map[string]any{
		"tenant":        tenant,
		"meetingId":     meetingId,
		"content":       entity.Content,
		"contentType":   entity.ContentType,
		"now":           utils.Now(),
		"source":        entity.Source,
		"sourceOfTruth": entity.SourceOfTruth,
		"appSource":     entity.AppSource,
	}
	return params, fmt.Sprintf(query, tenant, tenant, tenant)
}

func (r *noteRepository) GetMentionedByNotesForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteRepository.GetMentionedByNotesForIssues")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (i:Issue_%s)-[:TAGGED]->(tag:Tag)<-[:MENTIONED]-(n:Note_%s)
			WHERE i.id IN $issueIds
			RETURN n, i.id ORDER BY tag.name`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]any{
				"tenant":   tenant,
				"issueIds": issueIds,
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
