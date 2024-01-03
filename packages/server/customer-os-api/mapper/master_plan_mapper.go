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
	}
}

func MapEntitiesToMasterPlans(entities *neo4jentity.MasterPlanEntities) []*model.MasterPlan {
	var masterPlans []*model.MasterPlan
	for _, masterPlanEntity := range *entities {
		masterPlans = append(masterPlans, MapEntityToMasterPlan(&masterPlanEntity))
	}
	return masterPlans
}
