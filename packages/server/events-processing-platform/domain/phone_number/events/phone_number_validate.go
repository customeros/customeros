package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

type PhoneNumberValidateEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewPhoneNumberValidateEvent(aggr eventstore.Aggregate) (eventstore.Event, error) {
	eventData := PhoneNumberValidateEvent{
		Tenant: aggr.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate PhoneNumberValidateEvent")
	}

	event := eventstore.NewBaseEvent(aggr, PhoneNumberValidateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for PhoneNumberValidateEvent")
	}
	return event, nil
}
