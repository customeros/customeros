package generic

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
)

type LinkEntityWithEntity struct {
	events.BaseEvent
	Tenant string `json:"tenant" validate:"required"`

	WithEntityId   string            `json:"withEntityId"`
	WithEntityType events.EntityType `json:"withEntityType"`

	RelationshipName string `json:"relationshipName"`

	//todo enhance with relationship properties
}

func NewLinkEntityWithEntity(aggregate eventstore.Aggregate, eventData LinkEntityWithEntity) (eventstore.Event, error) {
	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactAddLocationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, LinkEntityWithEntityV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactAddLocationEvent")
	}
	return event, nil
}
