package mapper

import (
	localentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"time"
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
			Version:       entity.AggregateVersion,
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
		LogoURL:            utils.StringPtr(entity.LogoUrl),
		Icon:               utils.StringPtr(entity.IconUrl),
		IconURL:            utils.StringPtr(entity.IconUrl),
		IcpFit:             entity.IcpFit,
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
		EnrichDetails:    prepareOrganizationEnrichDetails(entity.EnrichDetails.EnrichRequestedAt, entity.EnrichDetails.EnrichedAt, entity.EnrichDetails.EnrichFailedAt),

		// TODO: All below fields are deprecated and should be removed
		IsPublic:                      utils.BoolPtr(entity.IsPublic),
		Note:                          utils.StringPtr(entity.Note),
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

func prepareOrganizationEnrichDetails(requestedAt, enrichedAt, failedAt *time.Time) *model.EnrichDetails {
	output := model.EnrichDetails{
		RequestedAt: requestedAt,
		EnrichedAt:  enrichedAt,
		FailedAt:    failedAt,
	}
	if enrichedAt == nil && failedAt == nil && requestedAt != nil {
		// if requested is older than 1 min, remove it
		if time.Since(*requestedAt) > time.Minute {
			output.RequestedAt = nil
		}
	}
	return &output
}

func MapOrganizationSaveInputToEntity(input model.OrganizationSaveInput) *repository.OrganizationSaveFields {
	mapped := repository.OrganizationSaveFields{
		SourceFields: neo4jmodel.Source{
			Source:        constants.SourceOpenline,
			SourceOfTruth: constants.SourceOpenline,
			AppSource:     constants.AppSourceCustomerOsApi,
		},

		Domains: input.Domains,
	}

	if input.ReferenceID != nil {
		mapped.ReferenceId = *input.ReferenceID
		mapped.UpdateReferenceId = true
	}
	if input.Name != nil {
		mapped.Name = *input.Name
		mapped.UpdateName = true
	}
	if input.Description != nil {
		mapped.Description = *input.Description
		mapped.UpdateDescription = true
	}
	if input.Website != nil {
		mapped.Website = *input.Website
		mapped.UpdateWebsite = true
	}
	if input.Industry != nil {
		mapped.Industry = *input.Industry
		mapped.UpdateIndustry = true
	}
	if input.SubIndustry != nil {
		mapped.SubIndustry = *input.SubIndustry
		mapped.UpdateSubIndustry = true
	}
	if input.IndustryGroup != nil {
		mapped.IndustryGroup = *input.IndustryGroup
		mapped.UpdateIndustryGroup = true
	}
	if input.Public != nil {
		mapped.IsPublic = *input.Public
		mapped.UpdateIsPublic = true
	}
	if input.Market != nil {
		mapped.Market = MapMarketFromModel(input.Market)
		mapped.UpdateMarket = true
	}
	if input.Employees != nil {
		mapped.Employees = *input.Employees
		mapped.UpdateEmployees = true
	}
	if input.Notes != nil {
		mapped.Note = *input.Notes
		mapped.UpdateNote = true
	}
	if input.TargetAudience != nil {
		mapped.TargetAudience = *input.TargetAudience
		mapped.UpdateTargetAudience = true
	}
	if input.ValueProposition != nil {
		mapped.ValueProposition = *input.ValueProposition
		mapped.UpdateValueProposition = true
	}
	if input.LogoURL != nil {
		mapped.LogoUrl = *input.LogoURL
		mapped.UpdateLogoUrl = true
	}
	if input.IconURL != nil {
		mapped.IconUrl = *input.IconURL
		mapped.UpdateIconUrl = true
	}
	if input.YearFounded != nil {
		mapped.YearFounded = *input.YearFounded
		mapped.UpdateYearFounded = true
	}
	if input.EmployeeGrowthRate != nil {
		mapped.EmployeeGrowthRate = *input.EmployeeGrowthRate
		mapped.UpdateEmployeeGrowthRate = true
	}
	if input.Headquarters != nil {
		mapped.Headquarters = *input.Headquarters
		mapped.UpdateHeadquarters = true
	}
	if input.SlackChannelID != nil {
		mapped.SlackChannelId = *input.SlackChannelID
		mapped.UpdateSlackChannelId = true
	}
	if input.LeadSource != nil {
		mapped.LeadSource = *input.LeadSource
		mapped.UpdateLeadSource = true
	}
	if input.Stage != nil {
		mapped.Stage = mapper.MapStageFromModel(*input.Stage)
		mapped.UpdateStage = true
	}
	if input.Relationship != nil {
		mapped.Relationship = mapper.MapRelationshipFromModel(*input.Relationship)
		mapped.UpdateRelationship = true
	}
	if input.LastFundingRound != nil {
		mapped.LastFundingRound = mapper.MapFundingRoundFromModel(input.LastFundingRound)
		mapped.UpdateLastFundingRound = true
	}
	if input.LastFundingAmount != nil {
		mapped.LastFundingAmount = *input.LastFundingAmount
		mapped.UpdateLastFundingAmount = true
	}
	if input.IcpFit != nil {
		mapped.IcpFit = *input.IcpFit
		mapped.UpdateIcpFit = true
	}

	return &mapped
}
