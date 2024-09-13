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

func MapEntitiesToFlows(entities *neo4jentity.FlowEntities) []*model.Flow {
	var mapped []*model.Flow
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlow(&entity))
	}
	return mapped
}

func MapEntityToFlowContact(entity *neo4jentity.FlowContactEntity) *model.FlowContact {
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
	}
}

func MapEntitiesToFlowContacts(entities *neo4jentity.FlowContactEntities) []*model.FlowContact {
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
		Mailbox:       entity.Mailbox,
		EmailsPerHour: entity.EmailsPerHour,
	}
}

func MapEntitiesToFlowActionSenders(entities *neo4jentity.FlowActionSenderEntities) []*model.FlowActionSender {
	var mapped []*model.FlowActionSender
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowActionSender(&entity))
	}
	return mapped
}

func MapEntityToFlowAction(entity *neo4jentity.FlowActionEntity) *model.FlowAction {
	if entity == nil {
		return nil
	}
	var actionData model.FlowActionData

	if entity.ActionType == neo4jentity.FlowActionTypeWait {
		actionData = &model.FlowActionDataWait{
			Minutes: entity.ActionData.Wait.Minutes,
		}
	}
	if entity.ActionType == neo4jentity.FlowActionTypeEmailNew {
		actionData = &model.FlowActionDataEmail{
			Subject:      entity.ActionData.EmailNew.Subject,
			BodyTemplate: entity.ActionData.EmailNew.BodyTemplate,
		}
	}
	if entity.ActionType == neo4jentity.FlowActionTypeEmailReply {
		actionData = &model.FlowActionDataEmail{
			ReplyToID:    entity.ActionData.EmailReply.ReplyToId,
			Subject:      entity.ActionData.EmailReply.Subject,
			BodyTemplate: entity.ActionData.EmailReply.BodyTemplate,
		}
	}
	if entity.ActionType == neo4jentity.FlowActionTypeLinkedinConnectionRequest {
		actionData = &model.FlowActionLinkedinMessage{
			MessageTemplate: entity.ActionData.LinkedinConnectionRequest.MessageTemplate,
		}
	}
	if entity.ActionType == neo4jentity.FlowActionTypeLinkedinMessage {
		actionData = &model.FlowActionLinkedinMessage{
			MessageTemplate: entity.ActionData.LinkedinMessage.MessageTemplate,
		}
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
		Index:      entity.Index,
		Name:       entity.Name,
		Status:     entity.Status,
		ActionType: entity.ActionType,
		ActionData: actionData,
	}
}

func MapEntitiesToFlowActions(entities *neo4jentity.FlowActionEntities) []*model.FlowAction {
	var mapped []*model.FlowAction
	for _, entity := range *entities {
		mapped = append(mapped, MapEntityToFlowAction(&entity))
	}
	return mapped
}

func MapFlowMergeInputToEntity(input model.FlowMergeInput) *neo4jentity.FlowEntity {
	return &neo4jentity.FlowEntity{
		Id:          utils.StringOrEmpty(input.ID),
		Name:        input.Name,
		Description: input.Description,
	}
}

func MapFlowActionMergeInputToEntity(input model.FlowActionMergeInput) *neo4jentity.FlowActionEntity {
	actionData := neo4jentity.FlowActionData{}

	if input.ActionType == neo4jentity.FlowActionTypeWait {
		actionData.Wait = &neo4jentity.FlowActionDataWait{
			Minutes: input.ActionData.Wait.Minutes,
		}
	}
	if input.ActionType == neo4jentity.FlowActionTypeEmailNew {
		actionData.EmailNew = &neo4jentity.FlowActionDataEmail{
			Subject:      input.ActionData.EmailNew.Subject,
			BodyTemplate: input.ActionData.EmailNew.BodyTemplate,
		}
	}
	if input.ActionType == neo4jentity.FlowActionTypeEmailReply {
		actionData.EmailReply = &neo4jentity.FlowActionDataEmail{
			ReplyToId:    input.ActionData.EmailReply.ReplyToID,
			Subject:      input.ActionData.EmailReply.Subject,
			BodyTemplate: input.ActionData.EmailReply.BodyTemplate,
		}
	}
	if input.ActionType == neo4jentity.FlowActionTypeLinkedinConnectionRequest {
		actionData.LinkedinConnectionRequest = &neo4jentity.FlowActionDataLinkedinConnectionRequest{
			MessageTemplate: input.ActionData.LinkedinConnectionRequest.MessageTemplate,
		}
	}
	if input.ActionType == neo4jentity.FlowActionTypeLinkedinMessage {
		actionData.LinkedinMessage = &neo4jentity.FlowActionDataLinkedinMessage{
			MessageTemplate: input.ActionData.LinkedinMessage.MessageTemplate,
		}
	}

	return &neo4jentity.FlowActionEntity{
		Id:         utils.StringOrEmpty(input.ID),
		Name:       input.Name,
		ActionType: input.ActionType,
		ActionData: actionData,
	}
}

func MapFlowActionSenderMergeInputToEntity(input model.FlowActionSenderMergeInput) *neo4jentity.FlowActionSenderEntity {
	return &neo4jentity.FlowActionSenderEntity{
		Id:            utils.StringOrEmpty(input.ID),
		Mailbox:       input.Mailbox,
		EmailsPerHour: input.EmailsPerHour,
		UserId:        input.UserID,
	}
}
