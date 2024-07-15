package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
)

type ReminderNotificationEvent struct {
	events.BaseEvent
	OrganizationId string `json:"organizationId" validate:"required"`
	UserId         string `json:"userId"` // who set the reminder
	Content        string `json:"content"`
}

func (e ReminderNotificationEvent) GetBaseEvent() events.BaseEvent {
	return e.BaseEvent
}

func (e ReminderNotificationEvent) SetEntityId(entityId string) {
}
