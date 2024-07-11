package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type InvoiceDeleteEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewInvoiceDeleteEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := InvoiceDeleteEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceDeleteEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceDeleteV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceDeleteEvent")
	}

	return event, nil
}
