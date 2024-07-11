package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type OrganizationAdjustIndustryEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewOrganizationAdjustIndustryEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := OrganizationAdjustIndustryEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationAdjustIndustryEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationAdjustIndustryV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationAdjustIndustryEvent")
	}
	return event, nil
}
