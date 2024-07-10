package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type ContractRefreshLtvEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewContractRefreshLtvEvent(aggr eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ContractRefreshLtvEvent{
		Tenant: aggr.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractRefreshLtvEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractRefreshLtvV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractRefreshLtvEvent")
	}
	return event, nil
}
