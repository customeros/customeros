package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapUserInputToEntity(input model.UserInput) *entity.UserEntity {
	userEntity := entity.UserEntity{
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
	}
	return &userEntity
}

func MapUserUpdateInputToEntity(input model.UserUpdateInput) *entity.UserEntity {
	userEntity := entity.UserEntity{
		Id:            input.ID,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		SourceOfTruth: entity.DataSourceOpenline,
	}
	return &userEntity
}

func MapEntityToUser(userEntity *entity.UserEntity) *model.User {
	return &model.User{
		ID:        userEntity.Id,
		FirstName: userEntity.FirstName,
		LastName:  userEntity.LastName,
		Email:     userEntity.Email,
		CreatedAt: userEntity.CreatedAt,
		Source:    MapDataSourceToModel(userEntity.Source),
	}
}

func MapEntitiesToUsers(userEntities *entity.UserEntities) []*model.User {
	var users []*model.User
	for _, userEntity := range *userEntities {
		users = append(users, MapEntityToUser(&userEntity))
	}
	return users
}
