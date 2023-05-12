package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type WorkspaceService interface {
	MergeToTenant(ctx context.Context, workspaceEntity entity.WorkspaceEntity, tenant string) (bool, error)
}

type workspaceService struct {
	repositories *repository.Repositories
}

func NewWorkspaceService(repository *repository.Repositories) WorkspaceService {
	return &workspaceService{
		repositories: repository,
	}
}

func (s *workspaceService) MergeToTenant(ctx context.Context, workspaceEntity entity.WorkspaceEntity, tenant string) (bool, error) {
	_, err := s.repositories.WorkspaceRepository.Merge(ctx, workspaceEntity)
	if err != nil {
		return false, fmt.Errorf("MergeToTenant: %w", err)
	}
	result, err := s.repositories.TenantRepository.LinkWithWorkspace(ctx, tenant, workspaceEntity)
	return result, err
}

func (s *workspaceService) mapDbNodeToWorkspaceEntity(dbNode dbtype.Node) *entity.WorkspaceEntity {
	props := utils.GetPropsFromNode(dbNode)
	domain := entity.WorkspaceEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "domain"),
		Provider:      utils.GetStringPropOrEmpty(props, "provider"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &domain
}
