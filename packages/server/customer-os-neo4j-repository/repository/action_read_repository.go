package neo4j_repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type ActionReadRepository interface {
	GetFor(ctx context.Context, tenant string, entityType model.EntityType, entityIds []string) ([]*utils.DbNodeAndId, error)
	GetLastAction(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType) (*dbtype.Node, error)
}

type actionReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewActionReadRepository(driver *neo4j.DriverWithContext, database string) ActionReadRepository {
	return &actionReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *actionReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *actionReadRepository) GetFor(ctx context.Context, tenant string, entityType model.EntityType, entityIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ActionReadRepository.GetFor")
	defer spans.Finish()
	spans.LogKV("entityType", entityType.String())
	spans.LogKV("entityIds", fmt.Sprintf("%v", entityIds))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf("MATCH (n:%s_%s)<-[:ACTION_ON]-(a:Action_%s) WHERE n.id IN $entityIds RETURN a, n.id",
		entityType.Neo4jLabel(), tenant, tenant)

	params := map[string]any{
		"tenant":    tenant,
		"entityIds": entityIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	return result.([]*utils.DbNodeAndId), err
}

func (r *actionReadRepository) GetLastAction(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ActionReadRepository.GetLastAction")
	defer spans.Finish()
	spans.TagEntity(entityId)
	spans.LogKV("entityType", entityType.String())
	spans.LogKV("actionType", string(actionType))

	cypher := fmt.Sprintf(`MATCH  (n:%s_%s {id:$entityId}) `, entityType.Neo4jLabel(), tenant)
	cypher += `WITH n
			  MATCH (n)<-[:ACTION_ON]-(a:Action {type:$type})
			  RETURN a ORDER BY a.createdAt DESC LIMIT 1`
	params := map[string]any{
		"entityId": entityId,
		"type":     actionType,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*dbtype.Node), nil
}
