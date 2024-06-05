package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

type OrganizationRefreshDerivedData struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewOrganizationRefreshDerivedData(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := OrganizationRefreshDerivedData{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationRefreshDerivedData")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationRefreshDerivedDataV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationRefreshDerivedData")
	}
	return event, nil
}
