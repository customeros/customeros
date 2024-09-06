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

type FlowWriteRepository interface {
	Merge(ctx context.Context, entity *entity.FlowEntity) (*dbtype.Node, error)
}

type flowWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowWriteRepository(driver *neo4j.DriverWithContext, database string) FlowWriteRepository {
	return &flowWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowWriteRepositoryImpl) Merge(ctx context.Context, entity *entity.FlowEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(f:Flow:Flow_%s { id: $id })
			ON MATCH SET
				f.name = $name,
				f.description = $description,
				f.updatedAt = $updatedAt,
				f.status = $status
			ON CREATE SET
				f.createdAt = $createdAt,
				f.updatedAt = $updatedAt,
				f.name = $name,
				f.description = $description,
				f.status = $status
			RETURN f`, common.GetTenantFromContext(ctx))

	params := map[string]any{
		"tenant":      common.GetTenantFromContext(ctx),
		"id":          entity.Id,
		"name":        entity.Name,
		"description": entity.Description,
		"status":      entity.Status,
		"createdAt":   utils.TimeOrNow(entity.CreatedAt),
		"updatedAt":   utils.TimeOrNow(entity.UpdatedAt),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
