package events

import (
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
)

const (
	ContactCreated = "CONTACT_CREATED"
	ContactUpdated = "CONTACT_UPDATED"
	ContactDeleted = "CONTACT_DELETED"
)

type ContactCreatedEvent struct {
	Uuid      string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
}

func NewContactCreatedEvent(aggregate eventstore.Aggregate, uuid string, firstName string, lastName string) (eventstore.Event, error) {
	eventData := ContactCreatedEvent{
		Uuid:      uuid,
		FirstName: firstName,
		LastName:  lastName,
	}
	event := eventstore.NewBaseEvent(aggregate, ContactCreated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactUpdatedEvent struct {
	Uuid      string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
}

func NewContactUpdatedEvent(aggregate eventstore.Aggregate, uuid string, firstName string, lastName string) (eventstore.Event, error) {
	eventData := ContactUpdatedEvent{
		Uuid:      uuid,
		FirstName: firstName,
		LastName:  lastName,
	}
	event := eventstore.NewBaseEvent(aggregate, ContactUpdated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactDeletedEvent struct {
	Uuid string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
}

func NewContactDeletedEvent(aggregate eventstore.Aggregate, uuid string) (eventstore.Event, error) {
	eventData := ContactDeletedEvent{
		Uuid: uuid,
	}
	event := eventstore.NewBaseEvent(aggregate, ContactDeleted)
	err := event.SetJsonData(&eventData)
	if err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
