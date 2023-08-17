package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapJobRoleInputToEntity(input *model.JobRoleInput) *entity.JobRoleEntity {
	if input == nil {
		return nil
	}
	jobRoleEntity := entity.JobRoleEntity{
		JobTitle:      utils.IfNotNilString(input.JobTitle),
		Primary:       utils.IfNotNilBool(input.Primary),
		Description:   input.Description,
		Company:       input.Company,
		StartedAt:     input.StartedAt,
		EndedAt:       input.EndedAt,
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilString(input.AppSource),
	}
	if len(jobRoleEntity.AppSource) == 0 {
		jobRoleEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return &jobRoleEntity
}

func MapJobRoleUpdateInputToEntity(input *model.JobRoleUpdateInput) *entity.JobRoleEntity {
	if input == nil {
		return nil
	}
	jobRoleEntity := entity.JobRoleEntity{
		Id:            input.ID,
		StartedAt:     input.StartedAt,
		EndedAt:       input.EndedAt,
		JobTitle:      utils.IfNotNilString(input.JobTitle),
		Primary:       utils.IfNotNilBool(input.Primary),
		SourceOfTruth: entity.DataSourceOpenline,
		Description:   input.Description,
		Company:       input.Company,
	}
	return &jobRoleEntity
}

func MapEntityToJobRole(entity *entity.JobRoleEntity) *model.JobRole {
	jobRole := model.JobRole{
		ID:            entity.Id,
		Primary:       entity.Primary,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		Description:   entity.Description,
		Company:       entity.Company,
		AppSource:     entity.AppSource,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		StartedAt:     entity.StartedAt,
		EndedAt:       entity.EndedAt,
	}
	if len(entity.JobTitle) > 0 {
		jobRole.JobTitle = utils.StringPtr(entity.JobTitle)
	}
	return &jobRole
}

func MapEntitiesToJobRoles(entities *entity.JobRoleEntities) []*model.JobRole {
	var jobRoles []*model.JobRole
	for _, jobRoleEntity := range *entities {
		jobRoles = append(jobRoles, MapEntityToJobRole(&jobRoleEntity))
	}
	return jobRoles
}
