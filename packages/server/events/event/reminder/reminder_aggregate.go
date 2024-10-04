package reminder

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"strings"
)

const ReminderAggregateType = "reminder"

type ReminderAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Reminder *Reminder
}

func NewReminderAggregateWithTenantAndID(tenant, id string) *ReminderAggregate {
	reminderAggregate := ReminderAggregate{}
	reminderAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(ReminderAggregateType, tenant, id)
	reminderAggregate.SetWhen(reminderAggregate.When)
	reminderAggregate.Reminder = &Reminder{}
	reminderAggregate.Tenant = tenant

	return &reminderAggregate
}

func (a *ReminderAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ReminderCreateV1:
		return a.whenReminderCreate(evt)
	case event.ReminderUpdateV1:
		return a.whenReminderUpdate(evt)
	case event.ReminderNotificationV1:
		return nil
	default:
		if strings.HasPrefix(evt.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *ReminderAggregate) whenReminderCreate(evt eventstore.Event) error {
	var reminderCreateEvent event.ReminderCreateEvent

	if err := evt.GetJsonData(&reminderCreateEvent); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Reminder = &Reminder{
		Content:        reminderCreateEvent.Content,
		DueDate:        reminderCreateEvent.DueDate,
		Dismissed:      reminderCreateEvent.Dismissed,
		CreatedAt:      reminderCreateEvent.CreatedAt,
		UserID:         reminderCreateEvent.UserId,
		OrganizationID: reminderCreateEvent.OrganizationId,
	}
	return nil
}

func (a *ReminderAggregate) whenReminderUpdate(evt eventstore.Event) error {
	var reminderUpdateEvent event.ReminderUpdateEvent

	if err := evt.GetJsonData(&reminderUpdateEvent); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

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
