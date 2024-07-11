package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type InvoicePayNotificationEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewInvoicePayNotificationEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := InvoicePayNotificationEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoicePayNotificationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoicePayNotificationV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoicePayNotificationEvent")
	}

	return event, nil
}
