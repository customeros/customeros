package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type CommonRepository interface {
	ExistsById(ctx context.Context, tenant, id, label string) (bool, error)
}

type commonRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommonRepository(driver *neo4j.DriverWithContext, database string) CommonRepository {
	return &commonRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commonRepository) ExistsById(ctx context.Context, tenant, id, label string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonRepository.ExistsById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("id", id), log.String("label", label))

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id:$id}) RETURN n.id LIMIT 1`, label, tenant)
	params := map[string]any{
		"id": id,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil
		}
	})
	if err != nil {
		return false, err
	}
	return result.(bool), err
}
