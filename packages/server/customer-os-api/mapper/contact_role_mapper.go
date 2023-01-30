package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapContactRoleInputToEntity(input *model.ContactRoleInput) *entity.ContactRoleEntity {
	if input == nil {
		return nil
	}
	contactRoleEntity := entity.ContactRoleEntity{
		JobTitle:            utils.IfNotNilString(input.JobTitle),
		Primary:             utils.IfNotNilBool(input.Primary),
		ResponsibilityLevel: utils.IfNotNilInt64(input.ResponsibilityLevel),
		Source:              entity.DataSourceOpenline,
		SourceOfTruth:       entity.DataSourceOpenline,
		AppSource:           utils.IfNotNilString(input.AppSource),
	}
	if len(contactRoleEntity.AppSource) == 0 {
		contactRoleEntity.AppSource = common.AppSourceCustomerOsApi
	}
	return &contactRoleEntity
}

func MapContactRoleUpdateInputToEntity(input *model.ContactRoleUpdateInput) *entity.ContactRoleEntity {
	if input == nil {
		return nil
	}
	contactRoleEntity := entity.ContactRoleEntity{
		Id:                  input.ID,
		JobTitle:            utils.IfNotNilString(input.JobTitle),
		Primary:             utils.IfNotNilBool(input.Primary),
		ResponsibilityLevel: utils.IfNotNilInt64(input.ResponsibilityLevel),
		SourceOfTruth:       entity.DataSourceOpenline,
	}
	return &contactRoleEntity
}

func MapEntityToContactRole(entity *entity.ContactRoleEntity) *model.ContactRole {
	contactRole := model.ContactRole{
		ID:                  entity.Id,
		Primary:             entity.Primary,
		Source:              MapDataSourceToModel(entity.Source),
		SourceOfTruth:       MapDataSourceToModel(entity.SourceOfTruth),
		ResponsibilityLevel: entity.ResponsibilityLevel,
		AppSource:           entity.AppSource,
		CreatedAt:           entity.CreatedAt,
		UpdatedAt:           entity.UpdatedAt,
	}
	if len(entity.JobTitle) > 0 {
		contactRole.JobTitle = utils.StringPtr(entity.JobTitle)
	}
	return &contactRole
}

func MapEntitiesToContactRoles(entities *entity.ContactRoleEntities) []*model.ContactRole {
	var contactRoles []*model.ContactRole
	for _, contactRoleEntity := range *entities {
		contactRoles = append(contactRoles, MapEntityToContactRole(&contactRoleEntity))
	}
	return contactRoles
}
