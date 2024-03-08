package events

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type ReminderUpdateEvent struct {
	Tenant     string    `json:"tenant" validate:"required"`
	ReminderId string    `json:"reminderId" validate:"required"`
	Content    string    `json:"content,omitempty"`
	DueDate    time.Time `json:"dueDate,omitempty"`
	Dismissed  bool      `json:"dismissed,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt"`
	FieldsMask []string  `json:"fieldsMask,omitempty"`
}

func NewReminderUpdateEvent(aggregate eventstore.Aggregate, reminderId, tenant string, content string, dueDate time.Time, dismissed bool, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := ReminderUpdateEvent{
		Tenant:     tenant,
		ReminderId: reminderId,
		UpdatedAt:  updatedAt,
		FieldsMask: fieldsMask,
	}
	if eventData.UpdateContent() {
		eventData.Content = content
	}
	if eventData.UpdateDueDate() {
		eventData.DueDate = dueDate
	}
	if eventData.UpdateDismissed() {
		eventData.Dismissed = dismissed
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ReminderUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ReminderUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ReminderUpdateEvent")
	}

	return event, nil
}

func (e ReminderUpdateEvent) UpdateContent() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskContent)
}

func (e ReminderUpdateEvent) UpdateDueDate() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskDueDate)
}

func (e ReminderUpdateEvent) UpdateDismissed() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskDismissed)
}
