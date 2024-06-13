package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OpportunityCloseLooseEvent struct {
	Tenant   string    `json:"tenant" validate:"required"`
	ClosedAt time.Time `json:"closedAt" validate:"required"`
}

func NewOpportunityCloseLooseEvent(aggregate eventstore.Aggregate, closedAt time.Time) (eventstore.Event, error) {
	eventData := OpportunityCloseLooseEvent{
		Tenant:   aggregate.GetTenant(),
		ClosedAt: closedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OpportunityCloseLooseEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityCloseLooseV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OpportunityCloseLooseEvent")
	}
	return event, nil
}
