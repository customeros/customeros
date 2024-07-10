package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type InvoicePaidEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewInvoicePaidEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := InvoicePaidEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoicePaidEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoicePaidV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoicePaidEvent")
	}

	return event, nil
}
