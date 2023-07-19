package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type MeetingRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error)
	LinkWithParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, meetingId, participantId string, entityType entity.EntityType, relation entity.MeetingRelation) error
	UnlinkParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, meetingId, participantId string, entityType entity.EntityType, relation entity.MeetingRelation) error
	GetParticipantsForMeetings(ctx context.Context, tenant string, ids []string, relation entity.MeetingRelation) ([]*utils.DbNodeWithRelationAndId, error)
	GetMeetingForInteractionEvent(ctx context.Context, tenant string, id string) (*dbtype.Node, error)
	GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetPaginatedMeetings(ctx context.Context, session neo4j.SessionWithContext, externalSystemID string, externalID *string, tenant, userEmail string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
}

type meetingRepository struct {
	driver *neo4j.DriverWithContext
}

func NewMeetingRepository(driver *neo4j.DriverWithContext) MeetingRepository {
	return &meetingRepository{
		driver: driver,
	}
}

func (r *meetingRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MERGE (m:Meeting_%s {id:randomUUID()}) " +
		" ON CREATE SET m:Meeting, " +
		" 				m:TimelineEvent, " +
		" 				m:TimelineEvent_%s, " +
		"				m.name=$name, " +
		"				m.agenda=$agenda, " +
		"				m.agendaContentType=$agendaContentType, " +
		"				m.conferenceUrl=$conferenceUrl, " +
		"				m.meetingExternalUrl=$meetingExternalUrl, " +
		"				m.createdAt=$createdAt, " +
		"				m.updatedAt=$updatedAt, " +
		"				m.startedAt=$startedAt, " +
		"				m.endedAt=$endedAt, " +
		"				m.appSource=$appSource, " +
		"				m.source=$source, " +
		"				m.sourceOfTruth=$sourceOfTruth, " +
		"				m.status=$status " +
		" RETURN m"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
		map[string]any{
			"name":               entity.Name,
			"agenda":             utils.IfNotNilStringWithDefault(entity.Agenda, ""),
			"agendaContentType":  utils.IfNotNilStringWithDefault(entity.AgendaContentType, ""),
			"conferenceUrl":      utils.IfNotNilStringWithDefault(entity.ConferenceUrl, ""),
			"meetingExternalUrl": utils.IfNotNilStringWithDefault(entity.MeetingExternalUrl, ""),
			"createdAt":          entity.CreatedAt,
			"updatedAt":          entity.CreatedAt,
			"startedAt":          utils.IfNotNilTimeWithDefault(entity.StartedAt, utils.Now()),
			"endedAt":            utils.IfNotNilTimeWithDefault(entity.EndedAt, utils.Now()),
			"appSource":          entity.AppSource,
			"source":             entity.Source,
			"sourceOfTruth":      entity.SourceOfTruth,
			"status":             entity.Status,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *meetingRepository) LinkWithParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, meetingId, participantId string, entityType entity.EntityType, relation entity.MeetingRelation) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.LinkWithParticipantInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = fmt.Sprintf(`MATCH (p:Contact_%s {id:$participantId}) `, tenant)
	case entity.USER:
		query = fmt.Sprintf(`MATCH (p:User_%s {id:$participantId}) `, tenant)
	case entity.ORGANIZATION:
		query = fmt.Sprintf(`MATCH (p:Organization_%s {id:$participantId}) `, tenant)
	}
	query += fmt.Sprintf(`MATCH (m:Meeting_%s {id:$meetingId}) `, tenant)
	query += fmt.Sprintf(`MERGE (m)-[r:%s]->(p) RETURN r`, relation)

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"participantId": participantId,
			"meetingId":     meetingId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *meetingRepository) UnlinkParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, meetingId, participantId string, entityType entity.EntityType, relation entity.MeetingRelation) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.UnlinkParticipantInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = fmt.Sprintf(`MATCH (p:Contact_%s {id:$participantId}) `, tenant)
	case entity.USER:
		query = fmt.Sprintf(`MATCH (p:User_%s {id:$participantId}) `, tenant)
	case entity.ORGANIZATION:
		query = fmt.Sprintf(`MATCH (p:Organization_%s {id:$participantId}) `, tenant)
	}
	query += fmt.Sprintf(`MATCH (m:Meeting_%s {id:$meetingId}) `, tenant)
	query += fmt.Sprintf(`MATCH (m)-[r:%s]->(p) DELETE r return m`, relation)

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"participantId": participantId,
			"meetingId":     meetingId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *meetingRepository) GetParticipantsForMeetings(ctx context.Context, tenant string, ids []string, relation entity.MeetingRelation) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.GetParticipantsForMeetings")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (m:Meeting_%s)-[rel:%s]->(p) " +
		" WHERE m.id IN $ids AND (p:Contact OR p:User OR p:Organization)" +
		" RETURN p, rel, m.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, relation),
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}

	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *meetingRepository) GetMeetingForInteractionEvent(ctx context.Context, tenant string, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.GetMeetingForInteractionEvent")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (m:Meeting_%s)<-[:PART_OF]-(e:InteractionEvent) " +
		" WHERE e.id= $id " +
		" RETURN m, e.id"
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"tenant": tenant,
				"id":     id,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}

	convertedResult, isOk := result.([]*dbtype.Node)
	if !isOk {
		return nil, errors.New("GetMeetingForInteractionEvent: cannot convert result")
	}
	if len(convertedResult) == 0 {
		return nil, nil
	}
	return convertedResult[0], err
}

