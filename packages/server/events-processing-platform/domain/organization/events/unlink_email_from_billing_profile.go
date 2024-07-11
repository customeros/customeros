package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type UnlinkEmailFromBillingProfileEvent struct {
	Tenant           string    `json:"tenant" validate:"required"`
	UpdatedAt        time.Time `json:"updatedAt"`
	BillingProfileId string    `json:"billingProfileId" validate:"required"`
	EmailId          string    `json:"emailId" validate:"required"`
}

func NewUnlinkEmailFromBillingProfileEvent(aggregate eventstore.Aggregate, billingProfileId, emailId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UnlinkEmailFromBillingProfileEvent{
		Tenant:           aggregate.GetTenant(),
		UpdatedAt:        updatedAt,
		BillingProfileId: billingProfileId,
		EmailId:          emailId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UnlinkEmailFromBillingProfileEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationEmailUnlinkFromBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UnlinkEmailFromBillingProfileEvent")
	}
	return event, nil
}
