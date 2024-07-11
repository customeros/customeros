package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationUnlinkDomainEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	Domain string `json:"domain" validate:"required"`
}

func NewOrganizationUnlinkDomainEvent(aggregate eventstore.Aggregate, domain string) (eventstore.Event, error) {
	eventData := OrganizationUnlinkDomainEvent{
		Tenant: aggregate.GetTenant(),
		Domain: domain,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationUnlinkDomainEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationUnlinkDomainV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationUnlinkDomainEvent")
	}
	return event, nil
}
