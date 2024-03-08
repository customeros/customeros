package events

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type ReminderCreateEvent struct {
	Tenant         string    `json:"tenant" validate:"required"`
	Content        string    `json:"content"`
	DueDate        time.Time `json:"dueDate"`
	UserId         string    `json:"userId" validate:"required"`
	OrganizationId string    `json:"organizationId" validate:"required"`
	Dismissed      bool      `json:"dismissed"`
	CreatedAt      time.Time `json:"createdAt"`
	Id             string    `json:"id"`
}

func NewReminderCreateEvent(aggregate eventstore.Aggregate, tenant, content, id, userId, organizationId string, dismissed bool, createdAt, dueDate time.Time) (eventstore.Event, error) {
	eventData := ReminderCreateEvent{
		Id:             id,
		Tenant:         tenant,
		Content:        content,
		DueDate:        dueDate,
		UserId:         userId,
		OrganizationId: organizationId,
		Dismissed:      dismissed,
		CreatedAt:      createdAt,
	}
	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ReminderCreateEvent")
	}
	event := eventstore.NewBaseEvent(aggregate, ReminderCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ReminderCreateEvent")
	}
	return event, nil
}
