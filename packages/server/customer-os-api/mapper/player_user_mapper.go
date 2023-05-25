package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToPlayerUser(entity *entity.UserEntity) *model.PlayerUser {
	return &model.PlayerUser{
		User:    MapEntityToUser(entity),
		Default: entity.DefaultForPlayer,
		Tenant:  entity.Tenant,
	}
}

func MapEntitiesToPlayerUsers(entities *entity.UserEntities) []*model.PlayerUser {
	var playerUsers []*model.PlayerUser
	for _, entity := range *entities {
		playerUsers = append(playerUsers, MapEntityToPlayerUser(&entity))
	}
	return playerUsers
}
