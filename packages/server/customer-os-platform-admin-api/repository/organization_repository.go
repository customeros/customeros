package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OrganizationRepository interface {
	CountOrganizationsForLastTouchpointRefresh(ctx context.Context, tenant string) (int64, error)
	GetOrganizationsForLastTouchpointRefresh(ctx context.Context, tenant string, skip, limit int) ([]string, error)
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) CountOrganizationsForLastTouchpointRefresh(ctx context.Context, tenant string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.CountOrganizationsForLastTouchpointRefresh")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {hide: false}) RETURN count(o)`
	span.LogFields(log.String("cypher", cypher))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher,
			map[string]any{
				"tenant": tenant,
			})
		return utils.ExtractSingleRecordFirstValue(ctx, queryResult, err)
	})
	if err != nil {
		return 0, err
	}
	return dbRecord.(int64), err
}

func (r *organizationRepository) GetOrganizationsForLastTouchpointRefresh(ctx context.Context, tenant string, skip, limit int) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationsForLastTouchpointRefresh")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("skip", skip))
	span.LogFields(log.Int("limit", limit))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {hide: false}) WHERE o:Organization_%s RETURN o.id SKIP $skip LIMIT $limit`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"skip":   skip,
				"limit":  limit,
			})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.([]string), err
}
