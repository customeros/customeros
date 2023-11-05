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
)

type InteractionSessionRepository interface {
	GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.InteractionSessionEntity) (*dbtype.Node, error)
	GetAttendedByParticipantsForInteractionSessions(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error)

	LinkWithAttendedByEmailInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionSessionId, email string, sentType *string) error
	LinkWithAttendedByPhoneNumberInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionSessionId, e164 string, sentType *string) error
	LinkWithAttendedByParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, interactionSessionId, participantId string, sentType *string) error
}

type interactionSessionRepository struct {
	driver *neo4j.DriverWithContext
}

func NewInteractionSessionRepository(driver *neo4j.DriverWithContext) InteractionSessionRepository {
	return &interactionSessionRepository{
		driver: driver,
	}
}

func (r *interactionSessionRepository) LinkWithAttendedByEmailInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionSessionId, email string, sentType *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionRepository.LinkWithAttendedByEmailInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (e:Email_%s) `, tenant)
	query += fmt.Sprintf(`MATCH (is:InteractionSession_%s {id:$sessionId}) `, tenant)
	query += `WHERE e.email = $email OR e.rawEmail = $email `

	if sentType != nil {
		query += fmt.Sprintf(`MERGE (is)-[r:ATTENDED_BY {type:$sentType}]->(e) RETURN r`)
	} else {
		query += fmt.Sprintf(`MERGE (is)-[r:ATTENDED_BY]->(e) RETURN r`)
	}
	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"email":     email,
			"sessionId": interactionSessionId,
			"sentType":  sentType,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *interactionSessionRepository) LinkWithAttendedByPhoneNumberInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionSessionId, e164 string, sentType *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionRepository.LinkWithAttendedByPhoneNumberInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (p:PhoneNumber_%s) `, tenant)
	query += fmt.Sprintf(`MATCH (is:InteractionSession_%s {id:$sessionId}) `, tenant)
	query += `WHERE p.e164 = $e164 OR p.rawPhoneNumber = $e164 `

	if sentType != nil {
		query += fmt.Sprintf(`MERGE (is)-[r:ATTENDED_BY {type:$sentType}]->(p) RETURN r`)
	} else {
		query += fmt.Sprintf(`MERGE (is)-[r:ATTENDED_BY]->(p) RETURN r`)
	}

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"e164":      e164,
			"sessionId": interactionSessionId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *interactionSessionRepository) LinkWithAttendedByParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType entity.EntityType, interactionSessionId, participantId string, sentType *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionRepository.LinkWithAttendedByParticipantInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := ""
	switch entityType {
	case entity.CONTACT:
		query = fmt.Sprintf(`MATCH (p:Contact_%s {id:$participantId}) `, tenant)
	case entity.USER:
		query = fmt.Sprintf(`MATCH (p:User_%s {id:$participantId}) `, tenant)
	}
	query += fmt.Sprintf(`MATCH (is:InteractionSession_%s {id:$sessionId}) `, tenant)

	if sentType != nil {
		query += fmt.Sprintf(`MERGE (is)<-[r:ATTENDED_BY {type:$sentType}]-(p) RETURN r`)
	} else {
		query += fmt.Sprintf(`MERGE (is)<-[r:ATTENDED_BY]-(p) RETURN r`)
	}

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"participantId": participantId,
			"sessionId":     interactionSessionId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *interactionSessionRepository) GetAttendedByParticipantsForInteractionSessions(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionRepository.GetAttendedByParticipantsForInteractionSessions")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (is:InteractionSession_%s)-[rel:ATTENDED_BY]->(p) " +
		" WHERE is.id IN $ids " +
		" RETURN p, rel, is.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
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

func (r *interactionSessionRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.InteractionSessionEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MERGE (is:InteractionSession_%s {id:randomUUID()}) " +
		" ON CREATE SET is:InteractionSession, " +
		"				is.identifier=$identifier, " +
		"				is.source=$source, " +
		"				is.channel=$channel, " +
		"				is.channelData=$channelData, " +
		"				is.createdAt=$now, " +
		"				is.updatedAt=$now, " +
		"				is.name=$name, " +
		" 				is.status=$status, " +
		"				is.type=$type, " +
		"				is.sourceOfTruth=$sourceOfTruth, " +
		"				is.appSource=$appSource " +
		" RETURN is"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]any{
			"identifier":    entity.SessionIdentifier,
			"source":        entity.Source,
			"channel":       entity.Channel,
			"channelData":   entity.ChannelData,
			"now":           entity.CreatedAt,
			"name":          entity.Name,
			"status":        entity.Status,
			"type":          entity.Type,
			"sourceOfTruth": entity.SourceOfTruth,
			"appSource":     entity.AppSource,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *interactionSessionRepository) GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionRepository.GetAllForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	cypher := fmt.Sprintf(`MATCH (e:InteractionEvent)-[:PART_OF]->(s:InteractionSession_%s) 
		 WHERE e.id IN $ids AND e:InteractionEvent_%s
		 RETURN s, e.id`, tenant, tenant)
	params := map[string]any{
		"ids": ids,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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
