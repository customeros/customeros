package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantRepository interface {
	GetTenantsWithOrganizationsWithSlackChannels(ctx context.Context) ([]string, error)
	GetTenantDomain(ctx context.Context, tenant string) (string, error)
}

type tenantRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTenantRepository(driver *neo4j.DriverWithContext) TenantRepository {
	return &tenantRepository{
		driver: driver,
	}
}

func (r *tenantRepository) GetTenantsWithOrganizationsWithSlackChannels(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetTenantsWithOrganizationsWithSlackChannels")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)

	query := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization) 
				WHERE NOT o.slackChannelLink IS NULL AND o.slackChannelLink <> ''
				RETURN t.name`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	span.LogFields(log.Object("output - tenants:", records.([]string)))
	return records.([]string), err
}

func (r *tenantRepository) GetTenantDomain(ctx context.Context, tenant string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantRepository.GetTenantDomain")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)
	span.LogFields(log.String("tenant", tenant))

	query := `MATCH (t:Tenant {name:$tenant})--(w:Workspace) return w.name limit 1;`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	domain, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
		})
		return utils.ExtractSingleRecordFirstValueAsType[string](ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	span.LogFields(log.String("output - domain", domain.(string)))
	return domain.(string), err
}
