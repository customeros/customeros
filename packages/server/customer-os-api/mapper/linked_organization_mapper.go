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
	var linkedOrgainzations []*model.LinkedOrganization
	for _, organizationEntity := range *organizationEntities {
		linkedOrgainzations = append(linkedOrgainzations, MapEntityToLinkedOrganization(&organizationEntity))
	}
	return linkedOrgainzations
}
