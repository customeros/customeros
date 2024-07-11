package reminder

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type ReminderCreateEvent struct {
	Tenant         string        `json:"tenant" validate:"required"`
	Content        string        `json:"content"`
	DueDate        time.Time     `json:"dueDate"`
	UserId         string        `json:"userId" validate:"required"`
	OrganizationId string        `json:"organizationId" validate:"required"`
	Dismissed      bool          `json:"dismissed"`
	CreatedAt      time.Time     `json:"createdAt"`
	SourceFields   events.Source `json:"sourceFields" validate:"required"`
}

func NewReminderCreateEvent(aggregate eventstore.Aggregate, content, userId, organizationId string, dismissed bool, createdAt, dueDate time.Time, sourceFields events.Source) (eventstore.Event, error) {
	eventData := ReminderCreateEvent{
		Tenant:         aggregate.GetTenant(),
		Content:        content,
		DueDate:        dueDate,
		UserId:         userId,
		OrganizationId: organizationId,
		Dismissed:      dismissed,
		CreatedAt:      createdAt,
		SourceFields:   sourceFields,
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

func (e *ReminderCreateEvent) String() string {
	dismissedStr := "false"
	if e.Dismissed {
		dismissedStr = "true"
	}
	return "ReminderCreateEvent{" +
		"Tenant: " + e.Tenant +
		", Content: " + e.Content +
		", DueDate: " + e.DueDate.String() +
		", UserId: " + e.UserId +
		", OrganizationId: " + e.OrganizationId +
		", Dismissed: " + dismissedStr +
		", CreatedAt: " + e.CreatedAt.String() +
		", SourceFields: " + e.SourceFields.String() +
		"}"
}
