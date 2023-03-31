package events

import (
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
)

const (
	PhoneNumberCreated = "PHONE_NUMBER_CREATED"
)

type PhoneNumberCreatedEvent struct {
	Tenant         string `json:"tenant" validate:"required"`
	RawPhoneNumber string `json:"rawPhoneNumber" validate:"required"`
}

func NewPhoneNumberCreatedEvent(aggregate eventstore.Aggregate, tenant, rawPhoneNumber string) (eventstore.Event, error) {
	eventData := PhoneNumberCreatedEvent{
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
	}
	event := eventstore.NewBaseEvent(aggregate, PhoneNumberCreated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
