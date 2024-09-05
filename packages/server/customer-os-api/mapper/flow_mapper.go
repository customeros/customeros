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
		Status:      MapFlowStatusToModel(entity.Status),
	}
}

func MapEntitiesToFlows(entities []*postgresEntity.Flow) []*model.Flow {
	var mapped []*model.Flow
	for _, entity := range entities {
		mapped = append(mapped, MapEntityToFlow(entity))
	}
	return mapped
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
		Status:      MapFlowSequenceStatusToModel(entity.Status),
	}
}

func MapEntitiesToFlowSequence(entities []*postgresEntity.FlowSequence) []*model.FlowSequence {
	var mapped []*model.FlowSequence
	for _, entity := range entities {
		mapped = append(mapped, MapEntityToFlowSequence(entity))
	}
	return mapped
}

func MapEntityToFlowSequenceStep(entity *postgresEntity.FlowSequenceStep) *model.FlowSequenceStep {
	if entity == nil {
		return nil
	}
	return &model.FlowSequenceStep{
		Metadata: &model.Metadata{
			ID:            entity.ID,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        model.DataSourceOpenline,
			SourceOfTruth: model.DataSourceOpenline,
			AppSource:     "",
		},
		Name:    entity.Name,
		Status:  MapFlowSequenceStepStatusToModel(entity.Status),
		Type:    MapFlowSequenceStepTypeToModel(entity.Type),
		Subtype: MapFlowSequenceStepSubtypeToModel(entity.Subtype),
		Body:    entity.Body,
	}
}

func MapEntitiesToFlowSequenceSteps(entities []*postgresEntity.FlowSequenceStep) []*model.FlowSequenceStep {
	var mapped []*model.FlowSequenceStep
	for _, entity := range entities {
		mapped = append(mapped, MapEntityToFlowSequenceStep(entity))
	}
	return mapped
}

func MapFlowStatusToModel(flowSequenceStatus postgresEntity.FlowStatus) model.FlowStatus {
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
		return ""
	}
}

func MapFlowStatusToEntity(flowSequenceStatus model.FlowStatus) postgresEntity.FlowStatus {
	switch flowSequenceStatus {
	case model.FlowStatusActive:
		return postgresEntity.FlowStatusActive
	case model.FlowStatusInactive:
		return postgresEntity.FlowStatusInactive
	case model.FlowStatusPaused:
		return postgresEntity.FlowStatusPaused
	case model.FlowStatusArchived:
		return postgresEntity.FlowStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStatusToModel(flowSequenceStatus postgresEntity.FlowSequenceStatus) model.FlowSequenceStatus {
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
		return ""
	}
}

func MapFlowSequenceStatusToEntity(flowSequenceStatus model.FlowSequenceStatus) postgresEntity.FlowSequenceStatus {
	switch flowSequenceStatus {
	case model.FlowSequenceStatusActive:
		return postgresEntity.FlowSequenceStatusActive
	case model.FlowSequenceStatusInactive:
		return postgresEntity.FlowSequenceStatusInactive
	case model.FlowSequenceStatusPaused:
		return postgresEntity.FlowSequenceStatusPaused
	case model.FlowSequenceStatusArchived:
		return postgresEntity.FlowSequenceStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStepStatusToModel(flowSequenceStepStatus postgresEntity.FlowSequenceStepStatus) model.FlowSequenceStepStatus {
	switch flowSequenceStepStatus {
	case postgresEntity.FlowSequenceStepStatusActive:
		return model.FlowSequenceStepStatusActive
	case postgresEntity.FlowSequenceStepStatusInactive:
		return model.FlowSequenceStepStatusInactive
	case postgresEntity.FlowSequenceStepStatusPaused:
		return model.FlowSequenceStepStatusPaused
	case postgresEntity.FlowSequenceStepStatusArchived:
		return model.FlowSequenceStepStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStepStatusToEntity(enum model.FlowSequenceStepStatus) postgresEntity.FlowSequenceStepStatus {
	switch enum {
	case model.FlowSequenceStepStatusActive:
		return postgresEntity.FlowSequenceStepStatusActive
	case model.FlowSequenceStepStatusInactive:
		return postgresEntity.FlowSequenceStepStatusInactive
	case model.FlowSequenceStepStatusPaused:
		return postgresEntity.FlowSequenceStepStatusPaused
	case model.FlowSequenceStepStatusArchived:
		return postgresEntity.FlowSequenceStepStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStepTypeToModel(enum postgresEntity.FlowSequenceStepType) model.FlowSequenceStepType {
	switch enum {
	case postgresEntity.FlowSequenceStepTypeEmail:
		return model.FlowSequenceStepTypeEmail
	case postgresEntity.FlowSequenceStepTypeLinkedin:
		return model.FlowSequenceStepTypeLinkedin
	default:
		return ""
	}
}

func MapFlowSequenceStepSubtypeToModel(enum *postgresEntity.FlowSequenceStepSubtype) *model.FlowSequenceStepSubtype {
	if enum == nil {
		return nil
	}
	var v model.FlowSequenceStepSubtype
	switch *enum {
	case postgresEntity.FlowSequenceStepSubtypeLinkedinConnectionRequest:
		v = model.FlowSequenceStepSubtypeLinkedinConnectionRequest
	case postgresEntity.FlowSequenceStepSubtypeLinkedinMessage:
		v = model.FlowSequenceStepSubtypeLinkedinMessage
	default:
		v = ""
	}

	return &v
}
