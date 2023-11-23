package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type ContractUpdateStatusEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	Status string `json:"status" validate:"required"`
}

func NewContractUpdateStatusEvent(aggr eventstore.Aggregate, status string) (eventstore.Event, error) {
	eventData := ContractUpdateStatusEvent{
		Tenant: aggr.GetTenant(),
		Status: status,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractUpdateStatusEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractUpdateStatusV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractUpdateStatusEvent")
	}
	return event, nil
}
