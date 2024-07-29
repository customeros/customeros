package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type ReminderUpdateEvent struct {
	event.BaseEvent
	Content    string    `json:"content,omitempty"`
	DueDate    time.Time `json:"dueDate,omitempty"`
	Dismissed  bool      `json:"dismissed,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt"`
	FieldsMask []string  `json:"fieldsMask,omitempty"`
}

func (e ReminderUpdateEvent) GetBaseEvent() event.BaseEvent {
	return e.BaseEvent
}

func (e *ReminderUpdateEvent) SetEntityId(entityId string) {
}

func (e ReminderUpdateEvent) UpdateContent() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskReminderContent)
}

func (e ReminderUpdateEvent) UpdateDueDate() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskReminderDueDate)
}

func (e ReminderUpdateEvent) UpdateDismissed() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskReminderDismissed)
}
