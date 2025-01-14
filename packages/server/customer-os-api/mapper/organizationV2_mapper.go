package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	enummapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
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
	output.ValueProposition = utils.StringPtr(entity.ValueProposition)
	output.Public = utils.BoolPtr(entity.IsPublic)
	output.Employees = utils.Int64Ptr(entity.Employees)
	output.Market = MapMarketToModel(entity.Market)
	output.LastFundingRound = enummapper.MapFundingRoundToModel(entity.LastFundingRound)
	output.YearFounded = entity.YearFounded
	output.SlackChannelID = utils.StringPtr(entity.SlackChannelId)
	output.LogoURL = utils.StringPtr(entity.LogoUrl)
	output.IconURL = utils.StringPtr(entity.IconUrl)
	output.Notes = utils.StringPtr(entity.Note)
	output.Stage = utils.ToPtr(enummapper.MapStageToModel(entity.Stage))
	output.Relationship = utils.ToPtr(enummapper.MapRelationshipToModel(entity.Relationship))
	output.LeadSource = utils.StringPtr(entity.LeadSource)

	output.Ltv = utils.Float64Ptr(entity.DerivedData.Ltv)
	output.ChurnedAt = entity.DerivedData.ChurnedAt
	output.ContactCount = utils.IntPtr(int(entity.DerivedData.ContactCount))

	output.RenewalSummaryArrForecast = entity.RenewalSummary.ArrForecast
	output.RenewalSummaryMaxArrForecast = entity.RenewalSummary.MaxArrForecast
	output.RenewalSummaryNextRenewalAt = entity.RenewalSummary.NextRenewalAt
	output.RenewalSummaryRenewalLikelihood = MapOpportunityRenewalLikelihoodToModelPtr(entity.RenewalSummary.RenewalLikelihood)

	output.OnboardingStatus = enummapper.MapOnboardingStatusToModel(neo4jenum.DecodeOnboardingStatus(entity.OnboardingDetails.Status))
	output.OnboardingStatusUpdatedAt = entity.OnboardingDetails.UpdatedAt
	output.OnboardingComments = utils.StringPtr(entity.OnboardingDetails.Comments)

	output.LastTouchPointAt = entity.LastTouchpointAt
	output.LastTouchPointType = enummapper.MapLastTouchpointTypeToModel(entity.LastTouchpointType)

	output.EnrichedAt = entity.EnrichDetails.EnrichedAt
	output.EnrichedRequestedAt = entity.EnrichDetails.EnrichRequestedAt
	output.EnrichedFailedAt = entity.EnrichDetails.EnrichFailedAt
	if entity.EnrichDetails.EnrichedAt == nil && entity.EnrichDetails.EnrichFailedAt == nil && entity.EnrichDetails.EnrichRequestedAt != nil {
		// if requested is older than 1 min, remove it
		if time.Since(*entity.EnrichDetails.EnrichRequestedAt) > time.Minute {
			output.EnrichedRequestedAt = nil
		}
	}
}
