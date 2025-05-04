package neo4j_repository

import (
	context2 "context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type IndustryReadRepository interface {
	GetAllForOrganizationIds(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	GetInUseIndustries(ctx context2.Context, tenant string) ([]*dbtype.Node, error)
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

func (r *industryReadRepository) GetAllForOrganizationIds(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IndustryReadRepository.GetAllForOrganizationIds")
	defer spans.Finish()

	spans.LogKV("organizationIds", fmt.Sprintf("%v", organizationIds))

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_INDUSTRY]->(i:Industry) 
				WHERE o.id in $organizationIds RETURN i, o.id`

	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	return result.([]*utils.DbNodeAndId), err
}

func (r *industryReadRepository) GetInUseIndustries(ctx context2.Context, tenant string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IndustryReadRepository.GetInUseIndustries")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {hide:false})-[:HAS_INDUSTRY]->(i:Industry) 
				RETURN DISTINCT i ORDER BY i.code`

	params := map[string]any{
		"tenant": tenant,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), err
}

func (r *industryReadRepository) GetByCode(ctx context.Context, code string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IndustryReadRepository.GetByCode")
	defer spans.Finish()

	spans.LogKV("code", code)

	cypher := `MATCH (i:Industry {code:$code}) RETURN i`
	params := map[string]any{
		"code": code,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	if len(result.([]*dbtype.Node)) == 0 {
		return nil, nil
	}
	return result.([]*dbtype.Node)[0], nil
}
