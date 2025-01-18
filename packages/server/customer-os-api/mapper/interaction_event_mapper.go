package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToInteractionEvent(entity *neo4jentity.InteractionEventEntity) *model.InteractionEvent {
	return &model.InteractionEvent{
		ID:              entity.Id,
		CreatedAt:       entity.CreatedAt,
		EventIdentifier: utils.StringPtrNillable(entity.Identifier),
		Content:         utils.StringPtrNillable(entity.Content),
		ContentType:     utils.StringPtrNillable(entity.ContentType),
		Channel:         entity.Channel.String(),
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
