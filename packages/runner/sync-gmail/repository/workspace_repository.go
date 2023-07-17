package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type WorkspaceRepository interface {
	GetWorkspaceForTenantByName(ctx context.Context, tenant, name string) (*dbtype.Node, error)
}

type workspaceRepository struct {
	driver *neo4j.DriverWithContext
}

func NewWorkspaceRepository(driver *neo4j.DriverWithContext) WorkspaceRepository {
	return &workspaceRepository{
		driver: driver,
	}
}

func (r *workspaceRepository) GetWorkspaceForTenantByName(ctx context.Context, tenant, name string) (*dbtype.Node, error) {
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
		if err != nil && err.Error() == "Result contains no more records" {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		return result.(*dbtype.Node), nil
	}
}
