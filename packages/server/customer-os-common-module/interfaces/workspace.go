package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type WorkspaceService interface {
	MergeToTenant(ctx context.Context, workspaceEntity entity.WorkspaceEntity, tenant string) (bool, error)
	GetWorkspaceDomainsForTenant(ctx context.Context) ([]string, error)
	CheckEmailBelongsToTenant(ctx context.Context, email string) (bool, error)
}
