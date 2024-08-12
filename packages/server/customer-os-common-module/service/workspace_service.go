package service

import (
	"context"
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type WorkspaceService interface {
	MergeToTenant(ctx context.Context, workspaceEntity neo4jentity.WorkspaceEntity, tenant string) (bool, error)
}

type workspaceService struct {
	services *Services
}

func NewWorkspaceService(services *Services) WorkspaceService {
	return &workspaceService{
		services: services,
	}
}

func (s *workspaceService) MergeToTenant(ctx context.Context, workspaceEntity neo4jentity.WorkspaceEntity, tenant string) (bool, error) {
	_, err := s.services.Neo4jRepositories.WorkspaceWriteRepository.Merge(ctx, workspaceEntity)
	if err != nil {
		return false, fmt.Errorf("MergeToTenant: %w", err)
	}
	result, err := s.services.Neo4jRepositories.TenantWriteRepository.LinkWithWorkspace(ctx, tenant, workspaceEntity)
	return result, err
}
