package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

func MapMeetingBookingEventInputToEntity(input model.SaveMeetingBookingEventInput, entity *postgresEntity.MeetingBookingEvent) *postgresEntity.MeetingBookingEvent {
	if entity == nil {
		entity = &postgresEntity.MeetingBookingEvent{}
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

	if input.EmailNotificationEnabled != nil {
		entity.EmailNotificationEnabled = *input.EmailNotificationEnabled
	}
	if input.BookingConfirmationRedirectLink != nil {
		entity.BookingConfirmationRedirectLink = *input.BookingConfirmationRedirectLink
	}
	if input.AssignmentMethod != nil {
		entity.AssignmentMethod = enummapper.MapMeetingBookingAssignmentMethodFromModel(*input.AssignmentMethod)
	}
	if input.ShowLogo != nil {
		entity.ShowLogo = *input.ShowLogo
	}
	if input.Location != nil {
		entity.Location = *input.Location
	}

	return entity
}

func MapMeetingBookingEventEntityToModel(entity *postgresEntity.MeetingBookingEvent) *model.MeetingBookingEvent {
	if entity == nil {
		return nil
	}
	return &model.MeetingBookingEvent{
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
		EmailNotificationEnabled:            entity.EmailNotificationEnabled,
		BookingConfirmationRedirectLink:     entity.BookingConfirmationRedirectLink,
		AssignmentMethod:                    enummapper.MapMeetingBookingAssignmentMethodToModel(entity.AssignmentMethod),
		ShowLogo:                            entity.ShowLogo,
		Location:                            entity.Location,
		CreatedAt:                           entity.CreatedAt,
		UpdatedAt:                           entity.UpdatedAt,
	}
}

func MapMeetingBookingEventEntitiesToModels(entities []*postgresEntity.MeetingBookingEvent) []*model.MeetingBookingEvent {
	if entities == nil {
		return nil
	}

	models := make([]*model.MeetingBookingEvent, len(entities))
	for i, entity := range entities {
		models[i] = MapMeetingBookingEventEntityToModel(entity)
	}
	return models
}
