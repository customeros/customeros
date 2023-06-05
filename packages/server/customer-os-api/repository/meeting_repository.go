package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

type MeetingRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error)
	LinkWithParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, meetingId, participantId string, entityType entity.EntityType, relation entity.MeetingRelation) error
	UnlinkParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, meetingId, participantId string, entityType entity.EntityType, relation entity.MeetingRelation) error
	GetParticipantsForMeetings(ctx context.Context, tenant string, ids []string, relation entity.MeetingRelation) ([]*utils.DbNodeWithRelationAndId, error)
	GetMeetingForInteractionEvent(ctx context.Context, tenant string, id string) (*dbtype.Node, error)
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
		"				m.sourceOfTruth=$sourceOfTruth " +
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
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *meetingRepository) LinkWithParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, meetingId, participantId string, entityType entity.EntityType, relation entity.MeetingRelation) error {
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

	qb.WriteString("	m.updatedAt=$now ")
	qb.WriteString(" RETURN m")

	return qb.String(), params
}
