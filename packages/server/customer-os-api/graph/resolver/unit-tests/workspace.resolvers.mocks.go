package unit_tests

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/mock"
)

type MockedWorkspaceService struct {
	mock.Mock
}

type WorkspaceService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func (s *MockedWorkspaceService) MergeToTenant(ctx context.Context, workspaceEntity entity.WorkspaceEntity, tenant string) (bool, error) {
	return true, nil
}

func NewWorkspaceService(log logger.Logger, repository *repository.Repositories) *WorkspaceService {
	return &WorkspaceService{
		log:          log,
		repositories: repository,
	}
}

func (s *MockedWorkspaceService) mapDbNodeToWorkspaceEntity(dbNode dbtype.Node) *entity.WorkspaceEntity {
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
