package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type ContactShowEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewContactShowEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ContactShowEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactShowEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactShowV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactShowEvent")
	}
	return event, nil
}
