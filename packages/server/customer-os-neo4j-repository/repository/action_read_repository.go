package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ActionReadRepository interface {
	GetFor(ctx context.Context, tenant string, entityType model.EntityType, entityIds []string) ([]*utils.DbNodeAndId, error)
	GetSingleAction(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType) (*dbtype.Node, error)
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

func (r *actionReadRepository) GetFor(ctx context.Context, tenant string, entityType model.EntityType, entityIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AttachmentRepository.GetAttachmentsForXX")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	span.LogFields(log.String("entityType", entityType.String()), log.String("entityIds", fmt.Sprintf("%v", entityIds)))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	var query = "MATCH (n:%s_%s)<-[:ACTION_ON]-(a:Action_%s)"
	query += " WHERE n.id IN $entityIds "
	query += " RETURN a, n.id"

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, entityType.Neo4jLabel(), tenant, tenant),
			map[string]any{
				"tenant":    tenant,
				"entityIds": entityIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	span.LogFields(log.String("query", query))
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *actionReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *actionReadRepository) GetSingleAction(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionReadRepository.GetSingleAction")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
