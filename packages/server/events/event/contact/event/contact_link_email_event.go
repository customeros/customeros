package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type ContactLinkEmailEvent struct {
	Tenant    string    `json:"tenant" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
	EmailId   string    `json:"emailId" validate:"required"`
	Primary   bool      `json:"primary"`
	Email     string    `json:"email"`
}

func NewContactLinkEmailEvent(aggregate eventstore.Aggregate, emailId, email string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContactLinkEmailEvent{
		Tenant:    aggregate.GetTenant(),
		UpdatedAt: updatedAt,
		EmailId:   emailId,
		Email:     email,
		Primary:   primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContactLinkEmailEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, ContactEmailLinkV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContactLinkEmailEvent")
	}
	return event, nil
}
