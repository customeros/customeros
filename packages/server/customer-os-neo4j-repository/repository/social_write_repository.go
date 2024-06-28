package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type SocialFields struct {
	SocialId       string       `json:"socialId"`
	Url            string       `json:"url"`
	Alias          string       `json:"alias"`
	FollowersCount int64        `json:"followersCount"`
	CreatedAt      time.Time    `json:"createdAt"`
	SourceFields   model.Source `json:"sourceFields"`
}

type SocialWriteRepository interface {
	MergeSocialForEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, data SocialFields) error
	PermanentlyDelete(ctx context.Context, tenant, socialId string) error
	RemoveSocialForEntityById(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel, socialId string) error
	RemoveSocialForEntityByUrl(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel, socialUrl string) error
}

type socialWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewSocialWriteRepository(driver *neo4j.DriverWithContext, database string) SocialWriteRepository {
	return &socialWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *socialWriteRepository) MergeSocialForEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, data SocialFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialWriteRepository.MergeSocialForEntity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`
		MATCH (e:%s {id:$entityId})
		MERGE (e)-[:HAS]->(soc:Social {id:$id})
		ON CREATE SET 
			soc.createdAt=$createdAt, 
			soc.updatedAt=datetime(), 
			soc.source=$source, 
		  	soc.sourceOfTruth=$sourceOfTruth, 
		  	soc.appSource=$appSource, 
		  	soc.url=$url,
			soc.alias=$alias,
			soc.followersCount=$followersCount,
		  	soc.syncedWithEventStore=true,
		  	soc:Social_%s
		ON MATCH SET
			soc.updatedAt=datetime(),
			soc.alias = CASE WHEN $alias <> '' THEN $alias ELSE soc.alias END,
			soc.followersCount = CASE WHEN $followersCount <> 0 THEN $followersCount ELSE soc.followersCount END,
			soc.syncedWithEventStore=true`, linkedEntityNodeLabel+"_"+tenant, tenant)
	params := map[string]any{
		"entityId":       linkedEntityId,
		"id":             data.SocialId,
		"createdAt":      data.CreatedAt,
		"url":            data.Url,
		"alias":          data.Alias,
		"followersCount": data.FollowersCount,
		"source":         data.SourceFields.Source,
		"sourceOfTruth":  data.SourceFields.SourceOfTruth,
		"appSource":      data.SourceFields.AppSource,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *socialWriteRepository) PermanentlyDelete(ctx context.Context, tenant, socialId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialWriteRepository.PermanentlyDelete")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := fmt.Sprintf(`MATCH (soc:Social_%s {id:$socialId}) DETACH DELETE soc`, tenant)
	params := map[string]any{
		"socialId": socialId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *socialWriteRepository) RemoveSocialForEntityById(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel, socialId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialWriteRepository.RemoveSocialForEntityById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel), log.String("socialId", socialId))

	// delete social only if has no other relations
	cypher := fmt.Sprintf(`
		MATCH (e:%s {id:$entityId})-[r:HAS]->(soc:Social {id:$socialId})
		DELETE r
		WITH soc
		WHERE NOT (soc)--()
		DELETE soc`, linkedEntityNodeLabel+"_"+tenant)
	params := map[string]any{
		"entityId": linkedEntityId,
		"socialId": socialId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *socialWriteRepository) RemoveSocialForEntityByUrl(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel, socialUrl string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialWriteRepository.RemoveSocialForEntityByUrl")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel), log.String("socialUrl", socialUrl))

	// delete social only if has no other relations
	cypher := fmt.Sprintf(`
		MATCH (e:%s {id:$entityId})-[r:HAS]->(soc:Social {url:$url})
		DELETE r
		WITH soc
		WHERE NOT (soc)--()
		DELETE soc`, linkedEntityNodeLabel+"_"+tenant)
	params := map[string]any{
		"entityId": linkedEntityId,
		"url":      socialUrl,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
