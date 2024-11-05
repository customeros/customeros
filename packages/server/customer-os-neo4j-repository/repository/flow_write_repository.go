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
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowEntity) (*dbtype.Node, error)
}

type flowWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowWriteRepository(driver *neo4j.DriverWithContext, database string) FlowWriteRepository {
	return &flowWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowWriteRepositoryImpl) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(f:Flow:Flow_%s { id: $id })
			ON MATCH SET
				f.name = $name,
				f.updatedAt = $updatedAt,
				f.nodes = $nodes,
				f.edges = $edges,
				f.status = $status,
				f.onHold = $onHold,
				f.ready = $ready,
				f.scheduled = $scheduled,
				f.inProgress = $inProgress,
				f.completed = $completed,
				f.goalAchieved = $goalAchieved
			ON CREATE SET
				f.createdAt = $createdAt,
				f.updatedAt = $updatedAt,
				f.name = $name,
				f.nodes = $nodes,
				f.edges = $edges,
				f.status = $status,
				f.onHold = $onHold,
				f.ready = $ready,
				f.scheduled = $scheduled,
				f.inProgress = $inProgress,
				f.completed = $completed,
				f.goalAchieved = $goalAchieved
			RETURN f`, common.GetTenantFromContext(ctx))

	params := map[string]any{
		"tenant":       common.GetTenantFromContext(ctx),
		"id":           entity.Id,
		"name":         entity.Name,
		"nodes":        entity.Nodes,
		"edges":        entity.Edges,
		"status":       entity.Status,
		"createdAt":    utils.NowIfZero(entity.CreatedAt),
		"updatedAt":    utils.NowIfZero(entity.UpdatedAt),
		"onHold":       entity.OnHold,
		"ready":        entity.Ready,
		"scheduled":    entity.Scheduled,
		"inProgress":   entity.InProgress,
		"completed":    entity.Completed,
		"goalAchieved": entity.GoalAchieved,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

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
