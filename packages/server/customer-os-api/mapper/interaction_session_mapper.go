package mapper

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToInteractionSession(entity *entity.InteractionSessionEntity) *model.InteractionSession {
	if entity == nil {
		return nil
	}
	return &model.InteractionSession{
		ID:                entity.Id,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
		StartedAt:         entity.CreatedAt,
		EndedAt:           utils.TimePtr(entity.UpdatedAt),
		SessionIdentifier: entity.SessionIdentifier,
		Name:              entity.Name,
		Status:            entity.Status,
		Type:              entity.Type,
		Channel:           entity.Channel,
		ChannelData:       entity.ChannelData,
		AppSource:         entity.AppSource,
		Source:            MapDataSourceToModel(entity.Source),
		SourceOfTruth:     MapDataSourceToModel(entity.SourceOfTruth),
	}
}

func MapInteractionSessionInputToEntity(model *model.InteractionSessionInput) *entity.InteractionSessionEntity {
	if model == nil {
		return nil
	}
	return &entity.InteractionSessionEntity{
		SessionIdentifier: model.SessionIdentifier,
		CreatedAt:         utils.Now(),
		Name:              model.Name,
		Status:            model.Status,
		Type:              model.Type,
		Channel:           model.Channel,
		ChannelData:       model.ChannelData,
		AppSource:         model.AppSource,
		Source:            neo4jentity.DataSourceOpenline,
		SourceOfTruth:     neo4jentity.DataSourceOpenline,
	}
}

func MapEntitiesToInteractionSessions(entities *entity.InteractionSessionEntities) []*model.InteractionSession {
	var interactionSessions []*model.InteractionSession
	for _, interactionSessionEntity := range *entities {
		interactionSessions = append(interactionSessions, MapEntityToInteractionSession(&interactionSessionEntity))
	}
	return interactionSessions
}
