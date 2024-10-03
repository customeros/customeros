package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type ServiceLineItemPauseEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewServiceLineItemPauseEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ServiceLineItemPauseEvent{
		Tenant: aggregate.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ServiceLineItemPauseEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ServiceLineItemPauseV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ServiceLineItemPauseEvent")
	}
	return event, nil
}
