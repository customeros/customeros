package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapOrganizationInputToEntity(input *model.OrganizationInput) *entity.OrganizationEntity {
	return &entity.OrganizationEntity{
		Name:        input.Name,
		Description: utils.IfNotNilString(input.Description),
		Domain:      utils.IfNotNilString(input.Domain),
		Website:     utils.IfNotNilString(input.Website),
		Industry:    utils.IfNotNilString(input.Industry),
		IsPublic:    utils.IfNotNilBool(input.IsPublic),
	}
}

func MapEntityToOrganization(entity *entity.OrganizationEntity) *model.Organization {
	return &model.Organization{
		ID:          entity.Id,
		Name:        entity.Name,
		Description: utils.StringPtr(entity.Description),
		Domain:      utils.StringPtr(entity.Domain),
		Website:     utils.StringPtr(entity.Website),
		Industry:    utils.StringPtr(entity.Industry),
		IsPublic:    utils.BoolPtr(entity.IsPublic),
		CreatedAt:   entity.CreatedAt,
		Readonly:    utils.BoolPtr(entity.Readonly),
	}
}

func MapEntitiesToOrganizations(organizationEntities *entity.OrganizationEntities) []*model.Organization {
	var organizations []*model.Organization
	for _, organizationEntity := range *organizationEntities {
		organizations = append(organizations, MapEntityToOrganization(&organizationEntity))
	}
	return organizations
}
