package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OrganizationAddTagEvent struct {
	Tenant   string    `json:"tenant" validate:"required"`
	TagId    string    `json:"tagId" validate:"required"`
	TaggedAt time.Time `json:"taggedAt" validate:"required"`
}

func NewOrganizationAddTagEvent(aggregate eventstore.Aggregate, tagId string, taggedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationAddTagEvent{
		Tenant:   aggregate.GetTenant(),
		TagId:    tagId,
		TaggedAt: taggedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationAddTagEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationAddTagV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationAddTagEvent")
	}
	return event, nil
}
