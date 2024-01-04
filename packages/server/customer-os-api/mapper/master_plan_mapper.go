package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToMasterPlan(entity *neo4jentity.MasterPlanEntity) *model.MasterPlan {
	if entity == nil {
		return nil
	}
	return &model.MasterPlan{
		ID:            entity.Id,
		Name:          entity.Name,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		Retired:       entity.Retired,
	}
}

func MapEntitiesToMasterPlans(entities *neo4jentity.MasterPlanEntities) []*model.MasterPlan {
	var models []*model.MasterPlan
	for _, entity := range *entities {
		models = append(models, MapEntityToMasterPlan(&entity))
	}
	return models
}

func MapEntityToMasterPlanMilestone(entity *neo4jentity.MasterPlanMilestoneEntity) *model.MasterPlanMilestone {
	if entity == nil {
		return nil
	}
	return &model.MasterPlanMilestone{
		ID:            entity.Id,
		Name:          entity.Name,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		Order:         entity.Order,
		DurationHours: entity.DurationHours,
		Optional:      entity.Optional,
		Items:         entity.Items,
		Retired:       entity.Retired,
	}
}

func MapEntitiesToMasterPlanMilestones(entities *neo4jentity.MasterPlanMilestoneEntities) []*model.MasterPlanMilestone {
	var models []*model.MasterPlanMilestone
	for _, entity := range *entities {
		models = append(models, MapEntityToMasterPlanMilestone(&entity))
	}
	return models
}
