package mapper

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToContactUIDetails(entity *neo4jentity.ContactEntity, output *model.ContactUIDetails) {
	if entity == nil || output == nil {
		return
	}

	output.ID = entity.Id
	output.CreatedAt = entity.CreatedAt
	output.UpdatedAt = entity.UpdatedAt
	output.Hide = entity.Hide
	output.FirstName = entity.FirstName
	output.LastName = entity.LastName
	output.Name = entity.Name
	output.Prefix = entity.Prefix
	output.Description = entity.Description
	output.Timezone = entity.Timezone
	output.ProfilePhotoURL = entity.ProfilePhotoUrl

	mapEnrichDetails(entity.EnrichDetails, output)
}

func mapEnrichDetails(enrichDetails neo4jentity.ContactEnrichDetails, output *model.ContactUIDetails) {
	output.EnrichedRequestedAt = enrichDetails.EnrichRequestedAt
	output.EnrichedAt = enrichDetails.EnrichedAt
	output.EnrichedFailedAt = enrichDetails.EnrichFailedAt

	output.EnrichedEmailRequestedAt = enrichDetails.FindWorkEmailWithBetterContactRequestedAt
	output.EnrichedEmailEnrichedAt = enrichDetails.FindWorkEmailWithBetterContactCompletedAt
	output.EnrichedEmailFound = enrichDetails.FindWorkEmailWithBetterContactFound
	if enrichDetails.EnrichedAt == nil && enrichDetails.EnrichFailedAt == nil && enrichDetails.EnrichRequestedAt != nil {
		// if requested is older than 1 min, remove it
		if time.Since(*enrichDetails.EnrichRequestedAt) > time.Minute {
			output.EnrichedRequestedAt = nil
		}
	}
	if output.EnrichedEmailEnrichedAt == nil && output.EnrichedEmailFound == nil && output.EnrichedEmailRequestedAt != nil {
		// if requested is older than 10 min, set found to false
		if time.Since(*output.EnrichedEmailRequestedAt) > 10*time.Minute {
			output.EnrichedEmailRequestedAt = nil
			output.EnrichedEmailFound = utils.BoolPtr(false)
		}
	}
}
