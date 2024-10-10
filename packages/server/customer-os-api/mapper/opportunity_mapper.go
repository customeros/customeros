package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
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

func MapOpportunitySaveInputToEntity(input model.OpportunitySaveInput) *repository.OpportunitySaveFields {
	mapped := repository.OpportunitySaveFields{
		AppSource: constants.AppSourceCustomerOsApi,
		Source:    neo4jentity.DataSourceOpenline.String(),
	}

	if input.Name != nil {
		mapped.Name = *input.Name
		mapped.UpdateName = true
	}
	if input.Amount != nil {
		mapped.Amount = *input.Amount
		mapped.UpdateAmount = true
	}
	if input.MaxAmount != nil {
		mapped.MaxAmount = *input.MaxAmount
		mapped.UpdateMaxAmount = true
	}
	if input.ExternalStage != nil {
		mapped.ExternalStage = *input.ExternalStage
		mapped.UpdateExternalStage = true
	}
	if input.ExternalType != nil {
		mapped.ExternalType = *input.ExternalType
		mapped.UpdateExternalType = true
	}
	if input.EstimatedClosedDate != nil {
		mapped.EstimatedClosedAt = input.EstimatedClosedDate
		mapped.UpdateEstimatedClosedAt = true
	}
	if input.InternalStage != nil {
		mapped.InternalStage = MapInternalStageFromModel(*input.InternalStage).String()
		mapped.UpdateInternalStage = true
	}
	if input.InternalType != nil {
		mapped.InternalType = MapInternalStageFromModel(*input.InternalStage).String()
		mapped.UpdateInternalStage = true
	}
	if input.NextSteps != nil {
		mapped.NextSteps = *input.NextSteps
		mapped.UpdateNextSteps = true
	}
	if input.LikelihoodRate != nil {
		mapped.LikelihoodRate = *input.LikelihoodRate
		mapped.UpdateLikelihoodRate = true
	}
	if input.OwnerID != nil {
		mapped.OwnerId = *input.OwnerID
		mapped.UpdateOwnerId = true
	}

	if input.Currency != nil {
		mapped.Currency = mapper.MapCurrencyFromModel(*input.Currency)
		mapped.UpdateCurrency = true
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
