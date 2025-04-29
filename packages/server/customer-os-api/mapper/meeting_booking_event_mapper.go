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
		if entity.DurationMins%5 != 0 {
			entity.DurationMins = ((entity.DurationMins / 5) + 1) * 5
		}
	}
	if input.Description != nil {
		entity.Description = *input.Description
	}
	if input.ParticipantEmails != nil {
		entity.ParticipantEmails = input.ParticipantEmails
	}

	if input.BookingFormNameEnabled != nil {
		entity.BookingFormNameEnabled = *input.BookingFormNameEnabled
	}
	if input.BookingFormEmailEnabled != nil {
		entity.BookingFormEmailEnabled = *input.BookingFormEmailEnabled
	}
	if input.BookingFormPhoneEnabled != nil {
		entity.BookingFormPhoneEnabled = *input.BookingFormPhoneEnabled
	}
	if input.BookingFormPhoneRequired != nil {
		entity.BookingFormPhoneRequired = *input.BookingFormPhoneRequired
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
		ParticipantEmails:                   entity.ParticipantEmails,
		BookingFormNameEnabled:              entity.BookingFormNameEnabled,
		BookingFormEmailEnabled:             entity.BookingFormEmailEnabled,
		BookingFormPhoneEnabled:             entity.BookingFormPhoneEnabled,
		BookingFormPhoneRequired:            entity.BookingFormPhoneRequired,
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
