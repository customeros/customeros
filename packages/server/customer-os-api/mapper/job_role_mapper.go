package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapJobRoleInputToEntity(input *model.JobRoleInput) *neo4jentity.JobRoleEntity {
	if input == nil {
		return nil
	}
	jobRoleEntity := neo4jentity.JobRoleEntity{
		JobTitle:    utils.IfNotNilString(input.JobTitle),
		Primary:     utils.IfNotNilBool(input.Primary),
		Description: input.Description,
		Company:     input.Company,
		StartedAt:   input.StartedAt,
		EndedAt:     input.EndedAt,
		Source:      neo4jentity.DataSourceOpenline,
		AppSource:   utils.IfNotNilString(input.AppSource),
	}
	if len(jobRoleEntity.AppSource) == 0 {
		jobRoleEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return &jobRoleEntity
}

func MapJobRoleUpdateInputToEntity(input *model.JobRoleUpdateInput) *neo4jentity.JobRoleEntity {
	if input == nil {
		return nil
	}
	jobRoleEntity := neo4jentity.JobRoleEntity{
		Id:          input.ID,
		StartedAt:   input.StartedAt,
		EndedAt:     input.EndedAt,
		JobTitle:    utils.IfNotNilString(input.JobTitle),
		Primary:     utils.IfNotNilBool(input.Primary),
		Description: input.Description,
		Company:     input.Company,
	}
	return &jobRoleEntity
}

func MapEntityToJobRole(entity *neo4jentity.JobRoleEntity) *model.JobRole {
	jobRole := model.JobRole{
		ID:          entity.Id,
		Primary:     entity.Primary,
		Source:      MapDataSourceToModel(entity.Source),
		Description: entity.Description,
		Company:     entity.Company,
		AppSource:   entity.AppSource,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		StartedAt:   entity.StartedAt,
		EndedAt:     entity.EndedAt,
	}
	if len(entity.JobTitle) > 0 {
		jobRole.JobTitle = utils.StringPtr(entity.JobTitle)
	}
	return &jobRole
}

func MapEntitiesToJobRoles(entities *neo4jentity.JobRoleEntities) []*model.JobRole {
	var jobRoles []*model.JobRole
	for _, jobRoleEntity := range *entities {
		jobRoles = append(jobRoles, MapEntityToJobRole(&jobRoleEntity))
	}
	return jobRoles
}
