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

type InteractionEventReadRepository interface {
	GetInteractionEvent(ctx context.Context, tenant, interactionEventId string) (*dbtype.Node, error)
	GetInteractionEventByCustomerOSIdentifier(ctx context.Context, customerOSInternalIdentifier string) (*dbtype.Node, error)
	InteractionEventSentByUser(ctx context.Context, tenant, interactionEventId string) (bool, error)
}

type interactionEventReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInteractionEventReadRepository(driver *neo4j.DriverWithContext, database string) InteractionEventReadRepository {
	return &interactionEventReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *interactionEventReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *interactionEventReadRepository) GetInteractionEvent(ctx context.Context, tenant, interactionEventId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventReadRepository.GetInteractionEvent")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent {id:$id}) WHERE i:InteractionEvent_%s RETURN i`, tenant)
	params := map[string]any{
		"id": interactionEventId,
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
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *interactionEventReadRepository) GetInteractionEventByCustomerOSIdentifier(ctx context.Context, customerOSInternalIdentifier string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventReadRepository.GetInteractionEventByCustomerOSIdentifier")
	defer span.Finish()
	span.LogFields(log.String("customerOSInternalIdentifier", customerOSInternalIdentifier))

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent {customerOSInternalIdentifier:$customerOSInternalIdentifier}) WHERE i:InteractionEvent RETURN i`)
	params := map[string]any{
		"customerOSInternalIdentifier": customerOSInternalIdentifier,
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

func (r *interactionEventReadRepository) InteractionEventSentByUser(ctx context.Context, tenant, interactionEventId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventReadRepository.InteractionEventSentByUser")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent {id:$id}) WHERE i:InteractionEvent_%s AND (i)-[:SENT_BY]->(:User) OR (i)-[:SENT_BY]->(:Email|PhoneNumber)--(:User) return count(i) > 0`, tenant)
	params := map[string]any{
		"id": interactionEventId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsType[bool](ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result", result.(bool)))
	return result.(bool), nil
}
