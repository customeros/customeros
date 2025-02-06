package neo4j_repository

import (
	"context"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type WorkspaceWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, workspace neo4j_entity.WorkspaceEntity) (*dbtype.Node, error)
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

func (r *workspaceWriteRepository) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, workspace neo4j_entity.WorkspaceEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkspaceWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	tracing.TagTenant(span, tenant)

	cypher := `
				MATCH (t:Tenant {name: $tenant}) 
				MERGE (w:Workspace {name: $name, provider: $provider}) 
				ON CREATE SET 
				  w.id = randomUUID(), 
				  w.createdAt = datetime(), 
				  w.updatedAt = datetime(), 
				  w.source = $source, 
				  w.sourceOfTruth = $sourceOfTruth, 
				  w.appSource = $appSource 
				
				MERGE (t)-[:HAS_WORKSPACE]->(w) 
				
				RETURN t, w`
	params := map[string]any{
		"tenant":        tenant,
		"name":          workspace.Name,
		"provider":      workspace.Provider,
		"source":        utils.StringFirstNonEmpty(workspace.Source.String(), neo4j_entity.DataSourceOpenline.String()),
		"sourceOfTruth": utils.StringFirstNonEmpty(workspace.SourceOfTruth.String(), neo4j_entity.DataSourceOpenline.String()),
		"appSource":     workspace.AppSource,
	}
	tracing.LogObjectAsJson(span, "params", params)
	span.LogFields(log.String("cypher", cypher))

	queryResult, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		qr, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, qr, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return queryResult.(*neo4j.Node), nil
}
