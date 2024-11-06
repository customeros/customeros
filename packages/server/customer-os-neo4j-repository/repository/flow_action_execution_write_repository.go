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

type FlowActionExecutionWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowActionExecutionEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, id string) error
}

type flowActionExecutionWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowActionExecutionWriteRepository(driver *neo4j.DriverWithContext, database string) FlowActionExecutionWriteRepository {
	return &flowActionExecutionWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowActionExecutionWriteRepositoryImpl) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowActionExecutionEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fa:FlowAction_%s {id: $actionId})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fae:FlowActionExecution:FlowActionExecution_%s {id: $id})
			ON MATCH SET
				fae.updatedAt = $updatedAt,
				fae.scheduledAt = $scheduledAt,
				fae.executedAt = $executedAt,
				fae.statusUpdatedAt = $statusUpdatedAt,
				fae.error = $error,
				fae.status = $status
				
			ON CREATE SET
				fae.createdAt = $createdAt,
				fae.updatedAt = $updatedAt,
				fae.flowId = $flowId,
				fae.entityId = $entityId,
				fae.entityType = $entityType,
				fae.actionId = $actionId,
				fae.scheduledAt = $scheduledAt,
				fae.statusUpdatedAt = $statusUpdatedAt,
				fae.status = $status,

				fae.mailbox = $mailbox,
				fae.userId = $userId,
				fae.socialUrl = $socialUrl
			
			WITH f, fa, fae
			MERGE (fa)-[:HAS_EXECUTION]->(fae)
			RETURN fae`, tenant, tenant, tenant)

	params := map[string]any{
		"tenant":          tenant,
		"id":              entity.Id,
		"createdAt":       utils.NowIfZero(entity.CreatedAt),
		"updatedAt":       utils.NowIfZero(entity.UpdatedAt),
		"flowId":          entity.FlowId,
		"entityId":        entity.EntityId,
		"entityType":      entity.EntityType,
		"actionId":        entity.ActionId,
		"scheduledAt":     entity.ScheduledAt,
		"executedAt":      utils.TimePtrAsAny(entity.ExecutedAt),
		"statusUpdatedAt": entity.StatusUpdatedAt,
		"status":          entity.Status,
		"mailbox":         entity.Mailbox,
		"userId":          entity.UserId,
		"socialUrl":       entity.SocialUrl,
		"error":           entity.Error,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	queryResult, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
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
}

func (r *flowActionExecutionWriteRepositoryImpl) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[r:BELONGS_TO_TENANT]-(fc:FlowParticipant_%s {id:$id}) delete r, fc`, tenant)

	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
