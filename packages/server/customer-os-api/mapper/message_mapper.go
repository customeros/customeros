package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"time"
)

func MapMessageInputToEntity(input model.MessageInput) *entity.MessageEntity {
	messageEntity := new(entity.MessageEntity)
	messageEntity.Id = input.ID
	messageEntity.ConversationId = input.ConversationID
	messageEntity.Channel = MapMessageChannelFromModel(input.Channel)
	if input.StartedAt == nil {
		messageEntity.StartedAt = time.Now().UTC()
	} else {
		messageEntity.StartedAt = *input.StartedAt
	}
	return messageEntity
}

func MapEntityToMessage(entity *entity.MessageEntity) *model.Message {
	return &model.Message{
		ID:        entity.Id,
		StartedAt: entity.StartedAt,
		Channel:   MapMessageChannelToModel(entity.Channel),
	}
}

func MapEntityToMessageAction(entity *entity.MessageEntity) *model.MessageAction {
	return &model.MessageAction{
		ID:             entity.Id,
		StartedAt:      entity.StartedAt,
		ConversationID: entity.ConversationId,
		Channel:        MapMessageChannelToModel(entity.Channel),
	}
}
