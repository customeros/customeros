package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type WorkspaceRepository interface {
	Merge(ctx context.Context, workspace entity.WorkspaceEntity) (*dbtype.Node, error)
}

type workspaceRepository struct {
	driver *neo4j.DriverWithContext
}

func NewWorkspaceRepository(driver *neo4j.DriverWithContext) WorkspaceRepository {
	return &workspaceRepository{
		driver: driver,
	}
}

func (r *workspaceRepository) Merge(ctx context.Context, workspace entity.WorkspaceEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MERGE (w:Workspace {name:$name, provider:$provider}) " +
		" ON CREATE SET " +
		"  w.id=randomUUID(), " +
		"  w.createdAt=$now, " +
		"  w.updatedAt=datetime(), " +
		"  w.source=$source, " +
		"  w.sourceOfTruth=$sourceOfTruth, " +
		"  w.appSource=$appSource " +
		" RETURN w"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"name":          workspace.Name,
				"provider":      workspace.Provider,
				"source":        workspace.Source,
				"sourceOfTruth": workspace.SourceOfTruth,
				"appSource":     workspace.AppSource,
				"now":           utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
