package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type InvoiceUpdateEvent struct {
	Tenant     string    `json:"tenant" validate:"required"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Status     string    `json:"status,omitempty"`
	FieldsMask []string  `json:"fieldsMask,omitempty"`
}

func NewInvoiceUpdateEvent(aggregate eventstore.Aggregate, updatedAt time.Time, fieldsMask []string, status string) (eventstore.Event, error) {
	eventData := InvoiceUpdateEvent{
		Tenant:     aggregate.GetTenant(),
		UpdatedAt:  updatedAt,
		FieldsMask: fieldsMask,
	}
	if eventData.UpdateStatus() {
		eventData.Status = status
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
	return len(e.FieldsMask) == 0 || utils.Contains(e.FieldsMask, FieldMaskStatus)
}
