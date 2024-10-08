package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type ServiceLineItemResumeEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewServiceLineItemResumeEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ServiceLineItemResumeEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ServiceLineItemResumeEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ServiceLineItemResumeV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ServiceLineItemResumeEvent")
	}
	return event, nil
}
