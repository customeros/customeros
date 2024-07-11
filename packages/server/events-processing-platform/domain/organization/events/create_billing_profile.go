package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type BillingProfileCreateEvent struct {
	Tenant           string        `json:"tenant" validate:"required"`
	BillingProfileId string        `json:"billingProfileId" validate:"required"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	LegalName        string        `json:"legalName"`
	TaxId            string        `json:"taxId"`
	SourceFields     events.Source `json:"sourceFields" validate:"required"`
}

func NewBillingProfileCreateEvent(aggregate eventstore.Aggregate, billingProfileId, legalName, taxId string, sourceFields events.Source, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := BillingProfileCreateEvent{
		Tenant:           aggregate.GetTenant(),
		BillingProfileId: billingProfileId,
		LegalName:        legalName,
		TaxId:            taxId,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		SourceFields:     sourceFields,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate BillingProfileCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationCreateBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for BillingProfileCreateEvent")
	}
	return event, nil
}
