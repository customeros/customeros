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

type FlowSequenceContactWriteRepository interface {
	Merge(ctx context.Context, entity *entity.FlowSequenceContactEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, id string) error
}

type flowSequenceContactWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowSequenceContactWriteRepository(driver *neo4j.DriverWithContext, database string) FlowSequenceContactWriteRepository {
	return &flowSequenceContactWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowSequenceContactWriteRepositoryImpl) Merge(ctx context.Context, entity *entity.FlowSequenceContactEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fs:FlowSequenceContact:FlowSequenceContact_%s {id: $id})
			ON MATCH SET
				fs.updatedAt = $updatedAt
			ON CREATE SET
				fs.createdAt = $createdAt,
				fs.updatedAt = $updatedAt
			RETURN fs`, common.GetTenantFromContext(ctx))

	params := map[string]any{
		"tenant":    common.GetTenantFromContext(ctx),
		"id":        entity.Id,
		"createdAt": utils.TimeOrNow(entity.CreatedAt),
		"updatedAt": utils.TimeOrNow(entity.UpdatedAt),
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

func (r *flowSequenceContactWriteRepositoryImpl) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[r:BELONGS_TO_TENANT]-(n:FlowSequenceContact_%s {id:$id}) delete r, n`, tenant)

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
