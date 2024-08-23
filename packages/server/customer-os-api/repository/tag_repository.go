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
