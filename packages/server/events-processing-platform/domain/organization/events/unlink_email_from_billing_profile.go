package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type UnlinkEmailToBillingProfileEvent struct {
	Tenant           string    `json:"tenant" validate:"required"`
	UpdatedAt        time.Time `json:"updatedAt"`
	BillingProfileId string    `json:"billingProfileId" validate:"required"`
	EmailId          string    `json:"emailId" validate:"required"`
}

func NewUnlinkEmailToBillingProfileEvent(aggregate eventstore.Aggregate, billingProfileId, emailId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UnlinkEmailToBillingProfileEvent{
		Tenant:           aggregate.GetTenant(),
		UpdatedAt:        updatedAt,
		BillingProfileId: billingProfileId,
		EmailId:          emailId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UnlinkEmailToBillingProfileEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationEmailUnlinkFromBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UnlinkEmailToBillingProfileEvent")
	}
	return event, nil
}
