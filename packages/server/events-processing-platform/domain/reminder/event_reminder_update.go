package reminder

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type ReminderUpdateEvent struct {
	Tenant     string    `json:"tenant" validate:"required"`
	Content    string    `json:"content,omitempty"`
	DueDate    time.Time `json:"dueDate,omitempty"`
	Dismissed  bool      `json:"dismissed,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt"`
	FieldsMask []string  `json:"fieldsMask,omitempty"`
}

func NewReminderUpdateEvent(aggregate eventstore.Aggregate, content string, dueDate time.Time, dismissed bool, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := ReminderUpdateEvent{
		Tenant:     aggregate.GetTenant(),
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

func (e *ReminderUpdateEvent) String() string {
	dismissedStr := "false"
	if e.Dismissed {
		dismissedStr = "true"
	}
	return "ReminderCreateEvent{" +
		"Tenant: " + e.Tenant +
		", Content: " + e.Content +
		", DueDate: " + e.DueDate.String() +
		", Dismissed: " + dismissedStr +
		", UpdatedAt: " + e.UpdatedAt.String() +
		"}"
}
