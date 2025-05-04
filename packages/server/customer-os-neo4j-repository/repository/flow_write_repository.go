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

type FlowWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *neo4j_entity.FlowEntity) (*dbtype.Node, error)

	UpdateStatistics(ctx context.Context) ([]*utils.StringsWithTenant, error)
	UpdateFlowStatistics(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string) ([]*utils.StringsWithTenant, error)
}

type flowWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowWriteRepository(driver *neo4j.DriverWithContext, database string) FlowWriteRepository {
	return &flowWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowWriteRepositoryImpl) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *neo4j_entity.FlowEntity) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowWriteRepository.Merge")
	defer spans.Finish()

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(f:Flow:Flow_%s { id: $id })
			ON MATCH SET
				f.name = $name,
				f.updatedAt = $updatedAt,
				f.nodes = $nodes,
				f.edges = $edges,
				f.firstStartedAt = $firstStartedAt,
				f.defaultName = $defaultName,
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
				f.defaultName = $defaultName,
				f.name = $name,
				f.nodes = $nodes,
				f.edges = $edges,
				f.firstStartedAt = $firstStartedAt,
				f.status = $status,
				f.onHold = $onHold,
				f.ready = $ready,
				f.scheduled = $scheduled,
				f.inProgress = $inProgress,
				f.completed = $completed,
				f.goalAchieved = $goalAchieved
			RETURN f`, common.GetTenantFromContext(ctx))

	params := map[string]any{
		"tenant":         common.GetTenantFromContext(ctx),
		"id":             entity.Id,
		"defaultName":    entity.DefaultName,
		"name":           entity.Name,
		"nodes":          entity.Nodes,
		"edges":          entity.Edges,
		"firstStartedAt": utils.TimePtrAsAny(entity.FirstStartedAt),
		"status":         entity.Status,
		"createdAt":      utils.NowIfZero(entity.CreatedAt),
		"updatedAt":      utils.NowIfZero(entity.UpdatedAt),
		"onHold":         entity.OnHold,
		"ready":          entity.Ready,
		"scheduled":      entity.Scheduled,
		"inProgress":     entity.InProgress,
		"completed":      entity.Completed,
		"goalAchieved":   entity.GoalAchieved,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
			spans.TraceError(err)
			return nil, err
		}

		return queryResult.(*neo4j.Node), nil
	} else {
		queryResult, err := (*tx).Run(ctx, cypher, params)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *flowWriteRepositoryImpl) UpdateStatistics(ctx context.Context) ([]*utils.StringsWithTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowWriteRepository.UpdateStatistics")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant)<-[:BELONGS_TO_TENANT]-(f:Flow)-[:HAS]->(fp:FlowParticipant)
WITH t, f, fp.status AS flowStatus, COUNT(fp.status) AS fs
WITH t, f, 
    CASE flowStatus
        WHEN 'ON_HOLD' THEN 'onHold'
        WHEN 'READY' THEN 'ready'
        WHEN 'SCHEDULED' THEN 'scheduled'
        WHEN 'IN_PROGRESS' THEN 'inProgress'
        WHEN 'COMPLETED' THEN 'completed'
        WHEN 'GOAL_ACHIEVED' THEN 'goalAchieved'
        ELSE null
    END AS property, fs
WHERE property IS NOT NULL
WITH t, f, collect({property: property, count: fs}) AS updates
UNWIND updates AS update
WITH t, f, update.property AS property, update.count AS count
SET f[property] = count
WITH t, f, collect(property) AS updatedProperties
WITH t, f, updatedProperties, 
    ['onHold', 'ready', 'scheduled', 'inProgress', 'completed', 'goalAchieved'] AS allProperties
WITH t, f, [prop IN allProperties WHERE NOT prop IN updatedProperties] AS propsToReset
UNWIND propsToReset AS prop
SET f[prop] = 0

RETURN collect(f.id), t.name`

	params := map[string]any{}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		r, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsStringsWithTenant(ctx, r, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	flowsUpdated := result.([]*utils.StringsWithTenant)

	for _, flowUpdated := range flowsUpdated {
		spans.LogKV("flowsUpdated."+flowUpdated.Tenant, fmt.Sprintf("%v", flowUpdated.Strings))
	}

	return flowsUpdated, nil
}

func (r *flowWriteRepositoryImpl) UpdateFlowStatistics(ctx context.Context, tx *neo4j.ManagedTransaction, flowId string) ([]*utils.StringsWithTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "FlowWriteRepository.UpdateFlowStatistics")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant{name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fp:FlowParticipant_%s)
WITH t, f, fp.status AS flowParticipantStatus, COUNT(fp.status) AS fs
WITH t, f, 
    CASE flowParticipantStatus
        WHEN 'ON_HOLD' THEN 'onHold'
        WHEN 'READY' THEN 'ready'
        WHEN 'SCHEDULED' THEN 'scheduled'
        WHEN 'IN_PROGRESS' THEN 'inProgress'
        WHEN 'COMPLETED' THEN 'completed'
        WHEN 'GOAL_ACHIEVED' THEN 'goalAchieved'
        ELSE null
    END AS property, fs
WHERE property IS NOT NULL
WITH t, f, collect({property: property, count: fs}) AS updates
UNWIND updates AS update
WITH t, f, update.property AS property, update.count AS count
SET f[property] = count
WITH t, f, collect(property) AS updatedProperties
WITH t, f, updatedProperties, 
    ['onHold', 'ready', 'scheduled', 'inProgress', 'completed', 'goalAchieved'] AS allProperties
WITH t, f, [prop IN allProperties WHERE NOT prop IN updatedProperties] AS propsToReset
UNWIND propsToReset AS prop
SET f[prop] = 0
RETURN collect(f.id), t.name`, tenant, tenant)

	params := map[string]any{
		"tenant": tenant,
		"flowId": flowId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		r, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsStringsWithTenant(ctx, r, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	flowsUpdated := result.([]*utils.StringsWithTenant)

	for _, flowUpdated := range flowsUpdated {
		spans.LogKV("flowsUpdated."+flowUpdated.Tenant, fmt.Sprintf("%v", flowUpdated.Strings))
	}

	return flowsUpdated, nil
}
