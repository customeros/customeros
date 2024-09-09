package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToFlow(entity *neo4jentity.FlowEntity) *model.Flow {
	if entity == nil {
		return nil
	}
	return &model.Flow{
		Metadata: &model.Metadata{
			ID:            entity.Id,
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

func MapEntitiesToFlows(entities []*neo4jentity.FlowEntity) []*model.Flow {
	var mapped []*model.Flow
	for _, entity := range entities {
		mapped = append(mapped, MapEntityToFlow(entity))
	}
	return mapped
}

func MapEntityToFlowSequence(entity *neo4jentity.FlowSequenceEntity) *model.FlowSequence {
	if entity == nil {
		return nil
	}
	return &model.FlowSequence{
		Metadata: &model.Metadata{
			ID:            entity.Id,
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

func MapEntitiesToFlowSequence(entities *neo4jentity.FlowSequenceEntities) []*model.FlowSequence {
	var mapped []*model.FlowSequence
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowSequence(&entity))
	}
	return mapped
}

func MapEntityToFlowSequenceContact(entity *neo4jentity.FlowSequenceContactEntity) *model.FlowSequenceContact {
	if entity == nil {
		return nil
	}
	return &model.FlowSequenceContact{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        model.DataSourceOpenline,
			SourceOfTruth: model.DataSourceOpenline,
			AppSource:     "",
		},
	}
}

func MapEntitiesToFlowSequenceContacts(entities *neo4jentity.FlowSequenceContactEntities) []*model.FlowSequenceContact {
	var mapped []*model.FlowSequenceContact
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowSequenceContact(&entity))
	}
	return mapped
}

func MapEntityToFlowSequenceSender(entity *neo4jentity.FlowSequenceSenderEntity) *model.FlowSequenceSender {
	if entity == nil {
		return nil
	}
	return &model.FlowSequenceSender{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        model.DataSourceOpenline,
			SourceOfTruth: model.DataSourceOpenline,
			AppSource:     "",
		},
		Mailbox: entity.Mailbox,
	}
}

func MapEntitiesToFlowSequenceSenders(entities *neo4jentity.FlowSequenceSenderEntities) []*model.FlowSequenceSender {
	var mapped []*model.FlowSequenceSender
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowSequenceSender(&entity))
	}
	return mapped
}

