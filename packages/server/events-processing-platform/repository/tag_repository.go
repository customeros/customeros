package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type TagRepository interface {
	AddTagByIdTo(ctx context.Context, tenant, tagId, linkedEntityId, linkedEntityNodeLabel string, taggedAt time.Time) error
	RemoveTagByIdFrom(ctx context.Context, tenant, tagId, linkedEntityId, linkedEntityNodeLabel string) error
}

type tagRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTagRepository(driver *neo4j.DriverWithContext) TagRepository {
	return &tagRepository{
		driver: driver,
	}
}

func (r *tagRepository) AddTagByIdTo(ctx context.Context, tenant, tagId, linkedEntityId, linkedEntityNodeLabel string, taggedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.AddTagByIdTo")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("tagId", tagId), log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel), log.Object("taggedAt", taggedAt))

	query := fmt.Sprintf(`
		MATCH (e:%s {id:$entityId}),
			(t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
		MERGE (e)-[rel:TAGGED]->(tag)
		ON CREATE SET
			rel.taggedAt=$taggedAt `, linkedEntityNodeLabel+"_"+tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":   tenant,
		"id":       tagId,
		"taggedAt": taggedAt,
		"entityId": linkedEntityId,
	})
}

func (r *tagRepository) RemoveTagByIdFrom(ctx context.Context, tenant, tagId, linkedEntityId, linkedEntityNodeLabel string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.RemoveTagByIdFrom")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("tagId", tagId), log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))

	query := fmt.Sprintf(`
		MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})<-[rel:TAGGED]-(e:%s {id:$entityId})
		DELETE rel`, linkedEntityNodeLabel+"_"+tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":   tenant,
		"id":       tagId,
		"entityId": linkedEntityId,
	})
}

// Common database interaction method
func (r *tagRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
