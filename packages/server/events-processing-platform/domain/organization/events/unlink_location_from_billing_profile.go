package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type UnlinkLocationFromBillingProfileEvent struct {
	Tenant           string    `json:"tenant" validate:"required"`
	UpdatedAt        time.Time `json:"updatedAt"`
	BillingProfileId string    `json:"billingProfileId" validate:"required"`
	LocationId       string    `json:"locationId" validate:"required"`
}

func NewUnlinkLocationFromBillingProfileEvent(aggregate eventstore.Aggregate, billingProfileId, locationId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UnlinkLocationFromBillingProfileEvent{
		Tenant:           aggregate.GetTenant(),
		UpdatedAt:        updatedAt,
		BillingProfileId: billingProfileId,
		LocationId:       locationId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UnlinkLocationFromBillingProfileEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationLocationUnlinkFromBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UnlinkLocationFromBillingProfileEvent")
	}
	return event, nil
}
