package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionSessionReadRepository.GetForInteractionEvent")
	defer spans.Finish()

	spans.LogKV("interactionEventId", interactionEventId)

	cypher := fmt.Sprintf(`MATCH (e:InteractionEvent_%s{id: $id})-[:PART_OF]->(s:InteractionSession_%s) 
		 RETURN s`, tenant, tenant)
	params := map[string]any{
		"id": interactionEventId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionSessionReadRepository.GetAllForInteractionEvents")
	defer spans.Finish()

	spans.LogObjectAsJson("ids", ids)

	cypher := fmt.Sprintf(`MATCH (e:InteractionEvent)-[:PART_OF]->(s:InteractionSession_%s) 
		 WHERE e.id IN $ids AND e:InteractionEvent_%s
		 RETURN s, e.id`, tenant, tenant)
	params := map[string]any{
		"ids": ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionSessionReadRepository.GetByIdentifierAndChannel")
	defer spans.Finish()

	spans.LogObjectAsJson("identifier", identifier)
	spans.LogObjectAsJson("channel", channel)

	cypher := fmt.Sprintf(`MATCH (i:InteractionSession_%s {identifier:$identifier, channel:$channel}) RETURN i LIMIT 1`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"identifier": identifier,
		"channel":    channel,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionSessionReadRepository.GetAttendedByParticipantsForInteractionSessions")
	defer spans.Finish()

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
