package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
)

type WorkspaceWriteRepository interface {
	Merge(ctx context.Context, workspace entity.WorkspaceEntity) (*dbtype.Node, error)
}

type workspaceWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewWorkspaceWriteRepository(driver *neo4j.DriverWithContext, database string) WorkspaceWriteRepository {
	return &workspaceWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *workspaceWriteRepository) Merge(ctx context.Context, workspace entity.WorkspaceEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MERGE (w:Workspace {name:$name, provider:$provider}) " +
		" ON CREATE SET " +
		"  w.id=randomUUID(), " +
		"  w.createdAt=datetime(), " +
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
				"source":        utils.StringFirstNonEmpty(workspace.Source.String(), entity.DataSourceOpenline.String()),
				"sourceOfTruth": utils.StringFirstNonEmpty(workspace.SourceOfTruth.String(), entity.DataSourceOpenline.String()),
				"appSource":     workspace.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
