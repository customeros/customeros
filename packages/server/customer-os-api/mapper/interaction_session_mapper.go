package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
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
		Status:        entity.Status.String(),
		Type:          utils.StringPtr(entity.Type.String()),
		Channel:       utils.StringPtr(entity.Channel.String()),
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
