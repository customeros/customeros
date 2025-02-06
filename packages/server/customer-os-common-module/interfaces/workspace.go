package interfaces

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type WorkspaceService interface {
	MergeToTenant(ctx context.Context, tx *neo4j.ManagedTransaction, workspaceEntity neo4j_entity.WorkspaceEntity, tenant string) (bool, error)
	GetWorkspaceDomainsForTenant(ctx context.Context) ([]string, error)
}
