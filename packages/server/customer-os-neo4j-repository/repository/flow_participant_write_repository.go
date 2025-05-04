package neo4j_repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type FlowParticipantWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *neo4j_entity.FlowParticipantEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, tx *neo4j.ManagedTransaction, id string) error

	MarkReadyContactsAsScheduled(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string) ([]string, error)
}

type flowParticipantWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowParticipantWriteRepository(driver *neo4j.DriverWithContext, database string) FlowParticipantWriteRepository {
	return &flowParticipantWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowParticipantWriteRepositoryImpl) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *neo4j_entity.FlowParticipantEntity) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowParticipantWriteRepository.Merge")
	defer spans.Finish()

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fc:FlowParticipant:FlowParticipant_%s {id: $id})
			ON MATCH SET
				fc.updatedAt = $updatedAt,
				fc.entityId = $entityId,
				fc.entityType = $entityType,
				fc.requirementsUnmeet = $requirementsUnmeet,
				fc.status = $status
			ON CREATE SET
				fc.createdAt = $createdAt,
				fc.updatedAt = $updatedAt,
				fc.entityId = $entityId,
				fc.entityType = $entityType,
				fc.requirementsUnmeet = $requirementsUnmeet,
				fc.status = $status
			RETURN fc`, common.GetTenantFromContext(ctx))

	params := map[string]any{
		"tenant":             common.GetTenantFromContext(ctx),
		"id":                 entity.Id,
		"createdAt":          utils.NowIfZero(entity.CreatedAt),
		"updatedAt":          utils.NowIfZero(entity.UpdatedAt),
		"entityId":           entity.EntityId,
		"entityType":         entity.EntityType.String(),
		"requirementsUnmeet": entity.RequirementsUnmeet,
		"status":             entity.Status,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	queryResult, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		qr, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, qr, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return queryResult.(*neo4j.Node), nil
}

func (r *flowParticipantWriteRepositoryImpl) Delete(ctx context.Context, tx *neo4j.ManagedTransaction, id string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.Delete")
	defer spans.Finish()

	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[r:BELONGS_TO_TENANT]-(fc:FlowParticipant_%s {id:$id}) delete r, fc`, tenant)

	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *flowParticipantWriteRepositoryImpl) MarkReadyContactsAsScheduled(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.MarkReadyContactsAsScheduled")
	defer spans.Finish()

	spans.LogKV("flowId", flowId)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fc:FlowParticipant_%s) where fc.status = 'READY' set fc.status = 'SCHEDULING' return fc.id`, tenant, tenant)

	params := map[string]any{
		"tenant": tenant,
		"flowId": flowId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		qr, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsString(ctx, qr, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return result.([]string), nil
}
