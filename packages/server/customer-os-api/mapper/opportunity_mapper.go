package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToOpportunity(entity *entity.OpportunityEntity) *model.Opportunity {
	if entity == nil {
		return nil
	}
	return &model.Opportunity{
		ID:                     entity.Id,
		Name:                   entity.Name,
		CreatedAt:              entity.CreatedAt,
		UpdatedAt:              entity.UpdatedAt,
		Amount:                 entity.Amount,
		MaxAmount:              entity.MaxAmount,
		InternalType:           MapInternalTypeToModel(entity.InternalType),
		ExternalType:           entity.ExternalType,
		InternalStage:          MapInternalStageToModel(entity.InternalStage),
		ExternalStage:          entity.ExternalStage,
		EstimatedClosedAt:      entity.EstimatedClosedAt,
		GeneralNotes:           entity.GeneralNotes,
		NextSteps:              entity.NextSteps,
		RenewedAt:              entity.RenewedAt,
		RenewalLikelihood:      entity.RenewalLikelihood,
		RenewalUpdatedByUserAt: entity.RenewalUpdatedByUserAt,
		RenewalUpdatedByUserID: entity.RenewalUpdatedByUserId,
		Comments:               entity.Comments,
		Source:                 MapDataSourceToModel(entity.Source),
		SourceOfTruth:          MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:              entity.AppSource,
	}
}
func MapEntitiesToOpportunities(entities *entity.OpportunityEntities) []*model.Opportunity {
	var Opportunities []*model.Opportunity
	for _, OpportunityEntity := range *entities {
		Opportunities = append(Opportunities, MapEntityToOpportunity(&OpportunityEntity))
	}
	return Opportunities
}
