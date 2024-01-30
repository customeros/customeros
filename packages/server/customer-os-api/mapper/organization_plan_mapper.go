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
		MasterPlanID:  entity.MasterPlanId,
	}
}

func MapEntitiesToOrganizationPlans(entities *neo4jentity.OrganizationPlanEntities) []*model.OrganizationPlan {
	var models []*model.OrganizationPlan
	if len(*entities) == 0 {
		return models
	}
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

func MapEntityToOrganizationPlanMilestoneItems(entities []neo4jentity.OrganizationPlanMilestoneItem) []*model.OrganizationPlanMilestoneItem {
	var models []*model.OrganizationPlanMilestoneItem
	for _, entity := range entities {
		models = append(models, MapEntityToOrganizationPlanMilestoneItem(&entity))
	}
	return models
}

func MapEntityToOrganizationPlanMilestoneItem(entity *neo4jentity.OrganizationPlanMilestoneItem) *model.OrganizationPlanMilestoneItem {
	if entity == nil {
		return nil
	}
	status := model.OnboardingPlanMilestoneItemStatusNotDone
	switch entity.Status {
	case model.OnboardingPlanMilestoneItemStatusNotDone.String():
		status = model.OnboardingPlanMilestoneItemStatusNotDone
	case model.OnboardingPlanMilestoneItemStatusDone.String():
		status = model.OnboardingPlanMilestoneItemStatusDone
	case model.OnboardingPlanMilestoneItemStatusSkipped.String():
		status = model.OnboardingPlanMilestoneItemStatusSkipped
	case model.OnboardingPlanMilestoneItemStatusNotDoneLate.String():
		status = model.OnboardingPlanMilestoneItemStatusNotDoneLate
	case model.OnboardingPlanMilestoneItemStatusSkippedLate.String():
		status = model.OnboardingPlanMilestoneItemStatusSkippedLate
	case model.OnboardingPlanMilestoneItemStatusDoneLate.String():
		status = model.OnboardingPlanMilestoneItemStatusDoneLate
	}
	return &model.OrganizationPlanMilestoneItem{
		Status:    status,
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
