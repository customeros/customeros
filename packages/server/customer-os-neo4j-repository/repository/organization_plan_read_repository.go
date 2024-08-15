package repository

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OrganizationPlanReadRepository interface {
	GetOrganizationPlanById(ctx context.Context, tenant, organizationPlanId string) (*dbtype.Node, error)
	GetOrganizationPlanMilestoneById(ctx context.Context, tenant, organizationPlanMilestoneId string) (*dbtype.Node, error)
	GetOrganizationPlanMilestoneByPlanAndId(ctx context.Context, tenant, organizationPlanId, organizationPlanMilestoneId string) (*dbtype.Node, error)
	GetOrganizationPlansOrderByCreatedAt(ctx context.Context, tenant string, returnRetired *bool) ([]*dbtype.Node, error)
	GetOrganizationPlanMilestonesForOrganizationPlans(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetMaxOrderForOrganizationPlanMilestones(ctx context.Context, tenant, organizationPlanId string) (int64, error)
	GetOrganizationPlansForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error)
	GetMilestoneDueDate(ctx context.Context, tenant, organizationPlanMilestoneId string) (time.Time, error)
	GetMilestonesForOrganizationPlan(ctx context.Context, tenant, organizationPlanId string) ([]*dbtype.Node, error)
	GetOrganizationFromOrganizationPlan(ctx context.Context, tenant, organizationPlanId string) (*dbtype.Node, error)
}

type organizationPlanReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOrganizationPlanReadRepository(driver *neo4j.DriverWithContext, database string) OrganizationPlanReadRepository {
	return &organizationPlanReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *organizationPlanReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *organizationPlanReadRepository) GetOrganizationPlanById(ctx context.Context, tenant, organizationPlanId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetOrganizationPlanById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan {id:$id}) RETURN op`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationPlanId,
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

func (r *organizationPlanReadRepository) GetOrganizationPlanMilestoneById(ctx context.Context, tenant, organizationPlanMilestoneId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetOrganizationPlanMilestoneById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanMilestoneId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(:OrganizationPlan)-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone {id:$id}) RETURN m`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationPlanMilestoneId,
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

func (r *organizationPlanReadRepository) GetOrganizationPlansOrderByCreatedAt(ctx context.Context, tenant string, returnRetired *bool) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetOrganizationPlansOrderByCreatedAt")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan)`
	if returnRetired != nil {
		if *returnRetired {
			cypher += ` WHERE op.retired = true`
		} else {
			cypher += ` WHERE op.retired IS NULL OR op.retired = false`
		}
	}
	cypher += ` RETURN op ORDER BY op.createdAt`
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

func (r *organizationPlanReadRepository) GetOrganizationPlanMilestonesForOrganizationPlans(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetOrganizationPlanMilestonesForOrganizationPlans")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan)-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone)
		 WHERE op.id IN $ids 
		 RETURN m, op.id ORDER BY m.optional, m.order`
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

func (r *organizationPlanReadRepository) GetMilestonesForOrganizationPlan(ctx context.Context, tenant, organizationPlanId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetMilestonesForOrgPlan")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(:OrganizationPlan {id:$id})-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone)
		 WHERE m.retired IS NULL OR m.retired = false
		 RETURN m ORDER BY m.optional, m.order`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationPlanId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), err
}

func (r *organizationPlanReadRepository) GetMaxOrderForOrganizationPlanMilestones(ctx context.Context, tenant, organizationPlanId string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetMaxOrderForOrganizationPlanMilestones")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(:OrganizationPlan {id:$id})-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone)
			WHERE m.retired IS NULL OR m.retired = false
		 	RETURN coalesce(max(m.order),-1)`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationPlanId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return int64(0), err
	}
	span.LogFields(log.Int64("result.maxOrder", result.(int64)))
	return result.(int64), err
}

func (r *organizationPlanReadRepository) GetOrganizationPlanMilestoneByPlanAndId(ctx context.Context, tenant, organizationPlanId, organizationPlanMilestoneId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetOrganizationPlanMilestoneByPlanAndId")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanMilestoneId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(:OrganizationPlan {id:$organizationPlanId})-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone {id:$id}) RETURN m`
	params := map[string]any{
		"tenant":             tenant,
		"id":                 organizationPlanMilestoneId,
		"organizationPlanId": organizationPlanId,
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

func (r *organizationPlanReadRepository) GetOrganizationPlansForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetOrganizationPlansForOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(op:OrganizationPlan)-[:ORGANIZATION_PLAN_BELONGS_TO_ORGANIZATION]->(o:Organization {id:$organizationId}) RETURN op`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
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

func (r *organizationPlanReadRepository) GetMilestoneDueDate(ctx context.Context, tenant, organizationPlanMilestoneId string) (time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetMilestoneDueDate")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanMilestoneId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(:OrganizationPlan)-[:HAS_MILESTONE]->(m:OrganizationPlanMilestone {id:$id}) RETURN m.dueDate`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationPlanMilestoneId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[time.Time](ctx, queryResult, err)
	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return time.Time{}, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(time.Time), nil
}

func (r *organizationPlanReadRepository) GetOrganizationFromOrganizationPlan(ctx context.Context, tenant, organizationPlanId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanReadRepository.GetOrganizationFromOrganizationPlan")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_PLAN_BELONGS_TO_TENANT]-(:OrganizationPlan {id:$id})-[:ORGANIZATION_PLAN_BELONGS_TO_ORGANIZATION]->(o:Organization) RETURN o`
	params := map[string]any{
		"tenant": tenant,
		"id":     organizationPlanId,
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
