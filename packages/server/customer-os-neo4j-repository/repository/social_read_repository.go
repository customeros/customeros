package neo4j_repository

import (
	"context"
	"fmt"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type TenantSocialIdAndEntityId struct {
	Tenant         string
	SocialId       string
	LinkedEntityId string
}

type SocialReadRepository interface {
	GetDuplicatedSocialsForEntityType(ctx context.Context, linkedEntityNodeLabel string, minutesSinceLastUpdate, limit int) ([]TenantSocialIdAndEntityId, error)
	GetEmptySocialsForEntityType(ctx context.Context, linkedEntityNodeLabel string, minutesSinceLastUpdate, limit int) ([]TenantSocialIdAndEntityId, error)
	GetAllForEntities(ctx context.Context, tenant string, linkedEntityType neo4jenum.EntityType, linkedEntityIds []string) ([]*utils.DbNodeAndId, error)
	GetAllLinkedinForEntities(ctx context.Context, tenant string, linkedEntityType neo4jenum.EntityType, linkedEntityIds []string) ([]*utils.DbNodeAndId, error)
	GetById(ctx context.Context, tenant, socialId string) (*dbtype.Node, error)
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

func (r *socialReadRepository) GetById(ctx context.Context, tenant, socialId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "SocialReadRepository.GetById")
	defer spans.Finish()

	spans.TagEntity(socialId)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`MATCH (s:Social_%s {id: $socialId}) RETURN s`, tenant)
	params := map[string]any{
		"socialId": socialId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
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
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *socialReadRepository) GetDuplicatedSocialsForEntityType(ctx context.Context, linkedEntityNodeLabel string, minutesSinceLastUpdate, limit int) ([]TenantSocialIdAndEntityId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "SocialReadRepository.GetDuplicatedSocialsForEntityType")
	defer spans.Finish()
	spans.LogKV("linkedEntityNodeLabel", linkedEntityNodeLabel)
	spans.LogKV("minutesSinceLastUpdate", minutesSinceLastUpdate)
	spans.LogKV("limit", limit)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *socialReadRepository) GetEmptySocialsForEntityType(ctx context.Context, linkedEntityNodeLabel string, minutesSinceLastUpdate, limit int) ([]TenantSocialIdAndEntityId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "SocialReadRepository.GetDuplicatedSocialsForEntityType")
	defer spans.Finish()
	spans.LogKV("linkedEntityNodeLabel", linkedEntityNodeLabel)
	spans.LogKV("minutesSinceLastUpdate", minutesSinceLastUpdate)
	spans.LogKV("limit", limit)

	cypher := fmt.Sprintf(`MATCH (t:Tenant)--(e:%s)-[:HAS]->(s:Social)
					WHERE s.updatedAt < datetime() - duration({minutes: $minutesSinceLastUpdate})
					AND (s.url IS NULL OR s.url = "")
					RETURN t.name, e.id, s.id ORDER by s.createdAt ASC LIMIT $limit`, linkedEntityNodeLabel)
	params := map[string]any{
		"minutesSinceLastUpdate": minutesSinceLastUpdate,
		"limit":                  limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans.LogKV("result.count", len(output))
	return output, nil
}

func (r *socialReadRepository) GetAllForEntities(ctx context.Context, tenant string, linkedEntityType neo4jenum.EntityType, linkedEntityIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "SocialReadRepository.GetAllForEntities")
	defer spans.Finish()

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`MATCH (e:%s)-[:HAS]->(soc:Social)
			WHERE e.id IN $entityIds
			RETURN soc, e.id as entityId ORDER BY soc.url`, linkedEntityType.Neo4jLabel()+"_"+tenant)
	params := map[string]any{
		"entityIds": linkedEntityIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	dbNodeAndIds := result.([]*utils.DbNodeAndId)
	spans.LogKV("result.count", len(dbNodeAndIds))
	return dbNodeAndIds, err
}

func (r *socialReadRepository) GetAllLinkedinForEntities(ctx context.Context, tenant string, linkedEntityType neo4jenum.EntityType, linkedEntityIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "SocialReadRepository.GetAllLinkedinForEntities")
	defer spans.Finish()

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := fmt.Sprintf(`MATCH (e:%s)-[:HAS]->(soc:Social)
			WHERE e.id IN $entityIds and soc.url contains 'linkedin.com'
			RETURN soc, e.id as entityId ORDER BY soc.url`, linkedEntityType.Neo4jLabel()+"_"+tenant)
	params := map[string]any{
		"entityIds": linkedEntityIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	dbNodeAndIds := result.([]*utils.DbNodeAndId)
	spans.LogKV("result.count", len(dbNodeAndIds))
	return dbNodeAndIds, err
}
