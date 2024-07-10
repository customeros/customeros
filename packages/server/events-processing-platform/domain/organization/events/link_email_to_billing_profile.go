package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type LinkEmailToBillingProfileEvent struct {
	Tenant           string    `json:"tenant" validate:"required"`
	UpdatedAt        time.Time `json:"updatedAt"`
	BillingProfileId string    `json:"billingProfileId" validate:"required"`
	EmailId          string    `json:"emailId" validate:"required"`
	Primary          bool      `json:"primary"`
}

func NewLinkEmailToBillingProfileEvent(aggregate eventstore.Aggregate, billingProfileId, emailId string, primary bool, updatedAt time.Time) (eventstore.Event, error) {
	eventData := LinkEmailToBillingProfileEvent{
		Tenant:           aggregate.GetTenant(),
		UpdatedAt:        updatedAt,
		BillingProfileId: billingProfileId,
		EmailId:          emailId,
		Primary:          primary,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate LinkEmailToBillingProfileEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationEmailLinkToBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for LinkEmailToBillingProfileEvent")
	}
	return event, nil
}
