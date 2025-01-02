package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToOpportunity(entity *neo4jentity.OpportunityEntity) *model.Opportunity {
	if entity == nil {
		return nil
	}
	return &model.Opportunity{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
		},
		Name:                   entity.Name,
		Amount:                 entity.Amount,
		MaxAmount:              entity.MaxAmount,
		InternalType:           MapInternalTypeToModel(entity.InternalType),
		ExternalType:           entity.ExternalType,
		InternalStage:          MapInternalStageToModel(entity.InternalStage),
		ExternalStage:          entity.ExternalStage,
		EstimatedClosedAt:      entity.EstimatedClosedAt,
		GeneralNotes:           entity.GeneralNotes,
		NextSteps:              entity.NextSteps,
		RenewedAt:              entity.RenewalDetails.RenewedAt,
		RenewalLikelihood:      MapOpportunityRenewalLikelihoodToModel(entity.RenewalDetails.RenewalLikelihood),
		RenewalUpdatedByUserAt: entity.RenewalDetails.RenewalUpdatedByUserAt,
		RenewalUpdatedByUserID: entity.RenewalDetails.RenewalUpdatedByUserId,
		RenewalApproved:        entity.RenewalDetails.RenewalApproved,
		RenewalAdjustedRate:    entity.RenewalDetails.RenewalAdjustedRate,
		Comments:               entity.Comments,
		ID:                     entity.Id,
		Currency:               utils.ToPtr(mapper.MapCurrencyToModel(entity.Currency)),
		LikelihoodRate:         entity.LikelihoodRate,
		StageLastUpdated:       entity.StageUpdatedAt,
	}
}

func MapOpportunitySaveInputToEntity(input model.OpportunitySaveInput) *data_fields.OpportunityFields {
	mapped := data_fields.OpportunityFields{
		AppSource:         utils.StringPtr(constants.AppSourceCustomerOsApi),
		Source:            utils.StringPtr(neo4jentity.DataSourceOpenline.String()),
		Name:              input.Name,
		Amount:            input.Amount,
		MaxAmount:         input.MaxAmount,
		ExternalStage:     input.ExternalStage,
		ExternalType:      input.ExternalType,
		EstimatedClosedAt: input.EstimatedClosedDate,
		NextSteps:         input.NextSteps,
		LikelihoodRate:    input.LikelihoodRate,
		OwnerId:           input.OwnerID,
	}
	if input.InternalStage != nil {
		mapped.InternalStage = utils.StringPtr(MapInternalStageFromModel(*input.InternalStage).String())
	}
	if input.InternalType != nil {
		mapped.InternalType = utils.StringPtr(MapInternalTypeFromModel(*input.InternalType).String())
	}
	if input.Currency != nil {
		mapped.Currency = utils.ToPtr(mapper.MapCurrencyFromModel(*input.Currency))
	}

	return &mapped
}

func MapEntitiesToOpportunities(entities *neo4jentity.OpportunityEntities) []*model.Opportunity {
	var Opportunities []*model.Opportunity
	for _, OpportunityEntity := range *entities {
		Opportunities = append(Opportunities, MapEntityToOpportunity(&OpportunityEntity))
	}
	return Opportunities
}
