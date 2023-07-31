package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToSuggestedMergeOrganization(entity *entity.OrganizationEntity) *model.SuggestedMergeOrganization {
	return &model.SuggestedMergeOrganization{
		Organization: MapEntityToOrganization(entity),
		SuggestedBy:  entity.SuggestedMerge.SuggestedBy,
		SuggestedAt:  entity.SuggestedMerge.SuggestedAt,
		Confidence:   entity.SuggestedMerge.Confidence,
	}
}

func MapEntitiesToSuggestedMergeOrganizations(organizationEntities *entity.OrganizationEntities) []*model.SuggestedMergeOrganization {
	var suggestedMergeOrganizations []*model.SuggestedMergeOrganization
	for _, organizationEntity := range *organizationEntities {
		suggestedMergeOrganizations = append(suggestedMergeOrganizations, MapEntityToSuggestedMergeOrganization(&organizationEntity))
	}
	return suggestedMergeOrganizations
}
