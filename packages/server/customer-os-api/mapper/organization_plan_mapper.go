package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToOrganizationPlan(entity *neo4jentity.OrganizationPlanEntity) *model.OrganizationPlan {
	if entity == nil {
		return nil
	}
	return &model.OrganizationPlan{
		ID:            entity.Id,
		Name:          entity.Name,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		Retired:       entity.Retired,
		StatusDetails: MapEntityToOrganizationPlanStatusDetails(&entity.StatusDetails),
	}
}

func MapEntitiesToOrganizationPlans(entities *neo4jentity.OrganizationPlanEntities) []*model.OrganizationPlan {
	var models []*model.OrganizationPlan
	for _, entity := range *entities {
		models = append(models, MapEntityToOrganizationPlan(&entity))
	}
	return models
}

func MapEntityToOrganizationPlanMilestone(entity *neo4jentity.OrganizationPlanMilestoneEntity) *model.OrganizationPlanMilestone {
	if entity == nil {
		return nil
	}
	return &model.OrganizationPlanMilestone{
		ID:            entity.Id,
		Name:          entity.Name,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		Order:         entity.Order,
		DueDate:       entity.DueDate,
		Optional:      entity.Optional,
		Items:         MapEntityToOrganizationPlanMilestoneItems(entity.Items),
		Retired:       entity.Retired,
		StatusDetails: MapEntityToOrganizationPlanMilestoneStatusDetails(&entity.StatusDetails),
	}
}

func MapEntitiesToOrganizationPlanMilestones(entities *neo4jentity.OrganizationPlanMilestoneEntities) []*model.OrganizationPlanMilestone {
	var models []*model.OrganizationPlanMilestone
	for _, entity := range *entities {
		models = append(models, MapEntityToOrganizationPlanMilestone(&entity))
	}
	return models
}

func MapEntityToOrganizationPlanMilestoneItems(entities []neo4jentity.OrganizationPlanMilestoneItem) []*model.MilestoneItem {
	var models []*model.MilestoneItem
	for _, entity := range entities {
		models = append(models, MapEntityToOrganizationPlanMilestoneItem(&entity))
	}
	return models
}

func MapEntityToOrganizationPlanMilestoneItem(entity *neo4jentity.OrganizationPlanMilestoneItem) *model.MilestoneItem {
	if entity == nil {
		return nil
	}
	return &model.MilestoneItem{
		Status:    entity.Status,
		UpdatedAt: entity.UpdatedAt,
		Text:      entity.Text,
	}
}

func MapEntityToOrganizationPlanStatusDetails(entity *neo4jentity.OrganizationPlanStatusDetails) *model.StatusDetails {
	if entity == nil {
		return nil
	}
	return &model.StatusDetails{
		Status:    entity.Status,
		UpdatedAt: entity.UpdatedAt,
		Text:      entity.Comments,
	}
}

func MapEntityToOrganizationPlanMilestoneStatusDetails(entity *neo4jentity.OrganizationPlanMilestoneStatusDetails) *model.StatusDetails {
	if entity == nil {
		return nil
	}
	return &model.StatusDetails{
		Status:    entity.Status,
		UpdatedAt: entity.UpdatedAt,
		Text:      entity.Comments,
	}
}
