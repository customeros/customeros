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

type FlowParticipantWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowParticipantEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, tx *neo4j.ManagedTransaction, id string) error
}

type flowParticipantWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowParticipantWriteRepository(driver *neo4j.DriverWithContext, database string) FlowParticipantWriteRepository {
	return &flowParticipantWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowParticipantWriteRepositoryImpl) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowParticipantEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fc:FlowParticipant:FlowParticipant_%s {id: $id})
			ON MATCH SET
				fc.updatedAt = $updatedAt,
				fc.entityId = $entityId,
				fc.entityType = $entityType,
				fc.status = $status
			ON CREATE SET
				fc.createdAt = $createdAt,
				fc.updatedAt = $updatedAt,
				fc.entityId = $entityId,
				fc.entityType = $entityType,
				fc.status = $status
			RETURN fc`, common.GetTenantFromContext(ctx))

	params := map[string]any{
		"tenant":     common.GetTenantFromContext(ctx),
		"id":         entity.Id,
		"createdAt":  utils.TimeOrNow(entity.CreatedAt),
		"updatedAt":  utils.TimeOrNow(entity.UpdatedAt),
		"entityId":   entity.EntityId,
		"entityType": entity.EntityType.String(),
		"status":     entity.Status,
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

func (r *flowParticipantWriteRepositoryImpl) Delete(ctx context.Context, tx *neo4j.ManagedTransaction, id string) error {
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

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
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
