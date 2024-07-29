package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
	"time"
)

type ReminderCreateEvent struct {
	event.BaseEvent
	Content        string    `json:"content"`
	DueDate        time.Time `json:"dueDate"`
	UserId         string    `json:"userId" validate:"required"`
	OrganizationId string    `json:"organizationId" validate:"required"`
	Dismissed      bool      `json:"dismissed"`
}

func (e ReminderCreateEvent) GetBaseEvent() event.BaseEvent {
	return e.BaseEvent
}

func (e *ReminderCreateEvent) SetEntityId(entityId string) {
	e.BaseEvent.EntityId = entityId
}
