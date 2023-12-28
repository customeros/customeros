package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/customer-os-neo4j-repository/model"
	"github.com/openline-ai/customer-os-neo4j-repository/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type SocialFields struct {
	SocialId     string       `json:"socialId"`
	Url          string       `json:"url"`
	PlatformName string       `json:"platformName"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	SourceFields model.Source `json:"sourceFields"`
}

type SocialWriteRepository interface {
	MergeSocialFor(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, data SocialFields) error
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

func (r *socialWriteRepository) MergeSocialFor(ctx context.Context, tenant, linkedEntityId, linkedEntityNodeLabel string, data SocialFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialWriteRepository.MergeSocialFor")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityNodeLabel", linkedEntityNodeLabel))
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`
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
	params := map[string]any{
		"entityId":      linkedEntityId,
		"id":            data.SocialId,
		"createdAt":     data.CreatedAt,
		"updatedAt":     data.UpdatedAt,
		"platformName":  data.PlatformName,
		"url":           data.Url,
		"source":        data.SourceFields.Source,
		"sourceOfTruth": data.SourceFields.SourceOfTruth,
		"appSource":     data.SourceFields.AppSource,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
