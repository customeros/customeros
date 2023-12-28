package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/customer-os-neo4j-repository/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type CommonReadRepository interface {
	ExistsById(ctx context.Context, tenant, id, label string) (bool, error)
}

type commonReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommonReadRepository(driver *neo4j.DriverWithContext, database string) CommonReadRepository {
	return &commonReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commonReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *commonReadRepository) ExistsById(ctx context.Context, tenant, id, label string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonReadRepository.ExistsById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("id", id), log.String("label", label))

	cypher := fmt.Sprintf(`MATCH (n:%s {id:$id}) WHERE n:%s_%s RETURN n.id LIMIT 1`, label, label, tenant)
	params := map[string]any{
		"id": id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result.exists", result.(bool)))
	return result.(bool), err
}
