package neo4j_repository

import (
	"context"
	"fmt"

	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/constants"
	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type ActionWriteRepository interface {
	//Deprecated: Use CreateV2 instead
	Create(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string) (*dbtype.Node, error)
	//Deprecated: Use CreateV2 instead
	CreateWithProperties(ctx context.Context, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string, extraProperties map[string]any) (*dbtype.Node, error)
	MergeByActionType(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string) (*dbtype.Node, error)
	CreateV2(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, actionId, entityId string, entityType model.EntityType, data data_fields.ActionFields) error
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ActionWriteRepository.CreateWithProperties")
	defer spans.Finish()
	spans.LogKV("entityId", entityId,
		"entityType", entityType.String(),
		"actionType", string(actionType),
		"content", content,
		"metadata", metadata,
		"createdAt", createdAt,
		"extraProperties", extraProperties)

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id:$entityId}) `, entityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(` MERGE (n)<-[:ACTION_ON]-(a:Action {id:randomUUID()}) 
				ON CREATE SET 	a.type=$type, 
								a.content=$content,
								a.metadata=$metadata,
								a.createdAt=$createdAt, 
								a.updatedAt=datetime(),
								a.source=$source, 
								a.appSource=$appSource, 
								a:Action_%s, 
								a:TimelineEvent, 
								a:TimelineEvent_%s`, tenant, tenant)
	if len(extraProperties) > 0 {
		cypher += ` SET a += $extraProperties `
	}
	cypher += ` return a `

	params := map[string]any{
		"tenant":    tenant,
		"entityId":  entityId,
		"type":      actionType,
		"content":   content,
		"metadata":  metadata,
		"source":    constants.SourceOpenline,
		"appSource": appSource,
		"createdAt": createdAt,
	}
	if len(extraProperties) > 0 {
		params["extraProperties"] = extraProperties
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *actionWriteRepository) MergeByActionType(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, actionType enum.ActionType, content, metadata string, createdAt time.Time, appSource string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ActionWriteRepository.MergeByActionType")
	defer spans.Finish()
	spans.LogKV("entityId", entityId,
		"entityType", entityType.String(),
		"actionType", string(actionType),
		"content", content)

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
								a.appSource=$appSource, 
								a:Action_%s, 
								a:TimelineEvent, 
								a:TimelineEvent_%s)`, tenant, tenant)
	cypher += ` WITH n
				MATCH (n)<-[:ACTION_ON]-(act:Action {type:$type})
				RETURN act `
	params := map[string]any{
		"tenant":    tenant,
		"entityId":  entityId,
		"type":      actionType,
		"content":   content,
		"metadata":  metadata,
		"source":    constants.SourceOpenline,
		"appSource": neo4jmodel.GetAppSource(appSource),
		"createdAt": createdAt,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *actionWriteRepository) CreateV2(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, actionId, entityId string, entityType model.EntityType, data data_fields.ActionFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ActionWriteRepository.CreateV2")
	defer spans.Finish()

	spans.LogKV("tenant", tenant,
		"actionId", actionId,
		"entityId", entityId,
		"entityType", entityType.String())
	spans.LogObjectAsJson("data", data)

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id:$entityId}) `, entityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(` MERGE (n)<-[:ACTION_ON]-(a:Action {id:$actionId}) 
				ON CREATE SET 	a.type=$type, 
								a.content=$content,
								a.metadata=$metadata,
								a.createdAt=$createdAt, 
								a.updatedAt=datetime(),
								a.source=$source, 
								a.appSource=$appSource, 
								a:Action_%s, 
								a:TimelineEvent, 
								a:TimelineEvent_%s`, tenant, tenant)
	if len(data.ExtraProperties) > 0 {
		cypher += ` SET a += $extraProperties `
	}

	params := map[string]any{
		"tenant":    tenant,
		"actionId":  actionId,
		"entityId":  entityId,
		"content":   utils.IfNotNilString(data.Content),
		"metadata":  utils.IfNotNilString(data.Metadata),
		"source":    utils.IfNotNilString(data.Source),
		"appSource": utils.IfNotNilString(data.AppSource),
		"createdAt": utils.IfNotNilTimeWithDefault(data.CreatedAt, utils.Now()),
	}
	if data.ActionType != nil {
		params["type"] = string(*data.ActionType)
	} else {
		params["type"] = ""
	}
	if len(data.ExtraProperties) > 0 {
		params["extraProperties"] = data.ExtraProperties
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx, cypher, params)
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
