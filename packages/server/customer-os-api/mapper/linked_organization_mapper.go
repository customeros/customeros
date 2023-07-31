package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToLinkedOrganization(entity *entity.OrganizationEntity) *model.LinkedOrganization {
	return &model.LinkedOrganization{
		Organization: MapEntityToOrganization(entity),
		Type:         entity.LinkedOrganizationType,
	}
}

func MapEntitiesToLinkedOrganizations(organizationEntities *entity.OrganizationEntities) []*model.LinkedOrganization {
	var linkedOrganizations []*model.LinkedOrganization
	for _, organizationEntity := range *organizationEntities {
		linkedOrganizations = append(linkedOrganizations, MapEntityToLinkedOrganization(&organizationEntity))
	}
	return linkedOrganizations
}
