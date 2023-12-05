package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type ServiceLineItemCloseEvent struct {
	Tenant     string    `json:"tenant" validate:"required"`
	EndedAt    time.Time `json:"endedAt" validate:"required"`
	UpdatedAt  time.Time `json:"updatedAt"`
	IsCanceled bool      `json:"isCanceled"`
}

func NewServiceLineItemCloseEvent(aggregate eventstore.Aggregate, endedAt, updatedAt time.Time, isCanceled bool) (eventstore.Event, error) {
	eventData := ServiceLineItemCloseEvent{
		Tenant:     aggregate.GetTenant(),
		EndedAt:    endedAt,
		UpdatedAt:  updatedAt,
		IsCanceled: isCanceled,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ServiceLineItemCloseEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ServiceLineItemCloseV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ServiceLineItemCloseEvent")
	}
	return event, nil
}
