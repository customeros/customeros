package mapper

import (
	localentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToOrganizationUIDetails(entity *neo4jentity.OrganizationEntity, output *model.OrganizationUIDetails) {
	if entity == nil || output == nil {
		return
	}

	output.ID = entity.ID
	output.CreatedAt = entity.CreatedAt
	output.UpdatedAt = entity.UpdatedAt
	output.Hide = entity.Hide
	output.CustomerOsID = entity.CustomerOsId
	output.ReferenceID = entity.ReferenceId
	output.Name = entity.Name
	output.Description = utils.StringPtr(entity.Description)
	output.Website = utils.StringPtr(entity.Website)
	output.Industry = utils.StringPtr(entity.Industry)
	output.ValueProposition = utils.StringPtr(entity.ValueProposition)
	output.Public = utils.BoolPtr(entity.IsPublic)
	output.Employees = utils.Int64Ptr(entity.Employees)
	output.Market = MapMarketToModel(entity.Market)
	output.LastFundingRound = mapper.MapFundingRoundToModel(entity.LastFundingRound)
	output.YearFounded = entity.YearFounded
	output.SlackChannelID = utils.StringPtr(entity.SlackChannelId)
	output.LogoURL = utils.StringPtr(entity.LogoUrl)
	output.IconURL = utils.StringPtr(entity.IconUrl)
	output.Notes = utils.StringPtr(entity.Note)
	output.Stage = utils.ToPtr(mapper.MapStageToModel(entity.Stage))
	output.Relationship = utils.ToPtr(mapper.MapRelationshipToModel(entity.Relationship))
	output.LeadSource = utils.StringPtr(entity.LeadSource)

	output.Ltv = utils.Float64Ptr(entity.DerivedData.Ltv)
	output.ChurnedAt = entity.DerivedData.ChurnedAt

	output.RenewalSummaryArrForecast = entity.RenewalSummary.ArrForecast
	output.RenewalSummaryMaxArrForecast = entity.RenewalSummary.MaxArrForecast
	output.RenewalSummaryNextRenewalAt = entity.RenewalSummary.NextRenewalAt
	output.RenewalSummaryRenewalLikelihood = MapOpportunityRenewalLikelihoodToModelPtr(entity.RenewalSummary.RenewalLikelihood)

	output.OnboardingStatus = MapOnboardingStatusToModel(localentity.GetOnboardingStatus(entity.OnboardingDetails.Status))
	output.OnboardingStatusUpdatedAt = entity.OnboardingDetails.UpdatedAt
	output.OnboardingComments = utils.StringPtr(entity.OnboardingDetails.Comments)

	output.LastTouchPointAt = entity.LastTouchpointAt
	output.LastTouchPointType = mapper.MapLastTouchpointTypeToModel(entity.LastTouchpointType)

	output.EnrichedAt = entity.EnrichDetails.EnrichedAt
	output.EnrichedRequestedAt = entity.EnrichDetails.EnrichRequestedAt
	output.EnrichedFailedAt = entity.EnrichDetails.EnrichFailedAt
}

//func MapEntitiesToOrganizationsV2(organizationEntities *neo4jentity.OrganizationEntities) []*model.OrganizationV2 {
//	var organizations []*model.OrganizationV2
//	for _, organizationEntity := range *organizationEntities {
//
//		organizations = append(organizations, MapEntityToOrganizationV2(&organizationEntity))
//	}
//	return organizations
//}
