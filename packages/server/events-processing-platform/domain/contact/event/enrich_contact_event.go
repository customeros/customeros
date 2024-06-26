package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContactRequestEnrich struct {
	Tenant      string    `json:"tenant" validate:"required"`
	RequestedAt time.Time `json:"requestedAt"`
}

func NewContactRequestEnrich(aggregate eventstore.Aggregate) (eventstore.Event, error) {
	eventData := ContactRequestEnrich{
		Tenant:      aggregate.GetTenant(),
		RequestedAt: utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactRequestEnrich")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactRequestEnrichV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactRequestEnrich")
	}
	return event, nil
}
