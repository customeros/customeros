package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToInteractionSession(entity *neo4jentity.InteractionSessionEntity) *model.InteractionSession {
	if entity == nil {
		return nil
	}
	return &model.InteractionSession{
		ID:            entity.Id,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Identifier:    entity.Identifier,
		Name:          entity.Name,
		Status:        entity.Status,
		Type:          &entity.Type,
		Channel:       &entity.Channel,
		ChannelData:   &entity.ChannelData,
		AppSource:     entity.AppSource,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
	}
}

func MapEntitiesToInteractionSessions(entities *neo4jentity.InteractionSessionEntities) []*model.InteractionSession {
	var interactionSessions []*model.InteractionSession
	for _, interactionSessionEntity := range *entities {
		interactionSessions = append(interactionSessions, MapEntityToInteractionSession(&interactionSessionEntity))
	}
	return interactionSessions
}
