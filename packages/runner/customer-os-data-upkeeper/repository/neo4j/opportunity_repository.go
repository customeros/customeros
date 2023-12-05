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

type TenantAndOpportunityId struct {
	Tenant        string
	OpportunityId string
}

type OpportunityRepository interface {
	GetRenewalOpportunitiesForClosingAsLost(ctx context.Context, referenceTime time.Time) ([]TenantAndOpportunityId, error)
}

type opportunityRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOpportunityRepository(driver *neo4j.DriverWithContext) OpportunityRepository {
	return &opportunityRepository{
		driver: driver,
	}
}

func (r *opportunityRepository) GetRenewalOpportunitiesForClosingAsLost(ctx context.Context, referenceTime time.Time) ([]TenantAndOpportunityId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.GetRenewalOpportunitiesForClosingAsLost")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(span)
	span.LogFields(log.Object("referenceTime", referenceTime))

	cypher := `MATCH (t:Tenant)<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:ACTIVE_RENEWAL]->(op:RenewalOpportunity)
				WHERE c.status = $endedStatus AND c.endedAt < $referenceTime AND op.internalStage = $internalStageOpen
				RETURN t.name, op.id LIMIT 100`
	params := map[string]any{
		"referenceTime":     referenceTime,
		"endedStatus":       "ENDED",
		"internalStageOpen": "OPEN",
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
	output := make([]TenantAndOpportunityId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantAndOpportunityId{
				Tenant:        v.Values[0].(string),
				OpportunityId: v.Values[1].(string),
			})
	}
	span.LogFields(log.Int("output - length", len(output)))
	return output, nil
}
