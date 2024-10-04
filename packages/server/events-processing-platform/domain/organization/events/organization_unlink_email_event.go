package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationUnlinkEmailEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	Email  string `json:"email"`
}

func NewOrganizationUnlinkEmailEvent(aggregate eventstore.Aggregate, email string) (eventstore.Event, error) {
	eventData := OrganizationUnlinkEmailEvent{
		Tenant: aggregate.GetTenant(),
		Email:  email,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationUnlinkEmailEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationEmailUnlinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationUnlinkEmailEvent")
	}
	return event, nil
}
