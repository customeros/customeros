package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
)

type ReminderNotificationEvent struct {
	event.BaseEvent
	OrganizationId string `json:"organizationId" validate:"required"`
	UserId         string `json:"userId"` // who set the reminder
	Content        string `json:"content"`
}

func (e ReminderNotificationEvent) GetBaseEvent() event.BaseEvent {
	return e.BaseEvent
}

func (e ReminderNotificationEvent) SetEntityId(entityId string) {
	e.EntityId = entityId
}
