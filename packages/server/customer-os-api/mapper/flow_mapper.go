package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

func MapEntityToFlow(entity *postgresEntity.Flow) *model.Flow {
	if entity == nil {
		return nil
	}
	return &model.Flow{
		Metadata: &model.Metadata{
			ID:            entity.ID,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        model.DataSourceOpenline,
			SourceOfTruth: model.DataSourceOpenline,
			AppSource:     "",
		},
		Name:        entity.Name,
		Description: entity.Description,
		Status:      mapFlowStatus(entity.Status),
	}
}

func MapEntityToFlowSequence(entity *postgresEntity.FlowSequence) *model.FlowSequence {
	if entity == nil {
		return nil
	}
	return &model.FlowSequence{
		Metadata: &model.Metadata{
			ID:            entity.ID,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        model.DataSourceOpenline,
			SourceOfTruth: model.DataSourceOpenline,
			AppSource:     "",
		},
		Name:        entity.Name,
		Description: entity.Description,
		Status:      mapFlowSequenceStatus(entity.Status),
	}
}

func MapEntitiesToFlowSequence(entities []*postgresEntity.FlowSequence) []*model.FlowSequence {
	var mapped []*model.FlowSequence
	for _, entity := range entities {
		mapped = append(mapped, MapEntityToFlowSequence(entity))
	}
	return mapped
}

func mapFlowStatus(flowSequenceStatus postgresEntity.FlowStatus) model.FlowStatus {
	switch flowSequenceStatus {
	case postgresEntity.FlowStatusActive:
		return model.FlowStatusActive
	case postgresEntity.FlowStatusInactive:
		return model.FlowStatusInactive
	case postgresEntity.FlowStatusPaused:
		return model.FlowStatusPaused
	case postgresEntity.FlowStatusArchived:
		return model.FlowStatusArchived
	default:
		return model.FlowStatusActive
	}
}

func mapFlowSequenceStatus(flowSequenceStatus postgresEntity.FlowSequenceStatus) model.FlowSequenceStatus {
	switch flowSequenceStatus {
	case postgresEntity.FlowSequenceStatusActive:
		return model.FlowSequenceStatusActive
	case postgresEntity.FlowSequenceStatusInactive:
		return model.FlowSequenceStatusInactive
	case postgresEntity.FlowSequenceStatusPaused:
		return model.FlowSequenceStatusPaused
	case postgresEntity.FlowSequenceStatusArchived:
		return model.FlowSequenceStatusArchived
	default:
		return model.FlowSequenceStatusActive
	}
}
