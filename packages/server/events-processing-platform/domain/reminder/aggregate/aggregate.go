package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

const ReminderAggregateType = "reminder"

type ReminderAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Reminder *model.Reminder
}

func NewReminderAggregateWithTenantAndID(tenant, id string) *ReminderAggregate {
	reminderAggregate := ReminderAggregate{}
	reminderAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(ReminderAggregateType, tenant, id)
	reminderAggregate.SetWhen(reminderAggregate.When)
	reminderAggregate.Reminder = &model.Reminder{}
	reminderAggregate.Tenant = tenant

	return &reminderAggregate
}

func (a *ReminderAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case events.ReminderCreateV1:
		return a.whenReminderCreate(event)
	case events.ReminderUpdateV1:
		return a.whenReminderUpdate(event)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *ReminderAggregate) whenReminderCreate(eventData interface{}) error {
	reminderCreateEvent := eventData.(*events.ReminderCreateEvent)
	a.Reminder = &model.Reminder{
		Content:        reminderCreateEvent.Content,
		DueDate:        reminderCreateEvent.DueDate,
		Dismissed:      reminderCreateEvent.Dismissed,
		CreatedAt:      reminderCreateEvent.CreatedAt,
		UserID:         reminderCreateEvent.UserId,
		OrganizationID: reminderCreateEvent.OrganizationId,
		ID:             reminderCreateEvent.Id,
	}
	return nil
}

func (a *ReminderAggregate) whenReminderUpdate(eventData interface{}) error {
	reminderUpdateEvent := eventData.(*events.ReminderUpdateEvent)
	if reminderUpdateEvent.UpdateContent() {
		a.Reminder.Content = reminderUpdateEvent.Content
	}
	if reminderUpdateEvent.UpdateDueDate() {
		a.Reminder.DueDate = reminderUpdateEvent.DueDate
	}
	if reminderUpdateEvent.UpdateDismissed() {
		a.Reminder.Dismissed = reminderUpdateEvent.Dismissed
	}
	return nil
}
