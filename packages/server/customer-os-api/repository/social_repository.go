package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SocialRepository interface {
	CreateSocialForEntity(ctx context.Context, tenant string, linkedEntityType entity.EntityType, linkedEntityId string, socialEntity entity.SocialEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tenant string, socialEntity entity.SocialEntity) (*dbtype.Node, error)
	GetAllForEntities(ctx context.Context, tenant string, linkedEntityType entity.EntityType, linkedEntityIds []string) ([]*utils.DbNodeAndId, error)
	Remove(ctx context.Context, socialId string) error
}

type socialRepository struct {
	driver *neo4j.DriverWithContext
}

func NewSocialRepository(driver *neo4j.DriverWithContext) SocialRepository {
	return &socialRepository{
		driver: driver,
	}
}

func (r *socialRepository) CreateSocialForEntity(ctx context.Context, tenant string, linkedEntityType entity.EntityType, linkedEntityId string, socialEntity entity.SocialEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialRepository.CreateSocialForEntity")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (e:%s {id:$entityId})
		 MERGE (e)-[:HAS]->(soc:Social {id:randomUUID()})
		 ON CREATE SET 
		  soc.createdAt=$now, 
		  soc.updatedAt=$now, 
		  soc.source=$source, 
		  soc.sourceOfTruth=$sourceOfTruth, 
		  soc.appSource=$appSource, 
		  soc.platformName=$platformName,
		  soc.url=$url,
		  soc:%s
		 RETURN soc`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, linkedEntityType.Neo4jLabel()+"_"+tenant, "Social_"+tenant),
			map[string]any{
				"tenant":        tenant,
				"now":           utils.Now(),
				"entityId":      linkedEntityId,
				"platformName":  socialEntity.PlatformName,
				"url":           socialEntity.Url,
				"source":        socialEntity.SourceFields.Source,
				"sourceOfTruth": socialEntity.SourceFields.SourceOfTruth,
				"appSource":     socialEntity.SourceFields.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *socialRepository) Update(ctx context.Context, tenant string, socialEntity entity.SocialEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (soc:Social_%s {id:$id})
			SET soc.updatedAt=$now,
				soc.platformName=$platformName,
				soc.url=$url,
				soc.sourceOfTruth=$sourceOfTruth
			RETURN soc`

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"now":           utils.Now(),
				"id":            socialEntity.Id,
				"platformName":  socialEntity.PlatformName,
				"url":           socialEntity.Url,
				"sourceOfTruth": socialEntity.SourceFields.SourceOfTruth,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *socialRepository) GetAllForEntities(ctx context.Context, tenant string, linkedEntityType entity.EntityType, linkedEntityIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialRepository.GetAllForEntities")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (e:%s)-[:HAS]->(soc:Social)
			WHERE e.id IN $entityIds
			RETURN soc, e.id as entityId ORDER BY soc.platformName`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, linkedEntityType.Neo4jLabel()+"_"+tenant),
			map[string]any{
				"entityIds": linkedEntityIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *socialRepository) Remove(ctx context.Context, socialId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialRepository.Remove")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("socialId", socialId))

	query := fmt.Sprintf(`MATCH (soc:Social_%s {id:$socialId}) DETACH DELETE soc`, common.GetTenantFromContext(ctx))
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]any{
				"socialId": socialId,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}
