package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type InvoicePdfRequestedEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewInvoicePdfRequestedEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := InvoicePdfRequestedEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoicePdfRequestedEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoicePdfRequestedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoicePdfRequestedEvent")
	}

	return event, nil
}
