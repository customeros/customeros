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
		ID:                input.ID,
		Name:              input.Name,
		Description:       utils.IfNotNilString(input.Description),
		Website:           utils.IfNotNilString(input.Website),
		Industry:          utils.IfNotNilString(input.Industry),
		SubIndustry:       utils.IfNotNilString(input.SubIndustry),
		IndustryGroup:     utils.IfNotNilString(input.IndustryGroup),
		IsPublic:          utils.IfNotNilBool(input.IsPublic),
		Employees:         utils.IfNotNilInt64(input.Employees),
		Market:            MapMarketFromModel(input.Market),
		TargetAudience:    utils.IfNotNilString(input.TargetAudience),
		ValueProposition:  utils.IfNotNilString(input.ValueProposition),
		LastFundingRound:  MapFundingRoundFromModel(input.LastFundingRound),
		LastFundingAmount: utils.IfNotNilString(input.LastFundingAmount),
		SlackChannelLink:  utils.IfNotNilString(input.SlackChannelLink),
		SourceOfTruth:     entity.DataSourceOpenline,
	}
}

func MapEntityToOrganization(entity *entity.OrganizationEntity) *model.Organization {
	if entity == nil {
		return nil
	}
	return &model.Organization{
		ID:                            entity.ID,
		Name:                          entity.Name,
		Description:                   utils.StringPtr(entity.Description),
		Website:                       utils.StringPtr(entity.Website),
		Industry:                      utils.StringPtr(entity.Industry),
		SubIndustry:                   utils.StringPtr(entity.SubIndustry),
		IndustryGroup:                 utils.StringPtr(entity.IndustryGroup),
		TargetAudience:                utils.StringPtr(entity.TargetAudience),
		ValueProposition:              utils.StringPtr(entity.ValueProposition),
		IsPublic:                      utils.BoolPtr(entity.IsPublic),
		Employees:                     utils.Int64Ptr(entity.Employees),
		Market:                        MapMarketToModel(entity.Market),
		LastFundingRound:              MapFundingRoundToModel(entity.LastFundingRound),
		LastFundingAmount:             utils.StringPtr(entity.LastFundingAmount),
		SlackChannelLink:              utils.StringPtr(entity.SlackChannelLink),
		CreatedAt:                     entity.CreatedAt,
		UpdatedAt:                     entity.UpdatedAt,
		Source:                        MapDataSourceToModel(entity.Source),
		SourceOfTruth:                 MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:                     entity.AppSource,
		LastTouchPointAt:              entity.LastTouchpointAt,
		LastTouchPointTimelineEventID: entity.LastTouchpointId,
		AccountDetails: &model.OrgAccountDetails{
			RenewalLikelihood: &model.RenewalLikelihood{
				Probability:         MapRenewalLikelihoodToModel(entity.RenewalLikelihood.RenewalLikelihood),
				PreviousProbability: MapRenewalLikelihoodToModel(entity.RenewalLikelihood.PreviousRenewalLikelihood),
				Comment:             entity.RenewalLikelihood.Comment,
				UpdatedAt:           entity.RenewalLikelihood.UpdatedAt,
				UpdatedBy:           entity.RenewalLikelihood.UpdatedBy,
			},
			RenewalForecast: &model.RenewalForecast{
				Amount:         entity.RenewalForecast.Amount,
				PreviousAmount: entity.RenewalForecast.PreviousAmount,
				Comment:        entity.RenewalForecast.Comment,
				UpdatedAt:      entity.RenewalForecast.UpdatedAt,
				UpdatedBy:      entity.RenewalForecast.UpdatedBy,
			},
			BillingDetails: &model.BillingDetails{
				Amount:            entity.BillingDetails.Amount,
				Frequency:         MapRenewalCycleToModel(entity.BillingDetails.Frequency),
				RenewalCycle:      MapRenewalCycleToModel(entity.BillingDetails.RenewalCycle),
				RenewalCycleStart: entity.BillingDetails.RenewalCycleStart,
			},
		},
	}
}

func MapEntitiesToOrganizations(organizationEntities *entity.OrganizationEntities) []*model.Organization {
	var organizations []*model.Organization
	for _, organizationEntity := range *organizationEntities {
		organizations = append(organizations, MapEntityToOrganization(&organizationEntity))
	}
	return organizations
}
