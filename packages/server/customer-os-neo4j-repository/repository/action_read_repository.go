package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ActionReadRepository interface {
	GetSingleAction(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType) (*dbtype.Node, error)
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

func (r *actionReadRepository) GetSingleAction(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionReadRepository.GetSingleAction")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("entityId", entityId),
		log.String("entityType", entityType.String()),
		log.String("actionType", string(actionType)))

	cypher := fmt.Sprintf(`MATCH  (n:%s_%s {id:$entityId}) `, entityType.Neo4jLabel(), tenant)
	cypher += `WITH n
			  MATCH (n)<-[:ACTION_ON]-(a:Action {type:$type})
			  return a limit 1`
	params := map[string]any{
		"entityId": entityId,
		"type":     actionType,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
