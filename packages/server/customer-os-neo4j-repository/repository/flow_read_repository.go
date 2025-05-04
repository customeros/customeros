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

type FlowReadRepository interface {
	GetList(ctx context.Context) ([]*dbtype.Node, error)
	GetWithParticipant(ctx context.Context, tx *neo4j.ManagedTransaction, flowParticipantId string) (*dbtype.Node, error)
	GetFlowsForParticipants(ctx context.Context, entityIds []string, entityType model.EntityType) ([]*utils.DbNodeAndId, error)
	GetListWithSender(ctx context.Context, senderIds []string) ([]*utils.DbNodeAndId, error)
	GetById(ctx context.Context, id string) (*dbtype.Node, error)
}

type flowReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowReadRepository(driver *neo4j.DriverWithContext, database string) FlowReadRepository {
	return &flowReadRepositoryImpl{driver: driver, database: database}
}

func (r flowReadRepositoryImpl) GetList(ctx context.Context) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowReadRepository.GetList")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s) where not f.status = 'ARCHIVED' RETURN f`, common.GetTenantFromContext(ctx))
	params := map[string]any{
		"tenant": common.GetTenantFromContext(ctx),
	}

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
	return result.([]*dbtype.Node), nil
}

func (r flowReadRepositoryImpl) GetWithParticipant(ctx context.Context, tx *neo4j.ManagedTransaction, flowParticipantId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowReadRepository.GetWithParticipant")
	defer spans.Finish()

	spans.LogObjectAsJson("flowParticipantId", flowParticipantId)

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant":        tenant,
		"participantId": flowParticipantId,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowParticipant_%s { id: $participantId }) RETURN f`, tenant, tenant)

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
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

func (r flowReadRepositoryImpl) GetFlowsForParticipants(ctx context.Context, entityIds []string, entityType model.EntityType) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowReadRepository.GetFlowsForParticipants")
	defer spans.Finish()

	spans.LogObjectAsJson("entityIds", entityIds)
	spans.LogObjectAsJson("entityType", entityType)

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant":     tenant,
		"entityIds":  entityIds,
		"entityType": entityType.String(),
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowParticipant_%s) where fc.entityType = $entityType and fc.entityId in $entityIds RETURN f, fc.entityId`, tenant, tenant)

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
	return result.([]*utils.DbNodeAndId), err
}

func (r flowReadRepositoryImpl) GetListWithSender(ctx context.Context, senderIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowReadRepository.GetListWithSender")
	defer spans.Finish()

	spans.LogObjectAsJson("senderIds", senderIds)

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant":    tenant,
		"senderIds": senderIds,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fs:FlowSender_%s) `, tenant, tenant)
	cypher += "where fs.id in $senderIds RETURN f, fs.id"

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
	return result.([]*utils.DbNodeAndId), err
}

func (r flowReadRepositoryImpl) GetById(ctx context.Context, id string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowReadRepository.GetById")
	defer spans.Finish()

	spans.LogObjectAsJson("id", id)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $id}) RETURN f`, common.GetTenantFromContext(ctx))
	params := map[string]any{
		"tenant": common.GetTenantFromContext(ctx),
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
