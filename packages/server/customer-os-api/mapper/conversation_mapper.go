package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"time"
)

func MapConversationInputToEntity(input model.ConversationInput) *entity.ConversationEntity {
	conversationEntity := entity.ConversationEntity{
		Id:      utils.IfNotNilString(input.ID),
		Channel: utils.IfNotNilString(input.Channel),
		Status:  MapConversationStatusFromModel(input.Status),
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
		Id:      utils.IfNotNilString(input.ID),
		Channel: utils.IfNotNilString(input.Channel),
	}
	if input.Status != nil {
		conversationEntity.Status = MapConversationStatusFromModel(*input.Status)
	}
	return &conversationEntity
}

func MapEntityToConversation(entity *entity.ConversationEntity) *model.Conversation {
	return &model.Conversation{
		ID:           entity.Id,
		StartedAt:    entity.StartedAt,
		EndedAt:      entity.EndedAt,
		Status:       MapConversationStatusToModel(entity.Status),
		Channel:      utils.StringPtr(entity.Channel),
		MessageCount: entity.MessageCount,
	}
}

func MapEntitiesToConversations(entities *entity.ConversationEntities) []*model.Conversation {
	var conversations []*model.Conversation
	for _, conversationEntity := range *entities {
		conversations = append(conversations, MapEntityToConversation(&conversationEntity))
	}
	return conversations
}
