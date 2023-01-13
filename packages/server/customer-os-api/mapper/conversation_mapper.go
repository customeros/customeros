package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"time"
)

func MapConversationInputToEntity(input model.ConversationInput) *entity.ConversationEntity {
	conversationEntity := entity.ConversationEntity{
		Id:            utils.IfNotNilString(input.ID),
		Channel:       utils.IfNotNilString(input.Channel),
		Status:        MapConversationStatusFromModel(input.Status),
		SourceOfTruth: entity.DataSourceOpenline,
		Source:        entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilString(input.AppSource),
	}
	if input.StartedAt == nil {
		conversationEntity.StartedAt = time.Now().UTC()
	} else {
		conversationEntity.StartedAt = *input.StartedAt
	}
	return &conversationEntity
}

func MapConversationUpdateInputToEntity(input model.ConversationUpdateInput) *entity.ConversationEntity {
	conversationEntity := entity.ConversationEntity{
		Id:            utils.IfNotNilString(input.ID),
		Channel:       utils.IfNotNilString(input.Channel),
		SourceOfTruth: entity.DataSourceOpenline,
	}
	if input.Status != nil {
		conversationEntity.Status = MapConversationStatusFromModel(*input.Status)
	}
	return &conversationEntity
}

func MapEntityToConversation(entity *entity.ConversationEntity) *model.Conversation {
	conversationModel := model.Conversation{
		ID:           entity.Id,
		StartedAt:    entity.StartedAt,
		EndedAt:      entity.EndedAt,
		Status:       MapConversationStatusToModel(entity.Status),
		Channel:      utils.StringPtr(entity.Channel),
		MessageCount: entity.MessageCount,
		Source:       MapDataSourceToModel(entity.Source),
	}
	if len(entity.AppSource) > 0 {
		conversationModel.AppSource = utils.StringPtr(entity.AppSource)
	}
	return &conversationModel
}

func MapEntitiesToConversations(entities *entity.ConversationEntities) []*model.Conversation {
	var conversations []*model.Conversation
	for _, conversationEntity := range *entities {
		conversations = append(conversations, MapEntityToConversation(&conversationEntity))
	}
	return conversations
}
