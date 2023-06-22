package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type HealthIndicatorRepository interface {
	CreateDefaultHealthIndicatorsForNewTenant(ctx context.Context, tenant string) error
}

type healthIndicatorRepository struct {
	driver *neo4j.DriverWithContext
}

func NewHealthIndicatorRepository(driver *neo4j.DriverWithContext) HealthIndicatorRepository {
	return &healthIndicatorRepository{
		driver: driver,
	}
}

func (r *healthIndicatorRepository) CreateDefaultHealthIndicatorsForNewTenant(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HealthIndicatorRepository.CreateDefaultHealthIndicatorsForNewTenant")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`WITH [
					{name: 'Green', order:10},
					{name: 'Yellow', order:20},
					{name: 'Orange', order:30},
					{name: 'Red', order:40}] AS indicators
				UNWIND indicators AS indicator
				MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:HEALTH_INDICATOR_BELONGS_TO_TENANT]-(h:HealthIndicator {name:indicator.name})
				ON CREATE SET 	h.id=randomUUID(), 
								h.order=indicator.order,
								h.createdAt=$now, 
								h:HealthIndicator_%s`, tenant)

	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"now":    utils.Now(),
			})
		return nil, err
	})
	return err
}
