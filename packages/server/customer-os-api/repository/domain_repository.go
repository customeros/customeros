package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type DomainRepository interface {
	GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
}

type domainRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewDomainRepository(driver *neo4j.DriverWithContext, database string) DomainRepository {
	return &domainRepository{
		driver:   driver,
		database: database,
	}
}

func (r *domainRepository) GetForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainRepository.GetForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain)
			WHERE o.id IN $organizationIds
			RETURN d, o.id ORDER BY d.domain`
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	result, err := utils.ExecuteQuery(ctx, *r.driver, r.database, cypher, params)
	if err != nil {
		return nil, err
	}
	return utils.ExtractAllRecordsAsDbNodeAndIdFromEagerResult(result), nil
}
