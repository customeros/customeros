package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContactAddTagEvent struct {
	Tenant   string    `json:"tenant" validate:"required"`
	TagId    string    `json:"tagId" validate:"required"`
	TaggedAt time.Time `json:"taggedAt" validate:"required"`
}

func NewContactAddTagEvent(aggregate eventstore.Aggregate, tagId string, taggedAt time.Time) (eventstore.Event, error) {
	eventData := ContactAddTagEvent{
		Tenant:   aggregate.GetTenant(),
		TagId:    tagId,
		TaggedAt: taggedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactAddTagEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactAddTagV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactAddTagEvent")
	}
	return event, nil
}
