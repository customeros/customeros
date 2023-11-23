package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
)

type ResetRenewalOpportunityEvent struct {
	Tenant string `json:"tenant" validate:"required"`
}

func NewResetRenewalOpportunityEvent(aggr eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ResetRenewalOpportunityEvent{
		Tenant: aggr.GetTenant(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ResetRenewalOpportunityEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractResetRenewalOpportunityV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ResetRenewalOpportunityEvent")
	}
	return event, nil
}
