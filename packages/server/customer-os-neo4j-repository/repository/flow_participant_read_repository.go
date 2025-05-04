package neo4j_repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type FlowParticipantReadRepository interface {
	GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error)
	CountWithStatus(ctx context.Context, flowId string, status neo4j_entity.FlowParticipantStatus) (int64, error)
	Identify(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*neo4j.Node, error)
	GetById(ctx context.Context, id string) (*neo4j.Node, error)
}

type flowParticipantReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowParticipantReadRepository(driver *neo4j.DriverWithContext, database string) FlowParticipantReadRepository {
	return &flowParticipantReadRepositoryImpl{driver: driver, database: database}
}

func (r flowParticipantReadRepositoryImpl) GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowParticipantReadRepository.GetList")
	defer spans.Finish()

	if len(flowIds) > 0 {
		spans.LogKV("flowIds", fmt.Sprintf("%v", flowIds))
	}

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant": tenant,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowParticipant_%s) `, tenant, tenant)
	if len(flowIds) > 0 {
		cypher += "WHERE f.id in $flowIds "
		params["flowIds"] = flowIds
	}
	cypher += "RETURN fc, f.id"

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

func (r flowParticipantReadRepositoryImpl) CountWithStatus(ctx context.Context, flowId string, status neo4j_entity.FlowParticipantStatus) (int64, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowParticipantReadRepository.CountWithStatus")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fc:FlowParticipant_%s {status: $status}) return count(fc)`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"flowId": flowId,
		"status": status,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Single(ctx)
		}
	})
	if err != nil {
		return 0, err
	}
	organizationsCount := dbRecord.(*db.Record).Values[0].(int64)
	spans.LogKV("result", organizationsCount)
	return organizationsCount, nil
}

func (r flowParticipantReadRepositoryImpl) Identify(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowParticipantReadRepository.Identify")
	defer spans.Finish()

	spans.LogKV("flowId", flowId)
	spans.LogKV("entityId", entityId)
	spans.LogKV("entityType", entityType.String())

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s { id: $flowId})-[:HAS]->(fc:FlowParticipant_%s { entityId: $entityId, entityType: $entityType}) RETURN fc`, tenant, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"flowId":     flowId,
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

func (r flowParticipantReadRepositoryImpl) GetById(ctx context.Context, id string) (*neo4j.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowParticipantReadRepository.GetById")
	defer spans.Finish()

	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowParticipant_%s {id: $id}) RETURN fc`, tenant, tenant)
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
