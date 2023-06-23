package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToHealthIndicator(entity *entity.HealthIndicatorEntity) *model.HealthIndicator {
	if entity == nil {
		return nil
	}
	return &model.HealthIndicator{
		ID:    entity.Id,
		Name:  entity.Name,
		Order: entity.Order,
	}
}

func MapEntitiesToHealthIndicators(entities *entity.HealthIndicatorEntities) []*model.HealthIndicator {
	var healthIndicators []*model.HealthIndicator
	for _, entity := range *entities {
		healthIndicators = append(healthIndicators, MapEntityToHealthIndicator(&entity))
	}
	return healthIndicators
}
