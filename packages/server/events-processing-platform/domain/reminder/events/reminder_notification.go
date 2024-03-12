package events

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type ReminderNotificationEvent struct {
	Tenant         string    `json:"tenant" validate:"required"`
	CreatedAt      time.Time `json:"createdAt"`
	OrganizationId string    `json:"organizationId" validate:"required"`
	UserId         string    `json:"userId"` // who set the reminder
	Content        string    `json:"content"`
}

func NewReminderNotificationEvent(aggregate eventstore.Aggregate, userId, organizationId, content string, createdAt time.Time) (eventstore.Event, error) {
	eventData := ReminderNotificationEvent{
		Tenant:         aggregate.GetTenant(),
		CreatedAt:      createdAt,
		OrganizationId: organizationId,
		UserId:         userId,
		Content:        content,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ReminderNotificationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ReminderNotificationV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ReminderNotificationEvent")
	}
	return event, nil
}
