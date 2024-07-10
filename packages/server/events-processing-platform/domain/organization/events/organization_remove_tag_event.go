package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OrganizationRemoveTagEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	TagId  string `json:"tagId" validate:"required"`
}

func NewOrganizationRemoveTagEvent(aggregate eventstore.Aggregate, tagId string, taggedAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationRemoveTagEvent{
		Tenant: aggregate.GetTenant(),
		TagId:  tagId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationRemoveTagEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationRemoveTagV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationRemoveTagEvent")
	}
	return event, nil
}
