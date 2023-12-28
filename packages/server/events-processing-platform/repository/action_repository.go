package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ActionRepository interface {
	Create(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType, content, metadata string, createdAt time.Time) (*dbtype.Node, error)
	CreateWithProperties(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType, content, metadata string, createdAt time.Time, extraProperties map[string]any) (*dbtype.Node, error)
	GetSingleAction(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType) (*dbtype.Node, error)
	MergeByActionType(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType, content, metadata string, createdAt time.Time) (*dbtype.Node, error)
}

type actionRepository struct {
	driver *neo4j.DriverWithContext
}

func NewActionRepository(driver *neo4j.DriverWithContext) ActionRepository {
	return &actionRepository{
		driver: driver,
	}
}

func (r *actionRepository) Create(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType, content, metadata string, createdAt time.Time) (*dbtype.Node, error) {
	return r.CreateWithProperties(ctx, tenant, entityId, entityType, actionType, content, metadata, createdAt, nil)
}

func (r *actionRepository) CreateWithProperties(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType, content, metadata string, createdAt time.Time, extraProperties map[string]any) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionRepository.CreateWithProperties")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("entityId", entityId),
		log.String("entityType", entityType.String()),
		log.String("actionType", string(actionType)),
		log.String("content", content),
		log.String("metadata", metadata),
		log.Object("createdAt", createdAt),
		log.Object("extraProperties", extraProperties))

	cypher := ""
	switch entityType {
	case entity.ORGANIZATION:
		cypher = fmt.Sprintf(`MATCH (n:Organization_%s {id:$entityId}) `, tenant)
	case entity.CONTRACT:
		cypher = fmt.Sprintf(`MATCH (n:Contract_%s {id:$entityId}) `, tenant)
	}

	cypher += fmt.Sprintf(` MERGE (n)<-[:ACTION_ON]-(a:Action {id:randomUUID()}) 
				ON CREATE SET 	a.type=$type, 
								a.content=$content,
								a.metadata=$metadata,
								a.createdAt=$createdAt, 
								a.source=$source, 
								a.sourceOfTruth=$sourceOfTruth,
								a.appSource=$appSource, 
								a:Action_%s, 
								a:TimelineEvent, 
								a:TimelineEvent_%s`, tenant, tenant)
	if extraProperties != nil && len(extraProperties) > 0 {
		cypher += ` SET a += $extraProperties `
	}
	cypher += ` return a `

	params := map[string]any{
		"tenant":        tenant,
		"entityId":      entityId,
		"type":          actionType,
		"content":       content,
		"metadata":      metadata,
		"source":        constants.SourceOpenline,
		"sourceOfTruth": constants.SourceOpenline,
		"appSource":     constants.AppSourceEventProcessingPlatform,
		"createdAt":     createdAt,
	}
	if extraProperties != nil && len(extraProperties) > 0 {
		params["extraProperties"] = extraProperties
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
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

func (r *actionRepository) GetSingleAction(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionRepository.GetSingleAction")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("entityId", entityId),
		log.String("entityType", entityType.String()),
		log.String("actionType", string(actionType)))

	query := ""
	switch entityType {
	case entity.ORGANIZATION:
		query = fmt.Sprintf(`MATCH  (n:Organization_%s {id:$entityId}) `, tenant)
	case entity.CONTRACT:
		query = fmt.Sprintf(`MATCH  (n:Contract_%s {id:$entityId}) `, tenant)
	}

	query += `WITH n
			  MATCH (n)<-[:ACTION_ON]-(a:Action {type:$type})
			  return a limit 1`

	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"entityId": entityId,
				"type":     actionType,
			}); err != nil {
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

func (r *actionRepository) MergeByActionType(ctx context.Context, tenant, entityId string, entityType entity.EntityType, actionType entity.ActionType, content, metadata string, createdAt time.Time) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionRepository.MergeByActionType")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("entityId", entityId),
		log.String("entityType", entityType.String()),
		log.String("actionType", string(actionType)),
		log.String("content", content))

	cypher := ""
	switch entityType {
	case entity.ORGANIZATION:
		cypher = fmt.Sprintf(`MATCH  (n:Organization_%s {id:$entityId}) `, tenant)
	case entity.CONTRACT:
		cypher = fmt.Sprintf(`MATCH  (n:Contract_%s {id:$entityId}) `, tenant)
	}

	cypher += fmt.Sprintf(`WITH n
								OPTIONAL MATCH (n)<-[:ACTION_ON]-(checkA:Action {type:$type})
								FOREACH (ignore IN CASE WHEN checkA IS NULL THEN [1] ELSE [] END |
								MERGE (n)<-[:ACTION_ON]-(a:Action {id:randomUUID()}) 
				ON CREATE SET 	a.type=$type,
								a.content=$content,
								a.metadata=$metadata,
								a.createdAt=$createdAt, 
								a.source=$source, 
								a.sourceOfTruth=$sourceOfTruth, 
								a.appSource=$appSource, 
								a:Action_%s, 
								a:TimelineEvent, 
								a:TimelineEvent_%s)`, tenant, tenant)
	cypher += ` 	WITH n
				MATCH (n)<-[:ACTION_ON]-(act:Action {type:$type})
				RETURN act `
	params := map[string]any{
		"tenant":        tenant,
		"entityId":      entityId,
		"type":          actionType,
		"content":       content,
		"metadata":      metadata,
		"source":        constants.SourceOpenline,
		"sourceOfTruth": constants.SourceOpenline,
		"appSource":     constants.AppSourceEventProcessingPlatform,
		"createdAt":     createdAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
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
