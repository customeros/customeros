package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type OpportunityCloseWinEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
	ClosedAt  time.Time `json:"closedAt" validate:"required"`
}

func NewOpportunityCloseWinEvent(aggregate eventstore.Aggregate, updatedAt, closedAt time.Time) (eventstore.Event, error) {
	eventData := OpportunityCloseWinEvent{
		Tenant:    aggregate.GetTenant(),
		UpdatedAt: updatedAt,
		ClosedAt:  closedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OpportunityCloseWinEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OpportunityCloseWinV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OpportunityCloseWinEvent")
	}
	return event, nil
}