func (r *meetingRepository) Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query, params := r.createQueryAndParams(tenant, entity)

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Meeting_"+tenant), params)
	if err != nil {
		return nil, err
	}

	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *meetingRepository) createQueryAndParams(tenant string, entity *entity.MeetingEntity) (string, map[string]interface{}) {
	var qb strings.Builder
	params := map[string]interface{}{
		"tenant":    tenant,
		"meetingId": entity.Id,
		"now":       utils.Now(),
	}

	qb.WriteString("MATCH (m:%s {id:$meetingId}) ")
	qb.WriteString(" SET ")
	if entity.Name != nil {
		qb.WriteString("	m.name=$name, ")
		params["name"] = entity.Name
	}
	if entity.StartedAt != nil {
		params["startedAt"] = *entity.StartedAt
		qb.WriteString("	m.startedAt=$startedAt, ")
	}
	if entity.EndedAt != nil {
		qb.WriteString("	m.endedAt=$endedAt, ")
		params["endedAt"] = *entity.EndedAt
	}
	if entity.ConferenceUrl != nil {
		qb.WriteString("	m.conferenceUrl=$conferenceUrl, ")
		params["conferenceUrl"] = entity.ConferenceUrl
	}
	if entity.MeetingExternalUrl != nil {
		qb.WriteString("	m.meetingExternalUrl=$meetingExternalUrl, ")
		params["meetingExternalUrl"] = entity.MeetingExternalUrl
	}

	if entity.Agenda != nil {
		qb.WriteString("	m.agenda=$agenda, ")
		params["agenda"] = entity.Agenda
	}

	if entity.AgendaContentType != nil {
		qb.WriteString("	m.agendaContentType=$agendaContentType, ")
		params["agendaContentType"] = entity.AgendaContentType
	}

	if entity.Recording != nil {
		qb.WriteString("	m.recording=$recording, ")
		params["recording"] = entity.Recording
	}

	if entity.Status != nil {
		qb.WriteString("	m.status=$status, ")
		params["status"] = entity.Status
	}

	qb.WriteString("	m.updatedAt=$now ")
	qb.WriteString(" RETURN m")

	return qb.String(), params
}

func (r *meetingRepository) GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.GetAllForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (e:InteractionEvent_%s)-[:PART_OF]->(m:Meeting) 
		 WHERE e.id IN $ids 
		 RETURN m, e.id`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
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
	return result.([]*utils.DbNodeAndId), err
}

func (r *meetingRepository) GetPaginatedMeetings(ctx context.Context, session neo4j.SessionWithContext, externalSystemID string, externalID *string, tenant, userEmail string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MeetingRepository.GetPaginatedMeetings")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("m")
		countParams := map[string]any{
			"externalId":       externalID,
			"externalSystemId": externalSystemID,
			"userEmail":        userEmail,
			"tenant":           tenant,
			"email":            userEmail,
		}
		var qb strings.Builder
		qb.WriteString("MATCH (m:Meeting_%s)")
		if externalID != nil {
			qb.WriteString("-[:IS_LINKED_WITH {externalId: $externalId}]->")
		} else {
			qb.WriteString("-[:IS_LINKED_WITH]->")
		}
		qb.WriteString("(ex:ExternalSystem_%s {id: $externalSystemId})")

		utils.MergeMapToMap(filterParams, countParams)
		countQuery := fmt.Sprintf(qb.String()+" %s RETURN count(m) as count", tenant, tenant, filterCypherStr)
		queryResult, err := tx.Run(ctx, countQuery, countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"externalId":       externalID,
			"externalSystemId": externalSystemID,
			"tenant":           tenant,
			"email":            userEmail,
			"skip":             skip,
			"limit":            limit,
		}
		utils.MergeMapToMap(filterParams, params)

		returnQuery := fmt.Sprintf(qb.String()+" %s RETURN m %s SKIP $skip LIMIT $limit", tenant, tenant, filterCypherStr, sort.SortingCypherFragment("m"))
		queryResult, err = tx.Run(ctx, returnQuery, params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}
