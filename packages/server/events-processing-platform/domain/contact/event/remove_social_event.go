package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

type ContactRemoveSocialEvent struct {
	Tenant   string `json:"tenant" validate:"required"`
	SocialId string `json:"socialId"`
	Url      string `json:"url"`
}

func NewContactRemoveSocialEvent(aggregate eventstore.Aggregate, socialId, url string) (eventstore.Event, error) {
	eventData := ContactRemoveSocialEvent{
		Tenant:   aggregate.GetTenant(),
		SocialId: socialId,
		Url:      url,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactRemoveSocialEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactRemoveSocialV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactRemoveSocialEvent")
	}
	return event, nil
}
