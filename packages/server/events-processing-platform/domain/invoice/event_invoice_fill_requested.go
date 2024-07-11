package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type InvoiceFillRequestedEvent struct {
	Tenant     string `json:"tenant" validate:"required"`
	ContractId string `json:"contractId" validate:"required"`
}

func NewInvoiceFillRequestedEvent(aggregate eventstore.Aggregate, contractId string) (eventstore.Event, error) {
	eventData := InvoiceFillRequestedEvent{
		Tenant:     aggregate.GetTenant(),
		ContractId: contractId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceFillRequestedEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceFillRequestedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceFillRequestedEvent")
	}

	return event, nil
}
