package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToPersonUser(entity *entity.UserEntity) *model.PersonUser {
	return &model.PersonUser{
		User:    MapEntityToUser(entity),
		Default: entity.DefaultForPerson,
		Tenant:  entity.Tenant,
	}
}

func MapEntitiesToPersonUsers(entities *entity.UserEntities) []*model.PersonUser {
	var personUsers []*model.PersonUser
	for _, entity := range *entities {
		personUsers = append(personUsers, MapEntityToPersonUser(&entity))
	}
	return personUsers
}
