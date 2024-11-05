package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowExecutionSettingsWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowExecutionSettingsEntity) (*dbtype.Node, error)
}

type flowExecutionSettingsWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowExecutionSettingsWriteRepository(driver *neo4j.DriverWithContext, database string) FlowExecutionSettingsWriteRepository {
	return &flowExecutionSettingsWriteRepository{driver: driver, database: database}
}

func (r *flowExecutionSettingsWriteRepository) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowExecutionSettingsEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowExecutionSettingsWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fes:FlowExecutionSettings:FlowExecutionSettings_%s {id: $id})<-[:HAS]-(f)
			ON MATCH SET
				fes.updatedAt = $updatedAt
			ON CREATE SET
				fes.createdAt = $createdAt,
				fes.updatedAt = $updatedAt,
				fes.flowId = $flowId,
				fes.entityId = $entityId,
				fes.entityType = $entityType,
				fes.mailbox = $mailbox,
				fes.userId = $userId
			RETURN fes`, tenant, tenant)

	params := map[string]any{
		"tenant":    tenant,
		"id":        entity.Id,
		"createdAt": utils.NowIfZero(entity.CreatedAt),
		"updatedAt": utils.NowIfZero(entity.UpdatedAt),

		"flowId": entity.FlowId,

		"entityId":   entity.EntityId,
		"entityType": entity.EntityType,

		"mailbox": entity.Mailbox,
		"userId":  entity.UserId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if tx == nil {
		session := utils.NewNeo4jWriteSession(ctx, *r.driver)
		defer session.Close(ctx)

		queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			qr, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return nil, err
			}
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, qr, err)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return queryResult.(*neo4j.Node), nil
	} else {
		queryResult, err := (*tx).Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}
