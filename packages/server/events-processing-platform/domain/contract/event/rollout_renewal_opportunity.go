package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type RolloutRenewalOpportunityEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewRolloutRenewalOpportunityEvent(aggr eventstore.Aggregate) (eventstore.Event, error) {
	eventData := RolloutRenewalOpportunityEvent{
		Tenant: aggr.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate RolloutRenewalOpportunityEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractRolloutRenewalOpportunityV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for RolloutRenewalOpportunityEvent")
	}
	return event, nil
}
