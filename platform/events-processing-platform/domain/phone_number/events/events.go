package events

import (
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
	"time"
)

const (
	PhoneNumberCreated = "PHONE_NUMBER_CREATED"
	PhoneNumberUpdated = "PHONE_NUMBER_UPDATED"
)

type PhoneNumberCreatedEvent struct {
	Tenant         string    `json:"tenant" validate:"required"`
	RawPhoneNumber string    `json:"rawPhoneNumber" validate:"required"`
	Source         string    `json:"source"`
	SourceOfTruth  string    `json:"sourceOfTruth"`
	AppSource      string    `json:"appSource"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func NewPhoneNumberCreatedEvent(aggregate eventstore.Aggregate, tenant, rawPhoneNumber, source, sourceOfTruth, appSource string, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := PhoneNumberCreatedEvent{
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
		Source:         source,
		SourceOfTruth:  sourceOfTruth,
		AppSource:      appSource,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
	event := eventstore.NewBaseEvent(aggregate, PhoneNumberCreated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type PhoneNumberUpdatedEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewPhoneNumberUpdatedEvent(aggregate eventstore.Aggregate, tenant, sourceOfTruth string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := PhoneNumberUpdatedEvent{
		Tenant:        tenant,
		SourceOfTruth: sourceOfTruth,
		UpdatedAt:     updatedAt,
	}
	event := eventstore.NewBaseEvent(aggregate, PhoneNumberUpdated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
