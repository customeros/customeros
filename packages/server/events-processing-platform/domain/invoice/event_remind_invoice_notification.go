package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type InvoiceRemindNotificationEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewInvoiceRemindNotificationEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := InvoiceRemindNotificationEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceRemindNotificationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceRemindNotificationV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceRemindNotificationEvent")
	}

	return event, nil
}
