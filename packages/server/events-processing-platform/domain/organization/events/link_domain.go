package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationLinkDomainEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	Domain string `json:"domain" validate:"required"`
}

func NewOrganizationLinkDomainEvent(aggregate eventstore.Aggregate, domain string) (eventstore.Event, error) {
	eventData := OrganizationLinkDomainEvent{
		Tenant: aggregate.GetTenant(),
		Domain: domain,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationLinkDomainEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationLinkDomainV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationLinkDomainEvent")
	}
	return event, nil
}
