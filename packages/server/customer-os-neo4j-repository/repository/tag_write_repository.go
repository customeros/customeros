package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TagWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, tag neo4jentity.TagEntity) (*dbtype.Node, error)
	LinkTagByIdToEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, tagId, linkedEntityId string, entityType model.EntityType) error
	UnlinkTagByIdFromEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, tagId, entityId string, entityType model.EntityType) error
	UnlinkAllAndDelete(ctx context.Context, tenant, tagId string) error
	UpdateName(ctx context.Context, tenant, tagId, name string) error
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:$name}) 
		 ON CREATE SET 
		  tag.id=randomUUID(),
		  tag.createdAt=$now,
		  tag.updatedAt=datetime(),
		  tag.source=$source,
		  tag:Tag_%s
		 RETURN tag`, tenant)
	params := map[string]any{
		"tenant": tenant,
		"name":   tag.Name,
		"source": tag.Source,
		"now":    utils.Now(),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result.(*dbtype.Node), nil
}

func (r *tagWriteRepository) LinkTagByIdToEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, tagId, entityId string, entityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagWriteRepository.LinkTagByIdToEntity")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("tagId", tagId), log.String("entityId", entityId), log.String("entityType", entityType.String()))

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

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
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

func (r *tagWriteRepository) UnlinkTagByIdFromEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, tagId, entityId string, entityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagWriteRepository.UnlinkTagByIdFromEntity")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("tagId", tagId), log.String("entityId", entityId), log.String("entityType", entityType.String()))

	cypher := fmt.Sprintf(`
		MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})<-[rel:TAGGED]-(e:%s {id:$entityId})
		DELETE rel
		SET e.updatedAt=datetime()`, entityType.Neo4jLabel()+"_"+tenant)
	params := map[string]any{
		"tenant":   tenant,
		"id":       tagId,
		"entityId": entityId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
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

func (r *tagWriteRepository) UnlinkAllAndDelete(ctx context.Context, tenant, tagId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagWriteRepository.UnlinkAllAndDelete")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("tagId", tagId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$tagId}) DETACH DELETE tag`
	params := map[string]any{
		"tenant": tenant,
		"tagId":  tagId,
	}

	return LogAndExecuteWriteQuery(ctx, *r.driver, cypher, params, span)
}

func (r *tagWriteRepository) UpdateName(ctx context.Context, tenant, tagId, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagWriteRepository.UpdateName")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("tagId", tagId), log.String("name", name))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$tagId}) 
				SET tag.name=$name, tag.updatedAt=datetime()`
	params := map[string]any{
		"tenant": tenant,
		"tagId":  tagId,
		"name":   name,
	}

	return LogAndExecuteWriteQuery(ctx, *r.driver, cypher, params, span)
}
