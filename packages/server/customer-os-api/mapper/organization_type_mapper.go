package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapOrganizationTypeInputToEntity(input model.OrganizationTypeInput) *entity.OrganizationTypeEntity {
	organizationTypeEntity := entity.OrganizationTypeEntity{
		Name: input.Name,
	}
	return &organizationTypeEntity
}

func MapOrganizationTypeUpdateInputToEntity(input model.OrganizationTypeUpdateInput) *entity.OrganizationTypeEntity {
	organizationTypeEntity := entity.OrganizationTypeEntity{
		Id:   input.ID,
		Name: input.Name,
	}
	return &organizationTypeEntity
}

func MapEntityToOrganizationType(entity *entity.OrganizationTypeEntity) *model.OrganizationType {
	return &model.OrganizationType{
		ID:        entity.Id,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
	}
}

func MapEntitiesToOrganizationTypes(entities *entity.OrganizationTypeEntities) []*model.OrganizationType {
	var organizationTypes []*model.OrganizationType
	for _, organizationTypeEntity := range *entities {
		organizationTypes = append(organizationTypes, MapEntityToOrganizationType(&organizationTypeEntity))
	}
	return organizationTypes
}
