package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type FlowSenderReadRepository interface {
	GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error)
	GetById(ctx context.Context, id string) (*neo4j.Node, error)
	GetFlowBySenderId(ctx context.Context, id string) (*neo4j.Node, error)
	GetUsersUsedInFlows(ctx context.Context, userIds []string) ([]string, error)
}

type flowSenderReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowSenderReadRepository(driver *neo4j.DriverWithContext, database string) FlowSenderReadRepository {
	return &flowSenderReadRepositoryImpl{driver: driver, database: database}
}

func (r flowSenderReadRepositoryImpl) GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowSenderReadRepository.GetList")
	defer spans.Finish()

	if len(flowIds) > 0 {
		spans.LogKV("flowIds", fmt.Sprintf("%v", flowIds))
	}

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant": tenant,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fs:FlowSender_%s) `, tenant, tenant)
	if len(flowIds) > 0 {
		cypher += "WHERE f.id in $flowIds "
		params["flowIds"] = flowIds
	}
	cypher += "RETURN fs, f.id"

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

func (r flowSenderReadRepositoryImpl) GetById(ctx context.Context, id string) (*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowSenderReadRepository.GetById")
	defer spans.Finish()

	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fs:FlowSender_%s {id: $id}) RETURN fs`, tenant, tenant)
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

func (r flowSenderReadRepositoryImpl) GetFlowBySenderId(ctx context.Context, id string) (*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowSenderReadRepository.GetFlowBySenderId")
	defer spans.Finish()

	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fs:FlowSender_%s {id: $id}) RETURN f`, tenant, tenant)
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

func (r flowSenderReadRepositoryImpl) GetUsersUsedInFlows(ctx context.Context, userIds []string) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowSenderReadRepository.GetUsersUsedInFlows")
	defer spans.Finish()

	spans.LogObjectAsJson("userIds", userIds)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fs:FlowSender_%s) where not f.status = 'ARCHIVED' and fs.userId in $userIds RETURN DISTINCT fs.userId AS userId`, tenant, tenant)
	params := map[string]any{
		"tenant":  tenant,
		"userIds": userIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.([]string), nil
}
