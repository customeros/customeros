package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type TagWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, tag neo4jentity.TagEntity) (*dbtype.Node, error)
	LinkTagByIdToEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, tagId, linkedEntityId string, entityType model.EntityType) error
	UnlinkTagByIdFromEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, tagId, entityId string, entityType model.EntityType) error
	UnlinkAllAndDelete(ctx context.Context, tenant, tagId string) error
	Update(ctx context.Context, tenant, tagId string, name, colorCode *string) error
}

type tagWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewTagWriteRepository(driver *neo4j.DriverWithContext, database string) TagWriteRepository {
	return &tagWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *tagWriteRepository) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, tag neo4jentity.TagEntity) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagWriteRepository.Merge")
	defer spans.Finish()

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:$name, entityType:$entityType}) 
		 ON CREATE SET 
		  tag.id=randomUUID(),
		  tag.createdAt=datetime(),
		  tag.updatedAt=datetime(),
		  tag.source=$source,
		  tag.appSource=$appSource,
		  tag.colorCode=$colorCode,
		  tag:Tag_%s
		 RETURN tag`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"name":       tag.Name,
		"source":     tag.Source,
		"appSource":  tag.AppSource,
		"entityType": tag.EntityType.String(),
		"colorCode":  tag.ColorCode,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return result.(*dbtype.Node), nil
}

func (r *tagWriteRepository) LinkTagByIdToEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, tagId, entityId string, entityType model.EntityType) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagWriteRepository.LinkTagByIdToEntity")
	defer spans.Finish()

	spans.LogKV("tagId", tagId)
	spans.LogKV("entityId", entityId)
	spans.LogKV("entityType", entityType.String())

	cypher := fmt.Sprintf(`
		MATCH (e:%s {id:$entityId}),
			(t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
		MERGE (e)-[rel:TAGGED]->(tag)
		ON CREATE SET
			rel.taggedAt=$taggedAt,
			e.updatedAt=datetime()`, entityType.Neo4jLabel()+"_"+tenant)
	params := map[string]any{
		"tenant":   tenant,
		"id":       tagId,
		"taggedAt": utils.Now(),
		"entityId": entityId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *tagWriteRepository) UnlinkTagByIdFromEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, tagId, entityId string, entityType model.EntityType) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagWriteRepository.UnlinkTagByIdFromEntity")
	defer spans.Finish()

	spans.LogKV("tagId", tagId)
	spans.LogKV("entityId", entityId)
	spans.LogKV("entityType", entityType.String())

	cypher := fmt.Sprintf(`
		MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})<-[rel:TAGGED]-(e:%s {id:$entityId})
		DELETE rel
		SET e.updatedAt=datetime()`, entityType.Neo4jLabel()+"_"+tenant)
	params := map[string]any{
		"tenant":   tenant,
		"id":       tagId,
		"entityId": entityId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *tagWriteRepository) UnlinkAllAndDelete(ctx context.Context, tenant, tagId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagWriteRepository.UnlinkAllAndDelete")
	defer spans.Finish()

	spans.LogKV("tagId", tagId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$tagId}) DETACH DELETE tag`
	params := map[string]any{
		"tenant": tenant,
		"tagId":  tagId,
	}

	return LogAndExecuteWriteQuery(ctx, *r.driver, cypher, params, spans)
}

func (r *tagWriteRepository) Update(ctx context.Context, tenant, tagId string, name, colorCode *string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TagWriteRepository.UpdateName")
	defer spans.Finish()

	spans.TagEntity(tagId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$tagId})
				SET tag.updatedAt=datetime()`
	params := map[string]any{
		"tenant": tenant,
		"tagId":  tagId,
	}
	if name != nil {
		cypher += `, tag.name=$name`
		params["name"] = *name
	}
	if colorCode != nil {
		cypher += `, tag.colorCode=$colorCode`
		params["colorCode"] = *colorCode
	}

	return LogAndExecuteWriteQuery(ctx, *r.driver, cypher, params, spans)
}
