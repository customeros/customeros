package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type ContractRefreshStatusEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewContractRefreshStatusEvent(aggr eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ContractRefreshStatusEvent{
		Tenant: aggr.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractRefreshStatusEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractRefreshStatusV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractRefreshStatusEvent")
	}
	return event, nil
}
