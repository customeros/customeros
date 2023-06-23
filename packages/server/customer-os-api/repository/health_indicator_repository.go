package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type HealthIndicatorRepository interface {
	CreateDefaultHealthIndicatorsForNewTenant(ctx context.Context, tenant string) error
	GetAll(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetAllForOrganizations(ctx context.Context, ids []string) ([]*utils.DbNodeAndId, error)
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

func (r *healthIndicatorRepository) GetAll(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HealthIndicatorRepository.GetAll")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:HEALTH_INDICATOR_BELONGS_TO_TENANT]-(h:HealthIndicator)
			RETURN h ORDER BY h.order ASC`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	return result.([]*dbtype.Node), err
}

func (r *healthIndicatorRepository) GetAllForOrganizations(ctx context.Context, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HealthIndicatorRepository.GetAllForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant}),
				(t)<-[:HEALTH_INDICATOR_BELONGS_TO_TENANT]-(h:HealthIndicator)<-[:HAS_INDICATOR]-(o:Organization)-[:ORGANIZATION_BELONGS_TO_TENANT]->(t)
		 		WHERE o.id IN $ids 
		 		RETURN h, o.id`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": common.GetTenantFromContext(ctx),
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}
