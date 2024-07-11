package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContractDeleteEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewContractDeleteEvent(aggregate eventstore.Aggregate, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContractDeleteEvent{
		Tenant:    aggregate.GetTenant(),
		UpdatedAt: updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractDeleteEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContractDeleteV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractDeleteEvent")
	}
	return event, nil
}
