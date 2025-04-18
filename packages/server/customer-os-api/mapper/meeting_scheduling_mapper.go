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
	if input.DurationMins != nil {
		entity.DurationMins = *input.DurationMins
	}
	if input.Description != nil {
		entity.Description = *input.Description
	}
	if input.AllowedParticipants != nil {
		entity.AllowedParticipants = input.AllowedParticipants
	}

	if input.BookingFormName != nil {
		entity.BookingFormName = *input.BookingFormName
	}
	if input.BookingFormEmail != nil {
		entity.BookingFormEmail = *input.BookingFormEmail
	}
	if input.BookingFormPhone != nil {
		entity.BookingFormPhone = *input.BookingFormPhone
	}

	if input.BookOptionEnabled != nil {
		entity.BookOptionEnabled = *input.BookOptionEnabled
	}
	if input.BookOptionBufferBetweenMeetingsMins != nil {
		entity.BookOptionBufferBetweenMeetingsMins = *input.BookOptionBufferBetweenMeetingsMins
	}
	if input.BookOptionDaysInAdvance != nil {
		entity.BookOptionDaysInAdvance = *input.BookOptionDaysInAdvance
	}
	if input.BookOptionMinNoticeMins != nil {
		entity.BookOptionMinNoticeMins = *input.BookOptionMinNoticeMins
	}
	if input.BookOptionRedirectLink != nil {
		entity.BookOptionRedirectLink = *input.BookOptionRedirectLink
	}

	if input.EmailNotificationEnabled != nil {
		entity.EmailNotificationEnabled = *input.EmailNotificationEnabled
	}

	return entity
}

func MapMeetingSchedulingEntityToModel(entity *postgres_entity.MeetingScheduling) *model.MeetingScheduling {
	if entity == nil {
		return nil
	}
	return &model.MeetingScheduling{
		ID:                                  entity.ID,
		Title:                               entity.Title,
		DurationMins:                        entity.DurationMins,
		Description:                         entity.Description,
		AllowedParticipants:                 entity.AllowedParticipants,
		BookingFormName:                     entity.BookingFormName,
		BookingFormEmail:                    entity.BookingFormEmail,
		BookingFormPhone:                    entity.BookingFormPhone,
		BookOptionEnabled:                   entity.BookOptionEnabled,
		BookOptionBufferBetweenMeetingsMins: entity.BookOptionBufferBetweenMeetingsMins,
		BookOptionDaysInAdvance:             entity.BookOptionDaysInAdvance,
		BookOptionMinNoticeMins:             entity.BookOptionMinNoticeMins,
		BookOptionRedirectLink:              entity.BookOptionRedirectLink,
		EmailNotificationEnabled:            entity.EmailNotificationEnabled,
		CreatedAt:                           entity.CreatedAt,
		UpdatedAt:                           entity.UpdatedAt,
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
