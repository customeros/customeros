package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContactRemoveTagEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	TagId  string `json:"tagId" validate:"required"`
}

func NewContactRemoveTagEvent(aggregate eventstore.Aggregate, tagId string, taggedAt time.Time) (eventstore.Event, error) {
	eventData := ContactRemoveTagEvent{
		Tenant: aggregate.GetTenant(),
		TagId:  tagId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactRemoveTagEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactRemoveTagV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactRemoveTagEvent")
	}
	return event, nil
}
