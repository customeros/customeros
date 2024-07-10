package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type EmailValidateEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewEmailValidateEvent(aggr eventstore.Aggregate) (eventstore.Event, error) {
	eventData := EmailValidateEvent{
		Tenant: aggr.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate EmailValidateEvent")
	}

	event := eventstore.NewBaseEvent(aggr, EmailValidateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for EmailValidateEvent")
	}
	return event, nil
}
