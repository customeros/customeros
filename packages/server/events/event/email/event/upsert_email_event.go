package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type EmailUpsertEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	RawEmail  string    `json:"rawEmail"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewEmailUpsertEvent(aggregate eventstore.Aggregate, tenant, rawEmail, source string, createdAt time.Time) (eventstore.Event, error) {
	eventData := EmailUpsertEvent{
		Tenant:    tenant,
		RawEmail:  rawEmail,
		Source:    source,
		CreatedAt: createdAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate EmailUpsertEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, EmailUpsertV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for EmailUpsertEvent")
	}
	return event, nil
}
