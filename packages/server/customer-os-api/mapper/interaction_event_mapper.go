package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToInteractionEvent(entity *entity.InteractionEventEntity) *model.InteractionEvent {
	return &model.InteractionEvent{
		ID:              entity.Id,
		CreatedAt:       entity.CreatedAt,
		EventIdentifier: utils.StringPtrNillable(entity.EventIdentifier),
		Content:         utils.StringPtrNillable(entity.Content),
		ContentType:     utils.StringPtrNillable(entity.ContentType),
		Channel:         utils.StringPtrNillable(entity.Channel),
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
