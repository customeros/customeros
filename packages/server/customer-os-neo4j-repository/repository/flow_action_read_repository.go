package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowActionReadRepository interface {
	GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error)
	GetById(ctx context.Context, id string) (*neo4j.Node, error)
	GetStartAction(ctx context.Context, flowId string) (*neo4j.Node, error)
	GetNext(ctx context.Context, actionId string) ([]*neo4j.Node, error)
	GetFlowByActionId(ctx context.Context, id string) (*neo4j.Node, error)
	GetFlowByContactId(ctx context.Context, id string) (*neo4j.Node, error)
}

type flowActionReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowActionReadRepository(driver *neo4j.DriverWithContext, database string) FlowActionReadRepository {
	return &flowActionReadRepositoryImpl{driver: driver, database: database}
}

func (r flowActionReadRepositoryImpl) GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionReadRepository.GetList")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if flowIds != nil && len(flowIds) > 0 {
		span.LogFields(log.String("flowIds", fmt.Sprintf("%v", flowIds)))
	}

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant": tenant,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fa:FlowAction_%s) `, tenant, tenant)
	if flowIds != nil && len(flowIds) > 0 {
		cypher += "WHERE f.id in $flowIds "
		params["flowIds"] = flowIds
	}
	cypher += "RETURN fa, f.id ORDER by fa.index"

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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
	if len(result.([]*utils.DbNodeAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r flowActionReadRepositoryImpl) GetById(ctx context.Context, id string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionReadRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fa:FlowAction_%s {id: $id}) RETURN fa`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r flowActionReadRepositoryImpl) GetStartAction(ctx context.Context, flowId string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionReadRepository.GetStartAction")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fa:FlowAction_%s {action: 'FLOW_START'}) RETURN fa`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"flowId": flowId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r flowActionReadRepositoryImpl) GetNext(ctx context.Context, actionId string) ([]*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionReadRepository.GetNext")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("actionId", actionId))

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant":   tenant,
		"actionId": actionId,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(:FlowAction_%s {id: $actionId})-[:NEXT]->(fa:FlowAction_%s) RETURN fa`, tenant, tenant)

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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
	if len(result.([]*dbtype.Node)) == 0 {
		return nil, nil
	}
	return result.([]*dbtype.Node), err
}

func (r flowActionReadRepositoryImpl) GetFlowByActionId(ctx context.Context, id string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionReadRepository.GetFlowByActionId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fa:FlowAction_%s {id: $id}) RETURN f`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r flowActionReadRepositoryImpl) GetFlowByContactId(ctx context.Context, id string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionReadRepository.GetFlowByContactId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowContact_%s {id: $id}) RETURN f`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
