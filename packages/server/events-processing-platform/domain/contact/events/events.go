package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

const (
	ContactCreated = "CONTACT_CREATED"
	ContactUpdated = "CONTACT_UPDATED"
	//ContactDeleted = "CONTACT_DELETED"
)

type ContactCreatedEvent struct {
	Tenant        string    `json:"tenant"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Name          string    `json:"name"`
	Prefix        string    `json:"prefix"`
	Source        string    `json:"source"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	AppSource     string    `json:"appSource"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewContactCreatedEvent(aggregate eventstore.Aggregate, tenant, firstName, lastName, name, prefix, source, sourceOfTruth, appSource string, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactCreatedEvent{
		Tenant:        tenant,
		FirstName:     firstName,
		LastName:      lastName,
		Name:          name,
		Prefix:        prefix,
		Source:        source,
		SourceOfTruth: sourceOfTruth,
		AppSource:     appSource,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
	event := eventstore.NewBaseEvent(aggregate, ContactCreated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type ContactUpdatedEvent struct {
	Tenant        string    `json:"tenant"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Name          string    `json:"name"`
	Prefix        string    `json:"prefix"`
}

func NewContactUpdatedEvent(aggregate eventstore.Aggregate, tenant, firstName, lastName, name, prefix, sourceOfTruth string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactUpdatedEvent{
		FirstName:     firstName,
		LastName:      lastName,
		Name:          name,
		Prefix:        prefix,
		Tenant:        tenant,
		UpdatedAt:     updatedAt,
		SourceOfTruth: sourceOfTruth,
	}
	event := eventstore.NewBaseEvent(aggregate, ContactUpdated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

// FIXME alexb implement
//type ContactDeletedEvent struct {
//	Uuid string `json:"uuid" bson:"uuid,omitempty" validate:"required"`
//}
//
//func NewContactDeletedEvent(aggregate eventstore.Aggregate, uuid string) (eventstore.Event, error) {
//	eventData := ContactDeletedEvent{
//		Uuid: uuid,
//	}
//	event := eventstore.NewBaseEvent(aggregate, ContactDeleted)
//	err := event.SetJsonData(&eventData)
//	if err != nil {
//		return eventstore.Event{}, err
//	}
//	return event, nil
//}
