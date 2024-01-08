package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type InvoicingCycleUpdateEvent struct {
	Tenant    string                         `json:"tenant" validate:"required"`
	Type      model.InvoicingCycleTypeString `json:"type,omitempty"`
	UpdatedAt time.Time                      `json:"updatedAt"`
}

func NewInvoicingCycleUpdateEvent(aggregate eventstore.Aggregate, updatedAt time.Time, invoicingCycleType model.InvoicingCycleType) (eventstore.Event, error) {
	eventData := InvoicingCycleUpdateEvent{
		Tenant:    aggregate.GetTenant(),
		UpdatedAt: updatedAt,
		Type:      invoicingCycleType.StringValue(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoicingCycleUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoicingCycleUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoicingCycleUpdateEvent")
	}

	return event, nil
}
