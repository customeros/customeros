package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type ActionWriteRepository interface {
	Create(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string) (*dbtype.Node, error)
	CreateWithProperties(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string, extraProperties map[string]any) (*dbtype.Node, error)
	MergeByActionType(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string) (*dbtype.Node, error)
}

type actionWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewActionWriteRepository(driver *neo4j.DriverWithContext, database string) ActionWriteRepository {
	return &actionWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *actionWriteRepository) Create(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string) (*dbtype.Node, error) {
	return r.CreateWithProperties(ctx, tenant, entityId, entityType, actionType, content, metadata, createdAt, appSource, nil)
}

func (r *actionWriteRepository) CreateWithProperties(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string, extraProperties map[string]any) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionRepository.CreateWithProperties")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("entityId", entityId),
		log.String("entityType", entityType.String()),
		log.String("actionType", string(actionType)),
		log.String("content", content),
		log.String("metadata", metadata),
		log.Object("createdAt", createdAt),
		log.Object("extraProperties", extraProperties))

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id:$entityId}) `, entityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(` MERGE (n)<-[:ACTION_ON]-(a:Action {id:randomUUID()}) 
				ON CREATE SET 	a.type=$type, 
								a.content=$content,
								a.metadata=$metadata,
								a.createdAt=$createdAt, 
								a.updatedAt=datetime(),
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
		"appSource":     appSource,
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
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *actionWriteRepository) MergeByActionType(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ActionRepository.MergeByActionType")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("entityId", entityId),
		log.String("entityType", entityType.String()),
		log.String("actionType", string(actionType)),
		log.String("content", content))

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id:$entityId}) `, entityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`WITH n
								OPTIONAL MATCH (n)<-[:ACTION_ON]-(checkA:Action {type:$type})
								FOREACH (ignore IN CASE WHEN checkA IS NULL THEN [1] ELSE [] END |
								MERGE (n)<-[:ACTION_ON]-(a:Action {id:randomUUID()}) 
				ON CREATE SET 	a.type=$type,
								a.content=$content,
								a.metadata=$metadata,
								a.createdAt=$createdAt, 
								a.updatedAt=datetime(),
								a.source=$source, 
								a.sourceOfTruth=$sourceOfTruth, 
								a.appSource=$appSource, 
								a:Action_%s, 
								a:TimelineEvent, 
								a:TimelineEvent_%s)`, tenant, tenant)
	cypher += ` WITH n
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
		"appSource":     neo4jmodel.GetAppSource(appSource),
		"createdAt":     createdAt,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
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
