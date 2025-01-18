package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToUser(userEntity *neo4jentity.UserEntity) *model.User {
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
		Bot:             userEntity.Bot,
		Test:            userEntity.Test,
		ProfilePhotoURL: utils.StringPtr(userEntity.ProfilePhotoUrl),
		Onboarding: &model.UserOnboardingDetails{
			ShowOnboardingPage:               userEntity.OnboardingDetails.ShowOnboardingPage,
			OnboardingInboundStepCompleted:   userEntity.OnboardingDetails.OnboardingInboundStepCompleted,
			OnboardingOutboundStepCompleted:  userEntity.OnboardingDetails.OnboardingOutboundStepCompleted,
			OnboardingCrmStepCompleted:       userEntity.OnboardingDetails.OnboardingCrmStepCompleted,
			OnboardingMailstackStepCompleted: userEntity.OnboardingDetails.OnboardingMailstackStepCompleted,
		},
	}
}

func MapRoleToModel(role string) model.Role {
	switch role {
	case "ADMIN":
		return model.RoleAdmin
	case "USER":
		return model.RoleUser
	case "PLATFORM_OWNER":
		return model.RolePlatformOwner
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

func MapEntitiesToUsers(userEntities *neo4jentity.UserEntities) []*model.User {
	var users []*model.User
	for _, userEntity := range *userEntities {
		users = append(users, MapEntityToUser(&userEntity))
	}
	return users
}
