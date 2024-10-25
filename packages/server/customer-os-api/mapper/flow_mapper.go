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
			OnHold:       entity.OnHold,
			Ready:        entity.Ready,
			Scheduled:    entity.Scheduled,
			InProgress:   entity.InProgress,
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

func MapEntityToFlowSender(entity *neo4jentity.FlowSenderEntity) *model.FlowSender {
	if entity == nil {
		return nil
	}
	return &model.FlowSender{
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

func MapEntitiesToFlowSenders(entities *neo4jentity.FlowSenderEntities) []*model.FlowSender {
	var mapped []*model.FlowSender
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowSender(&entity))
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

func MapFlowActionMergeInputToEntity(input model.FlowSenderMergeInput) *neo4jentity.FlowSenderEntity {
	return &neo4jentity.FlowSenderEntity{
		Id:     utils.StringOrEmpty(input.ID),
		UserId: input.UserID,
	}
}
