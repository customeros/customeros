package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type MasterPlanReadRepository interface {
	GetMasterPlanById(ctx context.Context, tenant, masterPlanId string) (*dbtype.Node, error)
	GetMasterPlanMilestoneById(ctx context.Context, tenant, masterPlanMilestoneId string) (*dbtype.Node, error)
	GetMasterPlansOrderByCreatedAt(ctx context.Context, tenant string, returnRetired *bool) ([]*dbtype.Node, error)
	GetMasterPlanMilestonesForMasterPlans(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
}

type masterPlanReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewMasterPlanReadRepository(driver *neo4j.DriverWithContext, database string) MasterPlanReadRepository {
	return &masterPlanReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *masterPlanReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *masterPlanReadRepository) GetMasterPlanById(ctx context.Context, tenant, masterPlanId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanReadRepository.GetMasterPlanById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan {id:$id}) RETURN mp`
	params := map[string]any{
		"tenant": tenant,
		"id":     masterPlanId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *masterPlanReadRepository) GetMasterPlanMilestoneById(ctx context.Context, tenant, masterPlanMilestoneId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanReadRepository.GetMasterPlanMilestoneById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, masterPlanMilestoneId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(:MasterPlan)-[:HAS_MILESTONE]->(m:MasterPlanMilestone {id:$id}) RETURN m`
	params := map[string]any{
		"tenant": tenant,
		"id":     masterPlanMilestoneId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

func (r *masterPlanReadRepository) GetMasterPlansOrderByCreatedAt(ctx context.Context, tenant string, returnRetired *bool) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanReadRepository.GetMasterPlansOrderByCreatedAt")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan)`
	if returnRetired != nil {
		if *returnRetired {
			cypher += ` WHERE mp.retired = true`
		} else {
			cypher += ` WHERE mp.retired IS NULL OR mp.retired = false`
		}
	}
	cypher += ` RETURN mp ORDER BY mp.createdAt`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), nil
}

func (r *masterPlanReadRepository) GetMasterPlanMilestonesForMasterPlans(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanReadRepository.GetMasterPlanMilestonesForMasterPlans")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:MASTER_PLAN_BELONGS_TO_TENANT]-(mp:MasterPlan)-[:HAS_MILESTONE]->(m:MasterPlanMilestone)
		 WHERE mp.id IN $ids 
		 RETURN m, mp.id ORDER BY m.optional, m.order`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}
