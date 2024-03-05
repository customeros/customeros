package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type ServiceLineItemDeleteEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewServiceLineItemDeleteEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ServiceLineItemDeleteEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ServiceLineItemDeleteEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ServiceLineItemDeleteV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ServiceLineItemDeleteEvent")
	}
	return event, nil
}
