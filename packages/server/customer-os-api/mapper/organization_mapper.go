package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToOrganization(entity *entity.OrganizationEntity) *model.Organization {
	if entity == nil {
		return nil
	}
	return &model.Organization{
		ID:                            entity.ID,
		ReferenceID:                   utils.StringPtrNillable(entity.ReferenceId),
		CustomerOsID:                  entity.CustomerOsId,
		Name:                          entity.Name,
		Description:                   utils.StringPtr(entity.Description),
		Note:                          utils.StringPtr(entity.Note),
		Website:                       utils.StringPtr(entity.Website),
		Industry:                      utils.StringPtr(entity.Industry),
		SubIndustry:                   utils.StringPtr(entity.SubIndustry),
		IndustryGroup:                 utils.StringPtr(entity.IndustryGroup),
		TargetAudience:                utils.StringPtr(entity.TargetAudience),
		ValueProposition:              utils.StringPtr(entity.ValueProposition),
		IsPublic:                      utils.BoolPtr(entity.IsPublic),
		IsCustomer:                    utils.BoolPtr(entity.IsCustomer),
		Employees:                     utils.Int64Ptr(entity.Employees),
		Market:                        MapMarketToModel(entity.Market),
		LastFundingRound:              MapFundingRoundToModel(entity.LastFundingRound),
		LastFundingAmount:             utils.StringPtr(entity.LastFundingAmount),
		CreatedAt:                     entity.CreatedAt,
		UpdatedAt:                     entity.UpdatedAt,
		Source:                        MapDataSourceToModel(entity.Source),
		SourceOfTruth:                 MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:                     entity.AppSource,
		LastTouchPointTimelineEventID: entity.LastTouchpointId,
		LastTouchPointAt:              entity.LastTouchpointAt,
		LastTouchPointType:            MapLastTouchpointTypeToModel(entity.LastTouchpointType),
		AccountDetails: &model.OrgAccountDetails{
			RenewalLikelihood: &model.RenewalLikelihood{
				Probability:         MapRenewalLikelihoodToModel(entity.RenewalLikelihood.RenewalLikelihood),
				PreviousProbability: MapRenewalLikelihoodToModel(entity.RenewalLikelihood.PreviousRenewalLikelihood),
				Comment:             entity.RenewalLikelihood.Comment,
				UpdatedAt:           entity.RenewalLikelihood.UpdatedAt,
				UpdatedByID:         entity.RenewalLikelihood.UpdatedBy,
			},
			RenewalForecast: &model.RenewalForecast{
				Amount:          entity.RenewalForecast.Amount,
				PotentialAmount: entity.RenewalForecast.PotentialAmount,
				Comment:         entity.RenewalForecast.Comment,
				UpdatedAt:       entity.RenewalForecast.UpdatedAt,
				UpdatedByID:     entity.RenewalForecast.UpdatedById,
				Arr:             entity.RenewalForecast.Arr,
				MaxArr:          entity.RenewalForecast.MaxArr,
			},
			BillingDetails: &model.BillingDetails{
				Amount:            entity.BillingDetails.Amount,
				Frequency:         MapRenewalCycleToModel(entity.BillingDetails.Frequency),
				RenewalCycle:      MapRenewalCycleToModel(entity.BillingDetails.RenewalCycle),
				RenewalCycleStart: entity.BillingDetails.RenewalCycleStart,
				RenewalCycleNext:  entity.BillingDetails.RenewalCycleNext,
			},
			RenewalSummary: &model.RenewalSummary{
				ArrForecast:       entity.RenewalSummary.ArrForecast,
				MaxArrForecast:    entity.RenewalSummary.MaxArrForecast,
				NextRenewalDate:   entity.RenewalSummary.NextRenewalAt,
				RenewalLikelihood: MapOpportunityRenewalLikelihoodToModelPtr(entity.RenewalSummary.RenewalLikelihood),
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
