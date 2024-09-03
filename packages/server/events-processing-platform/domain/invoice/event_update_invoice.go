package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type InvoiceUpdateEvent struct {
	Tenant                string     `json:"tenant" validate:"required"`
	UpdatedAt             time.Time  `json:"updatedAt"`
	Status                string     `json:"status,omitempty"`
	PaymentLink           string     `json:"paymentLink,omitempty"`
	PaymentLinkValidUntil *time.Time `json:"paymentLinkValidUntil,omitempty"`
	FieldsMask            []string   `json:"fieldsMask,omitempty"`
}

func NewInvoiceUpdateEvent(aggregate eventstore.Aggregate, updatedAt time.Time, fieldsMask []string, status, paymentLink string, paymentLinkValidUntil *time.Time) (eventstore.Event, error) {
	eventData := InvoiceUpdateEvent{
		Tenant:     aggregate.GetTenant(),
		UpdatedAt:  updatedAt,
		FieldsMask: fieldsMask,
	}
	if eventData.UpdateStatus() {
		eventData.Status = status
	}
	if eventData.UpdatePaymentLink() {
		eventData.PaymentLink = paymentLink
		eventData.PaymentLinkValidUntil = paymentLinkValidUntil
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceUpdateEvent")
	}

	return event, nil
}

func (e InvoiceUpdateEvent) UpdateStatus() bool {
	return utils.Contains(e.FieldsMask, FieldMaskStatus)
}

func (e InvoiceUpdateEvent) UpdatePaymentLink() bool {
	return utils.Contains(e.FieldsMask, FieldMaskPaymentLink)
}
