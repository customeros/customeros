package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
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

func MapOpportunityUpdateInputToEntity(input model.OpportunityUpdateInput) *neo4jentity.OpportunityEntity {
	opportunityEntity := neo4jentity.OpportunityEntity{
		Id:                input.OpportunityID,
		Name:              utils.IfNotNilString(input.Name),
		Amount:            utils.IfNotNilFloat64(input.Amount),
		ExternalType:      utils.IfNotNilString(input.ExternalType),
		ExternalStage:     utils.IfNotNilString(input.ExternalStage),
		EstimatedClosedAt: input.EstimatedClosedDate,
		Source:            neo4jentity.DataSourceOpenline,
		SourceOfTruth:     neo4jentity.DataSourceOpenline,
		AppSource:         constants.AppSourceCustomerOsApi,
	}
	return &opportunityEntity
}

func MapEntitiesToOpportunities(entities *neo4jentity.OpportunityEntities) []*model.Opportunity {
	var Opportunities []*model.Opportunity
	for _, OpportunityEntity := range *entities {
		Opportunities = append(Opportunities, MapEntityToOpportunity(&OpportunityEntity))
	}
	return Opportunities
}
