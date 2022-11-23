package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapEntityToConversation(entity *entity.ConversationEntity) *model.Conversation {
	return &model.Conversation{
		ID:        entity.Id,
		StartedAt: entity.StartedAt,
		ContactID: entity.ContactId,
		UserID:    entity.UserId,
	}
}

func MapEntitiesToConversations(entities *entity.ConversationEntities) []*model.Conversation {
	var conversations []*model.Conversation
	for _, conversationEntity := range *entities {
		conversations = append(conversations, MapEntityToConversation(&conversationEntity))
	}
	return conversations
}
