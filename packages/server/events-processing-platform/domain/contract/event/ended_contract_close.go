package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type EndedContractCloseEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	Status string `json:"status" validate:"required"`
}

func NewEndedContractCloseEvent(aggr eventstore.Aggregate) (eventstore.Event, error) {
	eventData := EndedContractCloseEvent{
		Tenant: aggr.GetTenant(),
		Status: string(model.ContractStatusStringEnded),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate EndedContractCloseEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractCloseEndedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for EndedContractCloseEvent")
	}
	return event, nil
}
