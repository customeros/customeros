package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
	"strings"
)

type MeetingRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error)
	LinkWithParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, meetingId, participantId string, sentType *string, relation entity.MeetingRelation) error
	GetParticipantsForMeetings(ctx context.Context, tenant string, ids []string, relation entity.MeetingRelation) ([]*utils.DbNodeWithRelationAndId, error)
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
		"				m.name=$name, " +
		"				m.createdAt=$now, " +
		"				m.updatedAt=$now, " +
		"				m.start=$now, " +
		"				m.end=$now, " +
		"				m.appSource=$appSource, " +
		"				m.source=$source, " +
		"				m.sourceOfTruth=$sourceOfTruth " +
		" RETURN m"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]any{
			"name":          entity.Name,
			"now":           entity.CreatedAt,
			"appSource":     entity.AppSource,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *meetingRepository) LinkWithParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, meetingId, participantId string, sentType *string, relation entity.MeetingRelation) error {
	query := ""
	switch entityType {
	case entity.CONTACT:
		query = fmt.Sprintf(`MATCH (p:Contact_%s {id:$participantId}) `, tenant)
	case entity.USER:
		query = fmt.Sprintf(`MATCH (p:User_%s {id:$participantId}) `, tenant)
	}
	query += fmt.Sprintf(`MATCH (is:Meeting_%s {id:$meetingId}) `, tenant)

	if sentType != nil {
		query += fmt.Sprintf(`MERGE (m)<-[r:%s {type:$sentType}]-(p) RETURN r`, relation)
	} else {
		query += fmt.Sprintf(`MERGE (m)<-[r:%s]-(p) RETURN r`, relation)
	}

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"participantId": participantId,
			"meetingId":     meetingId,
			"sentType":      sentType,
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

	query := "MATCH (is:Meeting_%s)-[rel:%s]->(p) " +
		" WHERE is.id IN $ids " +
		" RETURN p, rel, is.id"

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

func (r *meetingRepository) Update(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.MeetingEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query, params := r.createQueryAndParams(tenant, entity)

	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		txResult, err := tx.Run(ctx, fmt.Sprintf(query, "Meeting_"+tenant), params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, txResult, err)
	})
	if err != nil {
		return nil, err
	}
	return queryResult.(*dbtype.Node), nil
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
	if entity.Start != nil {
		params["start"] = *entity.Start
		qb.WriteString("	m.start=$start, ")
	}
	if entity.End != nil {
		qb.WriteString("	m.end=$end, ")
		params["end"] = *entity.End
	}
	if entity.Location != nil {
		qb.WriteString("	m.location=$location, ")
		params["location"] = entity.Location
	}

	if entity.Agenda != nil {
		qb.WriteString("	m.agenda=$agenda, ")
		params["agenda"] = entity.Agenda
	}
	if entity.AgendaContentType != nil {
		qb.WriteString("	m.agendaContentType=$agendaContentType, ")
		params["agendaContentType"] = entity.AgendaContentType
	}

	qb.WriteString("	m.updatedAt=$now ")
	qb.WriteString(" RETURN m")

	return qb.String(), params
}
