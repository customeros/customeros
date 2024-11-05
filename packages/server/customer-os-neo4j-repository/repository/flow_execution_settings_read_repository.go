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

type FlowExecutionSettingsReadRepository interface {
	GetForEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId, entityType string) (*dbtype.Node, error)
}

type flowExecutionSettingsReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowExecutionSettingsReadRepository(driver *neo4j.DriverWithContext, database string) FlowExecutionSettingsReadRepository {
	return &flowExecutionSettingsReadRepositoryImpl{driver: driver, database: database}
}

func (r flowExecutionSettingsReadRepositoryImpl) GetForEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId, entityType string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowReadRepository.GetForEntity")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId), log.String("entityId", entityId))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fes:FlowExecutionSettings_%s {entityId: $entityId, entityType: $entityType}) RETURN fes limit 1`, tenant, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"flowId":     flowId,
		"entityId":   entityId,
		"entityType": entityType,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result.(*dbtype.Node), nil
}
