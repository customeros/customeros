package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type EmailDeleteEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewEmailDeleteEvent(aggregate eventstore.Aggregate, tenant string) (eventstore.Event, error) {
	eventData := EmailDeleteEvent{
		Tenant: tenant,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate EmailDeleteEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, EmailDeleteV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for EmailDeleteEvent")
	}
	return event, nil
}
