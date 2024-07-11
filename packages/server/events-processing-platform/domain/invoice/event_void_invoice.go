package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type InvoiceVoidEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewInvoiceVoidEvent(aggregate eventstore.Aggregate, updatedAt time.Time) (eventstore.Event, error) {
	eventData := InvoiceVoidEvent{
		Tenant:    aggregate.GetTenant(),
		UpdatedAt: updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceVoidEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceVoidV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceVoidEvent")
	}

	return event, nil
}
