package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

func MapMeetingSchedulingInputToEntity(input model.SaveMeetingSchedulingInput, entity *postgres_entity.MeetingScheduling) *postgres_entity.MeetingScheduling {
	if entity == nil {
		entity = &postgres_entity.MeetingScheduling{}
	}

	if input.Title != nil {
		entity.Title = *input.Title
	}

	return entity
}

func MapMeetingSchedulingEntityToModel(entity *postgres_entity.MeetingScheduling) *model.MeetingScheduling {
	if entity == nil {
		return nil
	}
	return &model.MeetingScheduling{
		ID:        entity.ID,
		Title:     entity.Title,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func MapMeetingSchedulingEntitiesToModels(entities []*postgres_entity.MeetingScheduling) []*model.MeetingScheduling {
	if entities == nil {
		return nil
	}

	models := make([]*model.MeetingScheduling, len(entities))
	for i, entity := range entities {
		models[i] = MapMeetingSchedulingEntityToModel(entity)
	}
	return models
}
