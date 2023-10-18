package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SocialRepository interface {
	MergeSocialFor(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, event events.OrganizationAddSocialEvent) error
}

type socialRepository struct {
	driver *neo4j.DriverWithContext
}

func NewSocialRepository(driver *neo4j.DriverWithContext) SocialRepository {
	return &socialRepository{
		driver: driver,
	}
}

func (r *socialRepository) MergeSocialFor(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, event events.OrganizationAddSocialEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialRepository.MergeSocialFor")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))

	query := fmt.Sprintf(`
		MATCH (e:%s {id:$entityId})
		OPTIONAL MATCH (e)-[r:HAS]->(checkSoc:Social {url:$url})
		FOREACH (ignore IN CASE WHEN checkSoc IS NULL OR checkSoc.id = $id THEN [1] ELSE [] END |
		MERGE (e)-[:HAS]->(soc:Social {id:$id})
		ON CREATE SET 
			soc.createdAt=$createdAt, 
			soc.updatedAt=$updatedAt, 
			soc.source=$source, 
		  	soc.sourceOfTruth=$sourceOfTruth, 
		  	soc.appSource=$appSource, 
		  	soc.platformName=$platformName,
		  	soc.url=$url,
		  	soc.syncedWithEventStore=true,
		  	soc:Social_%s
		ON MATCH SET
			soc.syncedWithEventStore=true)`, linkedEntityNodeLabel+"_"+tenant, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	return r.executeQuery(ctx, query, map[string]any{
		"entityId":      linkedEntityId,
		"id":            event.SocialId,
		"createdAt":     event.CreatedAt,
		"updatedAt":     event.UpdatedAt,
		"platformName":  event.PlatformName,
		"url":           event.Url,
		"source":        helper.GetSource(event.Source),
		"sourceOfTruth": helper.GetSourceOfTruth(event.SourceOfTruth),
		"appSource":     helper.GetAppSource(event.AppSource),
	})
}

// Common database interaction method
func (r *socialRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteQuery(ctx, *r.driver, query, params)
}
