package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapOrganizationInputToEntity(input *model.OrganizationInput) *entity.OrganizationEntity {
	return &entity.OrganizationEntity{
		Name:          input.Name,
		Description:   utils.IfNotNilString(input.Description),
		Website:       utils.IfNotNilString(input.Website),
		Industry:      utils.IfNotNilString(input.Industry),
		IsPublic:      utils.IfNotNilBool(input.IsPublic),
		Employees:     utils.IfNotNilInt64(input.Employees),
		Market:        MapMarketFromModel(input.Market),
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
}

func MapOrganizationUpdateInputToEntity(input *model.OrganizationUpdateInput) *entity.OrganizationEntity {
	return &entity.OrganizationEntity{
		ID:            input.ID,
		Name:          input.Name,
		Description:   utils.IfNotNilString(input.Description),
		Website:       utils.IfNotNilString(input.Website),
		Industry:      utils.IfNotNilString(input.Industry),
		IsPublic:      utils.IfNotNilBool(input.IsPublic),
		Employees:     utils.IfNotNilInt64(input.Employees),
		Market:        MapMarketFromModel(input.Market),
		SourceOfTruth: entity.DataSourceOpenline,
	}
}

func MapEntityToOrganization(entity *entity.OrganizationEntity) *model.Organization {
	return &model.Organization{
		ID:            entity.ID,
		Name:          entity.Name,
		Description:   utils.StringPtr(entity.Description),
		Website:       utils.StringPtr(entity.Website),
		Industry:      utils.StringPtr(entity.Industry),
		IsPublic:      utils.BoolPtr(entity.IsPublic),
		Employees:     utils.Int64Ptr(entity.Employees),
		Market:        MapMarketToModel(entity.Market),
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		LastTouchPoint: utils.ToPtr(model.LastTouchpoint{
			TimelineEventID: entity.LastTouchpointId,
			At:              entity.LastTouchpointAt,
		}),
	}
}

func MapEntitiesToOrganizations(organizationEntities *entity.OrganizationEntities) []*model.Organization {
	var organizations []*model.Organization
	for _, organizationEntity := range *organizationEntities {
		organizations = append(organizations, MapEntityToOrganization(&organizationEntity))
	}
	return organizations
}
