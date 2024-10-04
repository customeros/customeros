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
		Name:   entity.Name,
		Nodes:  entity.Nodes,
		Edges:  entity.Edges,
		Status: entity.Status,
		Statistics: &model.FlowStatistics{
			Total:        entity.Total,
			Pending:      entity.Pending,
			Completed:    entity.Completed,
			GoalAchieved: entity.GoalAchieved,
		},
	}
}

func MapEntitiesToFlows(entities *neo4jentity.FlowEntities) []*model.Flow {
	var mapped []*model.Flow
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlow(&entity))
	}
	return mapped
}

func MapEntityToFlowAction(entity *neo4jentity.FlowActionEntity) *model.FlowAction {
	if entity == nil {
		return nil
	}
	return &model.FlowAction{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        model.DataSourceOpenline,
			SourceOfTruth: model.DataSourceOpenline,
			AppSource:     "",
		},
		JSON: entity.Json,
	}
}

func MapEntitiesToFlowActions(entities *neo4jentity.FlowActionEntities) []*model.FlowAction {
	var mapped []*model.FlowAction
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowAction(&entity))
	}
	return mapped
}

func MapEntityToFlowContact(entity *neo4jentity.FlowParticipantEntity) *model.FlowContact {
	if entity == nil {
		return nil
	}
	return &model.FlowContact{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        model.DataSourceOpenline,
			SourceOfTruth: model.DataSourceOpenline,
			AppSource:     "",
		},
		Status:          entity.Status,
		ScheduledAction: entity.ScheduledAction,
		ScheduledAt:     entity.ScheduledAt,
	}
}

func MapEntitiesToFlowContacts(entities *neo4jentity.FlowParticipantEntities) []*model.FlowContact {
	var mapped []*model.FlowContact
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowContact(&entity))
	}
	return mapped
}

func MapEntityToFlowActionSender(entity *neo4jentity.FlowActionSenderEntity) *model.FlowActionSender {
	if entity == nil {
		return nil
	}
	return &model.FlowActionSender{
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

func MapEntitiesToFlowActionSenders(entities *neo4jentity.FlowActionSenderEntities) []*model.FlowActionSender {
	var mapped []*model.FlowActionSender
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowActionSender(&entity))
	}
	return mapped
}

func MapFlowMergeInputToEntity(input model.FlowMergeInput) *neo4jentity.FlowEntity {
	return &neo4jentity.FlowEntity{
		Id:    utils.StringOrEmpty(input.ID),
		Name:  input.Name,
		Nodes: input.Nodes,
		Edges: input.Edges,
	}
}

func MapFlowActionSenderMergeInputToEntity(input model.FlowActionSenderMergeInput) *neo4jentity.FlowActionSenderEntity {
	return &neo4jentity.FlowActionSenderEntity{
		Id:      utils.StringOrEmpty(input.ID),
		Mailbox: input.Mailbox,
		UserId:  input.UserID,
	}
}
