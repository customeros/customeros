package neo4j_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type WorkspaceReadRepository interface {
	GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetByName(ctx context.Context, tenant, name string) (*dbtype.Node, error)
	GetByNameCrossTenant(ctx context.Context, name string) ([]*dbtype.Node, error)
}

type workspaceReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewWorkspaceReadRepository(driver *neo4j.DriverWithContext, database string) WorkspaceReadRepository {
	return &workspaceReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *workspaceReadRepository) GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "WorkspaceReadRepository.GetAllForTenant")
	defer spans.Finish()

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	cypher := `MATCH (t:Tenant {name:$tenant})--(w:Workspace) return w`
	params := map[string]any{
		"tenant": tenant,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), nil
}

func (r *workspaceReadRepository) GetByName(ctx context.Context, tenant, name string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "WorkspaceReadRepository.GetByName")
	defer spans.Finish()

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `MATCH (t:Tenant {name:$tenant})--(w:Workspace{name:$workspaceName}) return w`,
			map[string]any{
				"tenant":        tenant,
				"workspaceName": name,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		if err.Error() == "Result contains no more records" {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *workspaceReadRepository) GetByNameCrossTenant(ctx context.Context, name string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "WorkspaceReadRepository.GetByNameCrossTenant")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant)--(w:Workspace{name:$workspaceName}) return w`
	params := map[string]any{
		"workspaceName": name,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), nil
}
