package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type InteractionEventReadRepository interface {
	GetInteractionEvent(ctx context.Context, tenant, interactionEventId string) (*dbtype.Node, error)
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
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
