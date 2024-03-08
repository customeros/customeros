package event_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type EventHandlers struct {
	CreateReminderHandler CreateReminderHandler
	UpdateReminderHandler UpdateReminderHandler
}

func NewEventHandlers(log logger.Logger, es eventstore.AggregateStore) *EventHandlers {
	return &EventHandlers{
		CreateReminderHandler: NewCreateReminderHandler(log, es),
		UpdateReminderHandler: NewUpdateReminderHandler(log, es),
	}
}
