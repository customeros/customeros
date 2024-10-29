package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type SocialFields struct {
	SocialId       string             `json:"socialId"`
	Url            string             `json:"url"`
	Alias          string             `json:"alias"`
	FollowersCount int64              `json:"followersCount"`
	ExternalId     string             `json:"externalId"`
	CreatedAt      time.Time          `json:"createdAt"`
	SourceFields   model.SourceFields `json:"sourceFields"`
}

type SocialWriteRepository interface {
	MergeSocialForEntity(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, data SocialFields) error
	PermanentlyDelete(ctx context.Context, tenant, socialId string) error
	RemoveSocialForEntityById(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel, socialId string) error
	RemoveSocialForEntityByUrl(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel, socialUrl string) error
	Update(ctx context.Context, tenant string, socialEntity neo4jentity.SocialEntity) (*dbtype.Node, error)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`
		MATCH (e:%s {id:$entityId})
		MERGE (e)-[:HAS]->(soc:Social {id:$id})
		ON CREATE SET 
			soc.createdAt=$createdAt, 
			soc.updatedAt=datetime(), 
			soc.source=$source, 
		  	soc.appSource=$appSource, 
		  	soc.url=$url,
			soc.alias=$alias,
			soc.externalId=$externalId,
			soc.followersCount=$followersCount,
		  	soc:Social_%s
		ON MATCH SET
			soc.updatedAt=datetime(),
			soc.alias = CASE WHEN $alias <> '' THEN $alias ELSE soc.alias END,
			soc.externalId = CASE WHEN $externalId <> '' THEN $externalId ELSE soc.externalId END,
			soc.followersCount = CASE WHEN $followersCount <> 0 THEN $followersCount ELSE soc.followersCount END`, linkedEntityNodeLabel+"_"+tenant, tenant)
	params := map[string]any{
		"entityId":       linkedEntityId,
		"id":             data.SocialId,
		"createdAt":      data.CreatedAt,
		"url":            data.Url,
		"alias":          data.Alias,
		"followersCount": data.FollowersCount,
		"externalId":     data.ExternalId,
		"source":         data.SourceFields.Source,
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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

func (r *socialWriteRepository) Update(ctx context.Context, tenant string, socialEntity neo4jentity.SocialEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialWriteRepository.Update")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (soc:Social_%s {id:$id})
			SET soc.updatedAt=datetime(),
				soc.url=$url
			RETURN soc`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"id":  socialEntity.Id,
				"url": socialEntity.Url,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
