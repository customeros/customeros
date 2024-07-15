package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"
)

type ReminderCreateEvent struct {
	events.BaseEvent
	Content        string    `json:"content"`
	DueDate        time.Time `json:"dueDate"`
	UserId         string    `json:"userId" validate:"required"`
	OrganizationId string    `json:"organizationId" validate:"required"`
	Dismissed      bool      `json:"dismissed"`
}

func (e ReminderCreateEvent) GetBaseEvent() events.BaseEvent {
	return e.BaseEvent
}

func (e *ReminderCreateEvent) SetEntityId(entityId string) {
	e.BaseEvent.EntityId = entityId
}
