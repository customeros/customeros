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

type InteractionSessionReadRepository interface {
	GetByIdentifierAndChannel(ctx context.Context, tenant, identifier, channel string) (*neo4j.Node, error)
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

func (r *interactionSessionReadRepository) GetByIdentifierAndChannel(ctx context.Context, tenant, identifier, channel string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionReadRepository.GetByIdentifierAndChannel")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	tracing.LogObjectAsJson(span, "identifier", identifier)
	tracing.LogObjectAsJson(span, "channel", channel)

	cypher := fmt.Sprintf(`MATCH (i:InteractionSession_%s {identifier:$identifier, channel:$channel}) RETURN i`, tenant)
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

func (r *interactionSessionReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}
