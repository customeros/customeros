package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToLinkedOrganization(entity *neo4jentity.OrganizationEntity) *model.LinkedOrganization {
	return &model.LinkedOrganization{
		Organization: MapEntityToOrganization(entity),
		Type:         entity.LinkedOrganizationType,
	}
}

func MapEntitiesToLinkedOrganizations(organizationEntities *neo4jentity.OrganizationEntities) []*model.LinkedOrganization {
	var linkedOrganizations []*model.LinkedOrganization
	for _, organizationEntity := range *organizationEntities {
		linkedOrganizations = append(linkedOrganizations, MapEntityToLinkedOrganization(&organizationEntity))
	}
	return linkedOrganizations
}
