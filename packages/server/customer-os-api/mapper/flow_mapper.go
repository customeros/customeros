package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
		Status:      entity.Status,
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
		Status:      entity.Status,
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
		Name:   entity.Name,
		Status: entity.Status,
		Action: entity.Action,
	}
}

func MapEntitiesToFlowSequenceSteps(entities *neo4jentity.FlowSequenceStepEntities) []*model.FlowSequenceStep {
	var mapped []*model.FlowSequenceStep
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowSequenceStep(&entity))
	}
	return mapped
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

func MapFlowSequenceStepMergeInputToEntity(input model.FlowSequenceStepMergeInput) *neo4jentity.FlowSequenceStepEntity {
	return &neo4jentity.FlowSequenceStepEntity{
		Id:     utils.StringOrEmpty(input.ID),
		Name:   input.Name,
		Action: input.Action,
	}
}