func MapEntityToFlowSequenceStep(entity *neo4jentity.FlowSequenceStepEntity) *model.FlowSequenceStep {
	if entity == nil {
		return nil
	}
	return &model.FlowSequenceStep{
		Metadata: &model.Metadata{
			ID:            entity.Id,
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

func MapEntitiesToFlowSequenceSteps(entities *neo4jentity.FlowSequenceStepEntities) []*model.FlowSequenceStep {
	var mapped []*model.FlowSequenceStep
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowSequenceStep(&entity))
	}
	return mapped
}

func MapFlowStatusToModel(flowSequenceStatus neo4jentity.FlowStatus) model.FlowStatus {
	switch flowSequenceStatus {
	case neo4jentity.FlowStatusActive:
		return model.FlowStatusActive
	case neo4jentity.FlowStatusInactive:
		return model.FlowStatusInactive
	case neo4jentity.FlowStatusPaused:
		return model.FlowStatusPaused
	case neo4jentity.FlowStatusArchived:
		return model.FlowStatusArchived
	default:
		return ""
	}
}

func MapFlowStatusToEntity(flowSequenceStatus model.FlowStatus) neo4jentity.FlowStatus {
	switch flowSequenceStatus {
	case model.FlowStatusActive:
		return neo4jentity.FlowStatusActive
	case model.FlowStatusInactive:
		return neo4jentity.FlowStatusInactive
	case model.FlowStatusPaused:
		return neo4jentity.FlowStatusPaused
	case model.FlowStatusArchived:
		return neo4jentity.FlowStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStatusToModel(flowSequenceStatus neo4jentity.FlowSequenceStatus) model.FlowSequenceStatus {
	switch flowSequenceStatus {
	case neo4jentity.FlowSequenceStatusActive:
		return model.FlowSequenceStatusActive
	case neo4jentity.FlowSequenceStatusInactive:
		return model.FlowSequenceStatusInactive
	case neo4jentity.FlowSequenceStatusPaused:
		return model.FlowSequenceStatusPaused
	case neo4jentity.FlowSequenceStatusArchived:
		return model.FlowSequenceStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStatusToEntity(flowSequenceStatus model.FlowSequenceStatus) neo4jentity.FlowSequenceStatus {
	switch flowSequenceStatus {
	case model.FlowSequenceStatusActive:
		return neo4jentity.FlowSequenceStatusActive
	case model.FlowSequenceStatusInactive:
		return neo4jentity.FlowSequenceStatusInactive
	case model.FlowSequenceStatusPaused:
		return neo4jentity.FlowSequenceStatusPaused
	case model.FlowSequenceStatusArchived:
		return neo4jentity.FlowSequenceStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStepStatusToModel(flowSequenceStepStatus neo4jentity.FlowSequenceStepStatus) model.FlowSequenceStepStatus {
	switch flowSequenceStepStatus {
	case neo4jentity.FlowSequenceStepStatusActive:
		return model.FlowSequenceStepStatusActive
	case neo4jentity.FlowSequenceStepStatusInactive:
		return model.FlowSequenceStepStatusInactive
	case neo4jentity.FlowSequenceStepStatusPaused:
		return model.FlowSequenceStepStatusPaused
	case neo4jentity.FlowSequenceStepStatusArchived:
		return model.FlowSequenceStepStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStepStatusToEntity(enum model.FlowSequenceStepStatus) neo4jentity.FlowSequenceStepStatus {
	switch enum {
	case model.FlowSequenceStepStatusActive:
		return neo4jentity.FlowSequenceStepStatusActive
	case model.FlowSequenceStepStatusInactive:
		return neo4jentity.FlowSequenceStepStatusInactive
	case model.FlowSequenceStepStatusPaused:
		return neo4jentity.FlowSequenceStepStatusPaused
	case model.FlowSequenceStepStatusArchived:
		return neo4jentity.FlowSequenceStepStatusArchived
	default:
		return ""
	}
}

func MapFlowSequenceStepTypeToModel(enum neo4jentity.FlowSequenceStepType) model.FlowSequenceStepType {
	switch enum {
	case neo4jentity.FlowSequenceStepTypeEmail:
		return model.FlowSequenceStepTypeEmail
	case neo4jentity.FlowSequenceStepTypeLinkedin:
		return model.FlowSequenceStepTypeLinkedin
	default:
		return ""
	}
}

func MapFlowSequenceStepSubtypeToModel(enum *neo4jentity.FlowSequenceStepSubtype) *model.FlowSequenceStepSubtype {
	if enum == nil {
		return nil
	}
	var v model.FlowSequenceStepSubtype
	switch *enum {
	case neo4jentity.FlowSequenceStepSubtypeLinkedinConnectionRequest:
		v = model.FlowSequenceStepSubtypeLinkedinConnectionRequest
	case neo4jentity.FlowSequenceStepSubtypeLinkedinMessage:
		v = model.FlowSequenceStepSubtypeLinkedinMessage
	default:
		v = ""
	}

	return &v
}

func MapFlowSequenceCreateInputToEntity(input model.FlowSequenceCreateInput) *neo4jentity.FlowSequenceEntity {
	return &neo4jentity.FlowSequenceEntity{
		Name:        input.Name,
		Description: input.Description,
	}
}

func MapFlowSequenceUpdateInputToEntity(input model.FlowSequenceUpdateInput) *neo4jentity.FlowSequenceEntity {
	return &neo4jentity.FlowSequenceEntity{
		Id:          input.ID,
		Name:        input.Name,
		Description: input.Description,
	}
}

func MapFlowSequenceStepCreateInputToEntity(input model.FlowSequenceStepCreateInput) *neo4jentity.FlowSequenceStepEntity {
	return &neo4jentity.FlowSequenceStepEntity{
		Name: input.Name,
	}
}

func MapFlowSequenceStepUpdateInputToEntity(input model.FlowSequenceStepUpdateInput) *neo4jentity.FlowSequenceStepEntity {
	return &neo4jentity.FlowSequenceStepEntity{
		Id:   input.ID,
		Name: input.Name,
	}
}
