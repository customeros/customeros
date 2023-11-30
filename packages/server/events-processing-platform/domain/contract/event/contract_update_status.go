package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type ContractUpdateStatusEvent struct {
	Tenant           string     `json:"tenant" validate:"required"`
	Status           string     `json:"status" validate:"required"`
	ServiceStartedAt *time.Time `json:"serviceStartedAt,omitempty"`
	EndedAt          *time.Time `json:"endedAt,omitempty"`
}

func NewContractUpdateStatusEvent(aggr eventstore.Aggregate, status string, serviceStartedAt, endedAt *time.Time) (eventstore.Event, error) {
	eventData := ContractUpdateStatusEvent{
		Tenant:           aggr.GetTenant(),
		Status:           status,
		ServiceStartedAt: serviceStartedAt,
		EndedAt:          endedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractUpdateStatusEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractUpdateStatusV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractUpdateStatusEvent")
	}
	return event, nil
}
