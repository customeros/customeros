package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type TagReadRepository interface {
	GetByNameOptional(ctx context.Context, tenant, name string) (*dbtype.Node, error)
}

type tagReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewTagReadRepository(driver *neo4j.DriverWithContext, database string) TagReadRepository {
	return &tagReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *tagReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *tagReadRepository) GetByNameOptional(ctx context.Context, tenant, name string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagReadRepository.GetByNameOptional")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:$name}) return tag limit 1`
	params := map[string]any{
		"name": name,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}
