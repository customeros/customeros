package mapper

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapWorkspaceInputToEntity(input model.WorkspaceInput) entity.WorkspaceEntity {
	workspaceEntity := entity.WorkspaceEntity{
		Name:          input.Name,
		Provider:      input.Provider,
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilString(input.AppSource),
	}
	if len(workspaceEntity.AppSource) == 0 {
		workspaceEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return workspaceEntity
}

func MapEntitiesToWorkspaces(entities *entity.WorkspaceEntities) []*model.Workspace {
	var workspaces []*model.Workspace
	for _, workspaceEntity := range *entities {
		workspaces = append(workspaces, MapEntityToWorkspace(&workspaceEntity))
	}
	return workspaces
}

func MapEntityToWorkspace(entity *entity.WorkspaceEntity) *model.Workspace {
	return &model.Workspace{
		ID:            entity.Id,
		Name:          entity.Name,
		Provider:      entity.Provider,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
}
