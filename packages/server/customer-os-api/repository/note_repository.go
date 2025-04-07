package repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
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
	GetNotesForMeetings(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)

	CreateNoteForMeeting(ctx context.Context, tenant, meeting string, entity *entity.NoteEntity) (*dbtype.Node, error)
	CreateNoteForMeetingTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, meeting string, entity *entity.NoteEntity) (*dbtype.Node, error)

	UpdateNote(ctx context.Context, session neo4j.SessionWithContext, tenant string, entity entity.NoteEntity) (*dbtype.Node, error)

	Delete(ctx context.Context, tenant, noteId string) error
	SetNoteCreator(ctx context.Context, tenant, userId, noteId string) error
}

type noteRepository struct {
	driver *neo4j.DriverWithContext
}

func NewNoteRepository(driver *neo4j.DriverWithContext) NoteRepository {
	return &noteRepository{
		driver: driver,
	}
}

func (r *noteRepository) GetNotesForMeetings(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "NoteRepository.GetNotesForMeetings")
	defer spans.Finish()

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "NoteRepository.UpdateNote")
	defer spans.Finish()

	query := "MATCH (n:%s {id:$noteId}) " +
		" SET 	n.content=$content, " +
		"		n.contentType=$contentType, " +
		"		n.updatedAt=datetime() " +
		" RETURN n"
	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		txResult, err := tx.Run(ctx, fmt.Sprintf(query, "Note_"+tenant),
			map[string]interface{}{
				"tenant":      tenant,
				"noteId":      entity.Id,
				"content":     entity.Content,
				"contentType": entity.ContentType,
				"now":         utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, txResult, err)
	})
	if err != nil {
		return nil, err
	}
	return queryResult.(*dbtype.Node), nil
}

func (r *noteRepository) CreateNoteForMeeting(ctx context.Context, tenant, meetingId string, entity *entity.NoteEntity) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "NoteRepository.CreateNoteForMeeting")
	defer spans.Finish()

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "NoteRepository.CreateNoteForMeetingTx")
	defer spans.Finish()

	params, query := r.createMeetingQueryAndParams(tenant, meetingId, entity)
	result, err := tx.Run(ctx, query, params)
	if err != nil {
		return nil, err
	}

	return utils.ExtractSingleRecordFirstValueAsNode(ctx, result, err)
}

func (r *noteRepository) Delete(ctx context.Context, tenant, noteId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "NoteRepository.Delete")
	defer spans.Finish()

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "NoteRepository.SetNoteCreator")
	defer spans.Finish()

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

func (r *noteRepository) createMeetingQueryAndParams(tenant string, meetingId string, entity *entity.NoteEntity) (map[string]any, string) {
	query := "MATCH (m:Meeting_%s {id:$meetingId}) " +
		" MERGE (m)-[:NOTED]->(n:Note {id:randomUUID()}) " +
		" ON CREATE SET n.content=$content, " +
		"				n.contentType=$contentType, " +
		"				n.createdAt=$now, " +
		"				n.updatedAt=datetime(), " +
		"				n.source=$source, " +
		"				n.appSource=$appSource, " +
		"				n:Note_%s," +
		"				n:TimelineEvent," +
		"				n:TimelineEvent_%s " +
		" RETURN n"
	params := map[string]any{
		"tenant":      tenant,
		"meetingId":   meetingId,
		"content":     entity.Content,
		"contentType": entity.ContentType,
		"now":         utils.Now(),
		"source":      entity.Source,
		"appSource":   entity.AppSource,
	}
	return params, fmt.Sprintf(query, tenant, tenant, tenant)
}
