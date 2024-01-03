package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapPlayerInputToEntity(input *model.PlayerInput) *entity.PlayerEntity {
	if input == nil {
		return nil
	}
	playerEntity := entity.PlayerEntity{
		IdentityId:    input.IdentityID,
		AuthId:        input.AuthID,
		Provider:      input.Provider,
		CreatedAt:     utils.Now(),
		UpdatedAt:     utils.Now(),
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &playerEntity
}

func MapPlayerUpdateToEntity(id string, input *model.PlayerUpdate) *entity.PlayerEntity {
	if input == nil {
		return nil
	}
	playerEntity := entity.PlayerEntity{
		Id:            id,
		IdentityId:    input.IdentityID,
		UpdatedAt:     utils.Now(),
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &playerEntity
}

func MapEntityToPlayer(entity *entity.PlayerEntity) *model.Player {
	if entity == nil {
		return nil
	}
	return &model.Player{
		ID:            entity.Id,
		IdentityID:    entity.IdentityId,
		AuthID:        entity.AuthId,
		Provider:      entity.Provider,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
}
