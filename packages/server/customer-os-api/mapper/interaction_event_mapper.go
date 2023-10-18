package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToInteractionEvent(entity *entity.InteractionEventEntity) *model.InteractionEvent {
	return &model.InteractionEvent{
		ID:              entity.Id,
		CreatedAt:       *entity.CreatedAt,
		EventIdentifier: utils.StringPtrNillable(entity.EventIdentifier),
		Content:         utils.StringPtrNillable(entity.Content),
		ContentType:     utils.StringPtrNillable(entity.ContentType),
		Channel:         entity.Channel,
		ChannelData:     entity.ChannelData,
		EventType:       entity.EventType,
		Source:          MapDataSourceToModel(entity.Source),
		SourceOfTruth:   MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:       entity.AppSource,
	}
}

func MapEntitiesToInteractionEvents(entities *entity.InteractionEventEntities) []*model.InteractionEvent {
	var interactionEvents []*model.InteractionEvent
	for _, interactionEventEntity := range *entities {
		interactionEvents = append(interactionEvents, MapEntityToInteractionEvent(&interactionEventEntity))
	}
	return interactionEvents
}
func MapInteractionEventInputToEntity(input *model.InteractionEventInput) *entity.InteractionEventEntity {
	return &entity.InteractionEventEntity{
		EventIdentifier:  utils.IfNotNilString(input.EventIdentifier),
		ExternalId:       input.ExternalID,
		ExternalSystemId: input.ExternalSystemID,
		Content:          utils.IfNotNilString(input.Content),
		ContentType:      utils.IfNotNilString(input.ContentType),
		Channel:          input.Channel,
		ChannelData:      input.ChannelData,
		EventType:        input.EventType,
		AppSource:        input.AppSource,
		CreatedAt:        input.CreatedAt,
	}
}
