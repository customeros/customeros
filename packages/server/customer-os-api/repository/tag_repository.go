package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
)

type TagRepository interface {
	Merge(ctx context.Context, tenant string, tag neo4jentity.TagEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tenant string, tag neo4jentity.TagEntity) (*dbtype.Node, error)
	UnlinkAndDelete(ctx context.Context, tenant string, tagId string) error
}

type tagRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTagRepository(driver *neo4j.DriverWithContext) TagRepository {
	return &tagRepository{
		driver: driver,
	}
}

func (r *tagRepository) Merge(ctx context.Context, tenant string, tag neo4jentity.TagEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {name:$name}) 
		 ON CREATE SET 
		  tag.id=randomUUID(),
		  tag.createdAt=$now,
		  tag.updatedAt=datetime(),
		  tag.source=$source,
		  tag.sourceOfTruth=$sourceOfTruth,
		  tag.appSource=$appSource,
		  tag:Tag_%s
		 RETURN tag`, tenant)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":        tenant,
				"name":          tag.Name,
				"source":        tag.Source,
				"sourceOfTruth": tag.SourceOfTruth,
				"appSource":     tag.AppSource,
				"now":           utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tagRepository) Update(ctx context.Context, tenant string, tag neo4jentity.TagEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
			SET tag.name=$name, tag.updatedAt=datetime()
			RETURN tag`,
			map[string]any{
				"tenant": tenant,
				"id":     tag.Id,
				"name":   tag.Name,
				"now":    utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tagRepository) UnlinkAndDelete(ctx context.Context, tenant string, tagId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TagRepository.UnlinkAndDelete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:TAG_BELONGS_TO_TENANT]-(tag:Tag {id:$id})
			DETACH DELETE tag`,
			map[string]any{
				"tenant": tenant,
				"id":     tagId,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}
