package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type WorkspaceService interface {
	MergeToTenant(ctx context.Context, workspaceEntity neo4j_entity.WorkspaceEntity, tenant string) (bool, error)
	GetWorkspaceDomainsForTenant(ctx context.Context) ([]string, error)
	CheckEmailBelongsToTenant(ctx context.Context, email string) (bool, error)
}
