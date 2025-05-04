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

type FlowSenderWriteRepository interface {
	Merge(ctx context.Context, entity *neo4j_entity.FlowSenderEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, id string) error
}

type flowSenderWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowSenderWriteRepository(driver *neo4j.DriverWithContext, database string) FlowSenderWriteRepository {
	return &flowSenderWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowSenderWriteRepositoryImpl) Merge(ctx context.Context, entity *neo4j_entity.FlowSenderEntity) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowSenderWriteRepository.Merge")
	defer spans.Finish()

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fs:FlowSender:FlowSender_%s {id: $id})
			ON CREATE SET
				fs.createdAt = $createdAt,
				fs.updatedAt = $updatedAt,
				fs.userId = $userId
			RETURN fs`, common.GetTenantFromContext(ctx))

	params := map[string]any{
		"tenant":    common.GetTenantFromContext(ctx),
		"id":        entity.Id,
		"userId":    entity.UserId,
		"createdAt": utils.NowIfZero(entity.CreatedAt),
		"updatedAt": utils.NowIfZero(entity.UpdatedAt),
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}

	return result.(*dbtype.Node), nil
}

func (r *flowSenderWriteRepositoryImpl) Delete(ctx context.Context, id string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowSenderWriteRepository.Delete")
	defer spans.Finish()

	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[r:BELONGS_TO_TENANT]-(fs:FlowSender_%s {id:$id}) delete r, fs`, tenant)

	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return err
	}

	return nil
}
