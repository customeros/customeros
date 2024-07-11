package invoicing_cycle

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	InvoicingCycleCreateV1 = "V1_INVOICING_CYCLE_CREATE"
	InvoicingCycleUpdateV1 = "V1_INVOICING_CYCLE_UPDATE"
)

type InvoicingCycleCreateEvent struct {
	Tenant       string        `json:"tenant" validate:"required"`
	Type         string        `json:"type"`
	CreatedAt    time.Time     `json:"createdAt"`
	SourceFields events.Source `json:"sourceFields"`
}

func NewInvoicingCycleCreateEvent(aggregate eventstore.Aggregate, invoicingCycleType string, createdAt *time.Time, sourceFields events.Source) (eventstore.Event, error) {
	eventData := InvoicingCycleCreateEvent{
		Tenant:       aggregate.GetTenant(),
		Type:         invoicingCycleType,
		CreatedAt:    *createdAt,
		SourceFields: sourceFields,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoicingCycleCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoicingCycleCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoicingCycleCreateEvent")
	}

	return event, nil
}

type InvoicingCycleUpdateEvent struct {
	Tenant       string        `json:"tenant" validate:"required"`
	Type         string        `json:"type,omitempty"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	SourceFields events.Source `json:"sourceFields"`
}

func NewInvoicingCycleUpdateEvent(aggregate eventstore.Aggregate, invoicingCycleType string, updatedAt *time.Time, sourceFields events.Source) (eventstore.Event, error) {
	eventData := InvoicingCycleUpdateEvent{
		Tenant:       aggregate.GetTenant(),
		UpdatedAt:    *updatedAt,
		Type:         invoicingCycleType,
		SourceFields: sourceFields,
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

type EventHandlers struct {
	CreateInvoicingCycle CreateInvoicingCycleHandler
	UpdateInvoicingCycle UpdateInvoicingCycleHandler
}

func NewEventHandlers(log logger.Logger, es eventstore.AggregateStore) *EventHandlers {
	return &EventHandlers{
		CreateInvoicingCycle: NewCreateInvoicingCycleHandler(log, es),
		UpdateInvoicingCycle: NewUpdateInvoicingCycleHandler(log, es),
	}
}
