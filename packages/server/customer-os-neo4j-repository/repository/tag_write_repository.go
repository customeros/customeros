package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/customer-os-neo4j-repository/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type TagWriteRepository interface {
	LinkTagByIdToEntity(ctx context.Context, tenant, tagId, linkedEntityId, linkedEntityNodeLabel string, taggedAt time.Time) error
	UnlinkTagByIdFromEntity(ctx context.Context, tenant, tagId, linkedEntityId, linkedEntityNodeLabel string) error
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

func (r *tagWriteRepository) LinkTagByIdToEntity(ctx context.Context, tenant, tagId, linkedEntityId, linkedEntityNodeLabel string, taggedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagWriteRepository.LinkTagByIdToEntity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("tagId", tagId), log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel), log.Object("taggedAt", taggedAt))

	cypher := fmt.Sprintf(`
		MATCH (e:%s {id:$entityId}),
			(t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
		MERGE (e)-[rel:TAGGED]->(tag)
		ON CREATE SET
			rel.taggedAt=$taggedAt `, linkedEntityNodeLabel+"_"+tenant)
	params := map[string]any{
		"tenant":   tenant,
		"id":       tagId,
		"taggedAt": taggedAt,
		"entityId": linkedEntityId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *tagWriteRepository) UnlinkTagByIdFromEntity(ctx context.Context, tenant, tagId, linkedEntityId, linkedEntityNodeLabel string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagWriteRepository.UnlinkTagByIdFromEntity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("tagId", tagId), log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))

	cypher := fmt.Sprintf(`
		MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})<-[rel:TAGGED]-(e:%s {id:$entityId})
		DELETE rel`, linkedEntityNodeLabel+"_"+tenant)
	params := map[string]any{
		"tenant":   tenant,
		"id":       tagId,
		"entityId": linkedEntityId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
