package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToInteractionEvent(entity *neo4jentity.InteractionEventEntity) *model.InteractionEvent {
	return &model.InteractionEvent{
		ID:              entity.Id,
		CreatedAt:       entity.CreatedAt,
		EventIdentifier: utils.StringPtrNillable(entity.Identifier),
		Content:         utils.StringPtrNillable(entity.Content),
		ContentType:     utils.StringPtrNillable(entity.ContentType),
		Channel:         entity.Channel,
		ChannelData:     &entity.ChannelData,
		EventType:       &entity.EventType,
		Source:          MapDataSourceToModel(entity.Source),
		SourceOfTruth:   MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:       entity.AppSource,
	}
}

func MapEntitiesToInteractionEvents(entities *neo4jentity.InteractionEventEntities) []*model.InteractionEvent {
	var interactionEvents []*model.InteractionEvent
	for _, interactionEventEntity := range *entities {
		interactionEvents = append(interactionEvents, MapEntityToInteractionEvent(&interactionEventEntity))
	}
	return interactionEvents
}
