package mapper

import (
	localentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToOrganization(entity *neo4jentity.OrganizationEntity) *model.Organization {
	if entity == nil {
		return nil
	}
	organization := model.Organization{
		Metadata: &model.Metadata{
			ID:            entity.ID,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
		},
		CustomID:           utils.StringPtrNillable(entity.ReferenceId),
		CustomerOsID:       entity.CustomerOsId,
		Name:               entity.Name,
		Description:        utils.StringPtr(entity.Description),
		Website:            utils.StringPtr(entity.Website),
		Industry:           utils.StringPtr(entity.Industry),
		SubIndustry:        utils.StringPtr(entity.SubIndustry),
		IndustryGroup:      utils.StringPtr(entity.IndustryGroup),
		TargetAudience:     utils.StringPtr(entity.TargetAudience),
		ValueProposition:   utils.StringPtr(entity.ValueProposition),
		Public:             utils.BoolPtr(entity.IsPublic),
		Employees:          utils.Int64Ptr(entity.Employees),
		Market:             MapMarketToModel(entity.Market),
		LastFundingRound:   mapper.MapFundingRoundToModel(entity.LastFundingRound),
		LastFundingAmount:  utils.StringPtr(entity.LastFundingAmount),
		YearFounded:        entity.YearFounded,
		Headquarters:       utils.StringPtr(entity.Headquarters),
		EmployeeGrowthRate: utils.StringPtr(entity.EmployeeGrowthRate),
		SlackChannelID:     utils.StringPtr(entity.SlackChannelId),
		Logo:               utils.StringPtr(entity.LogoUrl),
		Icon:               utils.StringPtr(entity.IconUrl),
		AccountDetails: &model.OrgAccountDetails{
			RenewalSummary: &model.RenewalSummary{
				ArrForecast:       entity.RenewalSummary.ArrForecast,
				MaxArrForecast:    entity.RenewalSummary.MaxArrForecast,
				NextRenewalDate:   entity.RenewalSummary.NextRenewalAt,
				RenewalLikelihood: MapOpportunityRenewalLikelihoodToModelPtr(entity.RenewalSummary.RenewalLikelihood),
			},
			Onboarding: &model.OnboardingDetails{
				Status:    MapOnboardingStatusToModel(localentity.GetOnboardingStatus(entity.OnboardingDetails.Status)),
				UpdatedAt: entity.OnboardingDetails.UpdatedAt,
				Comments:  utils.StringPtr(entity.OnboardingDetails.Comments),
			},
			Churned:     entity.DerivedData.ChurnedAt,
			Ltv:         utils.Float64Ptr(entity.DerivedData.Ltv),
			LtvCurrency: utils.ToPtr(mapper.MapCurrencyToModel(entity.DerivedData.LtvCurrency)),
		},
		LastTouchpoint: &model.LastTouchpoint{
			LastTouchPointTimelineEventID: entity.LastTouchpointId,
			LastTouchPointAt:              entity.LastTouchpointAt,
			LastTouchPointType:            mapper.MapLastTouchpointTypeToModel(entity.LastTouchpointType),
		},
		Hide:             entity.Hide,
		Notes:            utils.StringPtr(entity.Note),
		Stage:            utils.ToPtr(mapper.MapStageToModel(entity.Stage)),
		Relationship:     utils.ToPtr(mapper.MapRelationshipToModel(entity.Relationship)),
		LeadSource:       utils.StringPtr(entity.LeadSource),
		StageLastUpdated: entity.StageUpdatedAt,

		// TODO: All below fields are deprecated and should be removed
		IsPublic:                      utils.BoolPtr(entity.IsPublic),
		Note:                          utils.StringPtr(entity.Note),
		LogoURL:                       utils.StringPtr(entity.LogoUrl),
		ID:                            entity.ID,
		ReferenceID:                   utils.StringPtrNillable(entity.ReferenceId),
		Source:                        MapDataSourceToModel(entity.Source),
		SourceOfTruth:                 MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:                     entity.AppSource,
		CreatedAt:                     entity.CreatedAt,
		UpdatedAt:                     entity.UpdatedAt,
		LastTouchPointTimelineEventID: entity.LastTouchpointId,
		LastTouchPointAt:              entity.LastTouchpointAt,
		LastTouchPointType:            mapper.MapLastTouchpointTypeToModel(entity.LastTouchpointType),
	}

	if organization.Relationship != nil && *organization.Relationship == model.OrganizationRelationshipCustomer {
		organization.IsCustomer = utils.BoolPtr(true)
	} else {
		organization.IsCustomer = utils.BoolPtr(false)
	}

	return &organization
}

func MapEntitiesToOrganizations(organizationEntities *neo4jentity.OrganizationEntities) []*model.Organization {
	var organizations []*model.Organization
	for _, organizationEntity := range *organizationEntities {
		organizations = append(organizations, MapEntityToOrganization(&organizationEntity))
	}
	return organizations
}
