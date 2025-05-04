package neo4j_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "WorkspaceWriteRepository.Merge")
	defer spans.Finish()

	cypher := `
				MATCH (t:Tenant {name: $tenant}) 
				MERGE (w:Workspace {name: $name, provider: $provider}) 
				ON CREATE SET 
				  w.id = randomUUID(), 
				  w.createdAt = datetime(), 
				  w.updatedAt = datetime(), 
				  w.source = $source, 
				  w.appSource = $appSource 
				
				MERGE (t)-[:HAS_WORKSPACE]->(w) 
				
				RETURN t, w`
	params := map[string]any{
		"tenant":    tenant,
		"name":      workspace.Name,
		"provider":  workspace.Provider,
		"source":    utils.StringFirstNonEmpty(workspace.Source.String(), neo4j_entity.DataSourceOpenline.String()),
		"appSource": workspace.AppSource,
	}
	spans.LogObjectAsJson("params", params)
	spans.LogKV("cypher", cypher)

	queryResult, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		qr, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, qr, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return queryResult.(*neo4j.Node), nil
}
