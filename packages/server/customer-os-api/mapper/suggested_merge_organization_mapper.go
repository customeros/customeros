package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToSuggestedMergeOrganization(entity *neo4jentity.OrganizationEntity) *model.SuggestedMergeOrganization {
	return &model.SuggestedMergeOrganization{
		Organization: MapEntityToOrganization(entity),
		SuggestedBy:  entity.SuggestedMerge.SuggestedBy,
		SuggestedAt:  entity.SuggestedMerge.SuggestedAt,
		Confidence:   entity.SuggestedMerge.Confidence,
	}
}

func MapEntitiesToSuggestedMergeOrganizations(organizationEntities *neo4jentity.OrganizationEntities) []*model.SuggestedMergeOrganization {
	var suggestedMergeOrganizations []*model.SuggestedMergeOrganization
	for _, organizationEntity := range *organizationEntities {
		suggestedMergeOrganizations = append(suggestedMergeOrganizations, MapEntityToSuggestedMergeOrganization(&organizationEntity))
	}
	return suggestedMergeOrganizations
}
