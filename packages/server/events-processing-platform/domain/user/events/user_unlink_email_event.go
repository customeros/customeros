package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type UserUnlinkEmailEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	Email  string `json:"email"`
}

func NewUserUnlinkEmailEvent(aggregate eventstore.Aggregate, email string) (eventstore.Event, error) {
	eventData := UserUnlinkEmailEvent{
		Tenant: aggregate.GetTenant(),
		Email:  email,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UserUnlinkEmailEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, UserEmailUnlinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UserUnlinkEmailEvent")
	}
	return event, nil
}
