package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	FieldMaskLegalName = "legalName"
	FieldMaskTaxId     = "taxId"
)

type BillingProfileUpdateEvent struct {
	Tenant           string    `json:"tenant" validate:"required"`
	BillingProfileId string    `json:"billingProfileId" validate:"required"`
	UpdatedAt        time.Time `json:"updatedAt"`
	LegalName        string    `json:"legalName,omitempty"`
	TaxId            string    `json:"taxId,omitempty"`
	FieldsMask       []string  `json:"fieldsMask,omitempty"`
}

func NewBillingProfileUpdateEvent(aggregate eventstore.Aggregate, billingProfileId, legalName, taxId string, updatedAt time.Time, fieldsMask []string) (eventstore.Event, error) {
	eventData := BillingProfileUpdateEvent{
		Tenant:           aggregate.GetTenant(),
		BillingProfileId: billingProfileId,
		UpdatedAt:        updatedAt,
		FieldsMask:       fieldsMask,
	}
	if eventData.UpdateLegalName() {
		eventData.LegalName = legalName
	}
	if eventData.UpdateTaxId() {
		eventData.TaxId = taxId
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate BillingProfileUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationUpdateBillingProfileV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for BillingProfileUpdateEvent")
	}
	return event, nil
}

func (e BillingProfileUpdateEvent) UpdateLegalName() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskLegalName)
}

func (e BillingProfileUpdateEvent) UpdateTaxId() bool {
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskTaxId)
}
