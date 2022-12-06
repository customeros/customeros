package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"time"
)

func MapMessageInputToEntity(input model.MessageInput) *entity.MessageEntity {
	messageEntity := new(entity.MessageEntity)
	messageEntity.Id = input.ID
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
