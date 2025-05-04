package neo4j_repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type FlowActionReadRepository interface {
	GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error)
	GetById(ctx context.Context, id string) (*neo4j.Node, error)
	GetStartAction(ctx context.Context, flowId string) (*neo4j.Node, error)
	GetPrevious(ctx context.Context, actionId string) ([]*neo4j.Node, error)
	GetNext(ctx context.Context, actionId string) ([]*neo4j.Node, error)
	GetFlowByActionId(ctx context.Context, id string) (*neo4j.Node, error)
	GetFlowByEntity(ctx context.Context, entityId string, entityType model.EntityType) (*neo4j.Node, error)
}

type flowActionReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowActionReadRepository(driver *neo4j.DriverWithContext, database string) FlowActionReadRepository {
	return &flowActionReadRepositoryImpl{driver: driver, database: database}
}

func (r flowActionReadRepositoryImpl) GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionReadRepository.GetList")
	defer spans.Finish()

	if len(flowIds) > 0 {
		spans.LogKV("flowIds", fmt.Sprintf("%v", flowIds))
	}

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant": tenant,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fa:FlowAction_%s) `, tenant, tenant)
	if len(flowIds) > 0 {
		cypher += "WHERE f.id in $flowIds "
		params["flowIds"] = flowIds
	}
	cypher += "RETURN fa, f.id ORDER by fa.index"

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	if len(result.([]*utils.DbNodeAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r flowActionReadRepositoryImpl) GetById(ctx context.Context, id string) (*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionReadRepository.GetById")
	defer spans.Finish()

	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fa:FlowAction_%s {id: $id}) RETURN fa`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionReadRepository.GetStartAction")
	defer spans.Finish()

	spans.LogKV("flowId", flowId)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fa:FlowAction_%s {action: 'FLOW_START'}) RETURN fa`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"flowId": flowId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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

func (r flowActionReadRepositoryImpl) GetPrevious(ctx context.Context, actionId string) ([]*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionReadRepository.GetPrevious")
	defer spans.Finish()

	spans.LogKV("actionId", actionId)

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant":   tenant,
		"actionId": actionId,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(previous:FlowAction_%s )-[:NEXT]->(current:FlowAction_%s {id: $actionId}) RETURN previous`, tenant, tenant)

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	if len(result.([]*dbtype.Node)) == 0 {
		return nil, nil
	}
	return result.([]*dbtype.Node), err
}

func (r flowActionReadRepositoryImpl) GetNext(ctx context.Context, actionId string) ([]*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionReadRepository.GetNext")
	defer spans.Finish()

	spans.LogKV("actionId", actionId)

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant":   tenant,
		"actionId": actionId,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(:FlowAction_%s {id: $actionId})-[:NEXT]->(fa:FlowAction_%s) RETURN fa`, tenant, tenant)

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	if len(result.([]*dbtype.Node)) == 0 {
		return nil, nil
	}
	return result.([]*dbtype.Node), err
}

func (r flowActionReadRepositoryImpl) GetFlowByActionId(ctx context.Context, id string) (*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionReadRepository.GetFlowByActionId")
	defer spans.Finish()

	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fa:FlowAction_%s {id: $id}) RETURN f`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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

func (r flowActionReadRepositoryImpl) GetFlowByEntity(ctx context.Context, entityId string, entityType model.EntityType) (*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowActionReadRepository.GetFlowByEntity")
	defer spans.Finish()

	spans.LogKV("entityId", entityId)
	spans.LogKV("entityType", entityType.String())

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowParticipant_%s { entityId: $entityId, entityType: $entityType }) RETURN f`, tenant, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"entityId":   entityId,
		"entityType": entityType.String(),
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
