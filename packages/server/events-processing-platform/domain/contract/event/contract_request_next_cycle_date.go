package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type ContractRequestNextCycleDateEvent struct {
	Tenant      string    `json:"tenant" validate:"required"`
	RequestedAt time.Time `json:"requestedAt"`
}

func NewContractRequestNextCycleDateEvent(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ContractRequestNextCycleDateEvent{
		Tenant:      aggregate.GetTenant(),
		RequestedAt: utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractRequestNextCycleDateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContractRequestNextCycleDateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractRequestNextCycleDateEvent")
	}
	return event, nil
}
