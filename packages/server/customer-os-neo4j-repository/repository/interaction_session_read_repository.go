package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type InteractionSessionReadRepository interface {
	GetForInteractionEvent(ctx context.Context, tenant, interactionEventId string) (*neo4j.Node, error)
	GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetByIdentifierAndChannel(ctx context.Context, tenant, identifier, channel string) (*neo4j.Node, error)
	GetAttendedByParticipantsForInteractionSessions(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error)
}

type interactionSessionReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInteractionSessionReadRepository(driver *neo4j.DriverWithContext, database string) InteractionSessionReadRepository {
	return &interactionSessionReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *interactionSessionReadRepository) GetForInteractionEvent(ctx context.Context, tenant, interactionEventId string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionReadRepository.GetForInteractionEvent")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.Object("interactionEventId", interactionEventId))

	cypher := fmt.Sprintf(`MATCH (e:InteractionEvent_%s{id: $id})-[:PART_OF]->(s:InteractionSession_%s) 
		 RETURN s`, tenant, tenant)
	params := map[string]any{
		"id": interactionEventId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *interactionSessionReadRepository) GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionReadRepository.GetAllForInteractionEvents")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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

func (r *interactionSessionReadRepository) GetByIdentifierAndChannel(ctx context.Context, tenant, identifier, channel string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionReadRepository.GetByIdentifierAndChannel")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "identifier", identifier)
	tracing.LogObjectAsJson(span, "channel", channel)

	cypher := fmt.Sprintf(`MATCH (i:InteractionSession_%s {identifier:$identifier, channel:$channel}) RETURN i LIMIT 1`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"identifier": identifier,
		"channel":    channel,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *interactionSessionReadRepository) GetAttendedByParticipantsForInteractionSessions(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionReadRepository.GetAttendedByParticipantsForInteractionSessions")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`
				MATCH (is:InteractionSession_%s)-[rel:ATTENDED_BY]->(p)
				WHERE is.id IN $ids
				RETURN distinct(p), rel, is.id
				UNION
				MATCH (is:InteractionSession_%s)<-[:PART_OF]-(ie:InteractionEvent_%s)-[rel:SENT_BY]->(p)
				WHERE is.id IN $ids
				RETURN distinct(p), rel, is.id
				UNION
				MATCH (is:InteractionSession_%s)<-[:PART_OF]-(ie:InteractionEvent_%s)-[rel:SENT_TO]->(p)
				WHERE is.id IN $ids
				RETURN distinct(p), rel, is.id`, tenant, tenant, tenant, tenant, tenant)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
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

func (r *interactionSessionReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}
