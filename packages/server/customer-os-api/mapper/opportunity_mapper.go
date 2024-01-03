package mapper

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
		RenewalLikelihood:      MapOpportunityRenewalLikelihoodToModel(entity.RenewalLikelihood),
		RenewalUpdatedByUserAt: entity.RenewalUpdatedByUserAt,
		RenewalUpdatedByUserID: entity.RenewalUpdatedByUserId,
		Comments:               entity.Comments,
		Source:                 MapDataSourceToModel(entity.Source),
		SourceOfTruth:          MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:              entity.AppSource,
	}
}

func MapOpportunityUpdateInputToEntity(input model.OpportunityUpdateInput) *entity.OpportunityEntity {
	opportunityEntity := entity.OpportunityEntity{
		Id:                input.OpportunityID,
		Name:              utils.IfNotNilString(input.Name),
		Amount:            utils.IfNotNilFloat64(input.Amount),
		ExternalType:      utils.IfNotNilString(input.ExternalType),
		ExternalStage:     utils.IfNotNilString(input.ExternalStage),
		GeneralNotes:      utils.IfNotNilString(input.GeneralNotes),
		NextSteps:         utils.IfNotNilString(input.NextSteps),
		EstimatedClosedAt: input.EstimatedClosedDate,
		Source:            neo4jentity.DataSourceOpenline,
		SourceOfTruth:     neo4jentity.DataSourceOpenline,
		AppSource:         utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &opportunityEntity
}

func MapEntitiesToOpportunities(entities *entity.OpportunityEntities) []*model.Opportunity {
	var Opportunities []*model.Opportunity
	for _, OpportunityEntity := range *entities {
		Opportunities = append(Opportunities, MapEntityToOpportunity(&OpportunityEntity))
	}
	return Opportunities
}
