package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type TenantAndOrganizationId struct {
	Tenant         string
	OrganizationId string
}

type OrganizationRepository interface {
	GetOrganizationsForNextCycleDateRenew(ctx context.Context, referenceTime time.Time) ([]TenantAndOrganizationId, error)
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) GetOrganizationsForNextCycleDateRenew(ctx context.Context, referenceTime time.Time) ([]TenantAndOrganizationId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRepository.GetOrganizationsForNextCycleDateRenew")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)
	span.LogFields(log.Object("referenceTime", referenceTime))

	cypher := `MATCH (t:Tenant)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)
				WHERE org.billingDetailsRenewalCycleNext <= $referenceTime
				RETURN t.name, org.id ORDER BY org.billingDetailsRenewalCycleNext ASC LIMIT 100`
	params := map[string]any{
		"referenceTime": referenceTime,
	}
	span.LogFields(log.String("query", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantAndOrganizationId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOrganizationId{
				Tenant:         v.Values[0].(string),
				OrganizationId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("output - length", len(output)))
	return output, nil
}
