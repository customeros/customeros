package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantSocialIdAndEntityId struct {
	Tenant         string
	SocialId       string
	LinkedEntityId string
}

type SocialReadRepository interface {
	GetDuplicatedSocialsForEntityType(ctx context.Context, linkedEntityNodeLabel string, minutesSinceLastUpdate, limit int) ([]TenantSocialIdAndEntityId, error)
	GetEmptySocialsForEntityType(ctx context.Context, linkedEntityNodeLabel string, minutesSinceLastUpdate, limit int) ([]TenantSocialIdAndEntityId, error)
}

type socialReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewSocialReadRepository(driver *neo4j.DriverWithContext, database string) SocialReadRepository {
	return &socialReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *socialReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *socialReadRepository) GetDuplicatedSocialsForEntityType(ctx context.Context, linkedEntityNodeLabel string, minutesSinceLastUpdate, limit int) ([]TenantSocialIdAndEntityId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialReadRepository.GetDuplicatedSocialsForEntityType")
	defer span.Finish()
	span.LogFields(log.String("linkedEntityNodeLabel", linkedEntityNodeLabel), log.Int("minutesSinceLastUpdate", minutesSinceLastUpdate), log.Int("limit", limit))

	cypher := fmt.Sprintf(`MATCH (t:Tenant)--(e:%s)-[:HAS]->(s:Social)
					WHERE s.updatedAt < datetime() - duration({minutes: $minutesSinceLastUpdate})
					WITH t, e, s.url AS url, COLLECT(s) AS socials
					WHERE size(socials) > 1
					WITH t, e, url, last(socials) AS lastSocial 
					RETURN t.name, e.id, lastSocial.id limit $limit`, linkedEntityNodeLabel)
	params := map[string]any{
		"minutesSinceLastUpdate": minutesSinceLastUpdate,
		"limit":                  limit,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantSocialIdAndEntityId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantSocialIdAndEntityId{
				Tenant:         v.Values[0].(string),
				LinkedEntityId: v.Values[1].(string),
				SocialId:       v.Values[2].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}

func (r *socialReadRepository) GetEmptySocialsForEntityType(ctx context.Context, linkedEntityNodeLabel string, minutesSinceLastUpdate, limit int) ([]TenantSocialIdAndEntityId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialReadRepository.GetDuplicatedSocialsForEntityType")
	defer span.Finish()
	span.LogFields(log.String("linkedEntityNodeLabel", linkedEntityNodeLabel), log.Int("minutesSinceLastUpdate", minutesSinceLastUpdate), log.Int("limit", limit))

	cypher := fmt.Sprintf(`MATCH (t:Tenant)--(e:%s)-[:HAS]->(s:Social)
					WHERE s.updatedAt < datetime() - duration({minutes: $minutesSinceLastUpdate})
					AND (s.url IS NULL OR s.url = "")
					RETURN t.name, e.id, s.id ORDER by s.createdAt ASC LIMIT $limit`, linkedEntityNodeLabel)
	params := map[string]any{
		"minutesSinceLastUpdate": minutesSinceLastUpdate,
		"limit":                  limit,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)

	})
	if err != nil {
		return nil, err
	}
	output := make([]TenantSocialIdAndEntityId, 0)
	for _, v := range records.([]*neo4j.Record) {
		output = append(output,
			TenantSocialIdAndEntityId{
				Tenant:         v.Values[0].(string),
				LinkedEntityId: v.Values[1].(string),
				SocialId:       v.Values[2].(string),
			})
	}
	span.LogFields(log.Int("result.count", len(output)))
	return output, nil
}
