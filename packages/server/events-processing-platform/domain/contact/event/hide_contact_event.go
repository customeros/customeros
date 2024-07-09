package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

type ContactHideEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewContactHideEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ContactHideEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactHideEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactHideV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactHideEvent")
	}
	return event, nil
}
