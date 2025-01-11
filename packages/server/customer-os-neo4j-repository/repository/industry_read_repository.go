package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type IndustryReadRepository interface {
	GetAllForOrganizationIds(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	GetByCode(ctx context.Context, code string) (*dbtype.Node, error)
}

type industryReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewIndustryReadRepository(driver *neo4j.DriverWithContext, database string) IndustryReadRepository {
	return &industryReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *industryReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *industryReadRepository) GetAllForOrganizationIds(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IndustryReadRepository.GetAllForOrganizationIds")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("organizationIds", fmt.Sprintf("%v", organizationIds)))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_INDUSTRY]->(i:Industry) 
				WHERE o.id in $organizationIds RETURN i, o.id`

	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *industryReadRepository) GetByCode(ctx context.Context, code string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IndustryReadRepository.GetByCode")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("code", code))

	cypher := `MATCH (i:Industry {code:$code}) RETURN i`
	params := map[string]any{
		"code": code,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if len(result.([]*dbtype.Node)) == 0 {
		return nil, nil
	}
	return result.(*dbtype.Node), nil
}
