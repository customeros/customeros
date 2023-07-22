package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/tracing"
	local_utils "github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantRepository interface {
	GetTenantsWithOrganizations(ctx context.Context, minOrganizations int) ([]string, error)
}

type tenantRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTenantRepository(driver *neo4j.DriverWithContext) TenantRepository {
	return &tenantRepository{
		driver: driver,
	}
}

func (r *tenantRepository) GetTenantsWithOrganizations(ctx context.Context, minOrganizations int) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetTenantsWithOrganizations")
	defer span.Finish()

	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization) 
				WITH t, count(o) as orgsCount 
				WHERE orgsCount >= $limit
				RETURN t.name order by orgsCount desc;`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"limit": minOrganizations,
		})
		return local_utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return records.([]string), err
}
