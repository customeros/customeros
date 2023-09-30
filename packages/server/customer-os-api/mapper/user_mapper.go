package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapUserInputToEntity(input model.UserInput) *entity.UserEntity {
	userEntity := entity.UserEntity{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Name:            utils.IfNotNilString(input.Name),
		Source:          entity.DataSourceOpenline,
		SourceOfTruth:   entity.DataSourceOpenline,
		Timezone:        utils.IfNotNilString(input.Timezone),
		ProfilePhotoUrl: utils.IfNotNilString(input.ProfilePhotoURL),
		AppSource:       utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &userEntity
}

func MapEntityToUser(userEntity *entity.UserEntity) *model.User {
	if userEntity == nil {
		return nil
	}
	return &model.User{
		ID:              userEntity.Id,
		FirstName:       userEntity.FirstName,
		LastName:        userEntity.LastName,
		Name:            utils.StringPtrNillable(userEntity.Name),
		Timezone:        utils.StringPtrNillable(userEntity.Timezone),
		CreatedAt:       userEntity.CreatedAt,
		UpdatedAt:       userEntity.UpdatedAt,
		Source:          MapDataSourceToModel(userEntity.Source),
		SourceOfTruth:   MapDataSourceToModel(userEntity.SourceOfTruth),
		Roles:           MapRolesToModel(userEntity.Roles),
		AppSource:       userEntity.AppSource,
		Internal:        userEntity.Internal,
		ProfilePhotoURL: utils.StringPtr(userEntity.ProfilePhotoUrl),
	}
}

func MapRoleToModel(role string) model.Role {
	switch role {
	case "ADMIN":
		return model.RoleAdmin
	case "USER":
		return model.RoleUser
	case "CUSTOMER_OS_PLATFORM_OWNER":
		return model.RoleCustomerOsPlatformOwner
	case "OWNER":
		return model.RoleOwner
	default:
		return model.RoleUser
	}
}

func MapRoleToEntity(role model.Role) string {
	return role.String()
}

func MapRolesToModel(roles []string) []model.Role {
	var modelRoles []model.Role
	for _, role := range roles {
		modelRoles = append(modelRoles, MapRoleToModel(role))
	}
	return modelRoles
}

func MapEntitiesToUsers(userEntities *entity.UserEntities) []*model.User {
	var users []*model.User
	for _, userEntity := range *userEntities {
		users = append(users, MapEntityToUser(&userEntity))
	}
	return users
}
