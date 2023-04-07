package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

const (
	EmailCreated = "EMAIL_CREATED"
	EmailUpdated = "EMAIL_UPDATED"
)

type EmailCreatedEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	RawEmail      string    `json:"rawEmail" validate:"required"`
	Source        string    `json:"source"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	AppSource     string    `json:"appSource"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewEmailCreatedEvent(aggregate eventstore.Aggregate, tenant, rawEmail, source, sourceOfTruth, appSource string, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := EmailCreatedEvent{
		Tenant:        tenant,
		RawEmail:      rawEmail,
		Source:        source,
		SourceOfTruth: sourceOfTruth,
		AppSource:     appSource,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
	event := eventstore.NewBaseEvent(aggregate, EmailCreated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type EmailUpdatedEvent struct {
	Tenant        string    `json:"tenant"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewEmailUpdatedEvent(aggregate eventstore.Aggregate, tenant, sourceOfTruth string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := EmailUpdatedEvent{
		Tenant:        tenant,
		SourceOfTruth: sourceOfTruth,
		UpdatedAt:     updatedAt,
	}
	event := eventstore.NewBaseEvent(aggregate, EmailUpdated)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
