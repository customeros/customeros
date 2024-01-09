package invoicing_cycle

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

const (
	InvoicingCycleCreateV1 = "V1_INVOICING_CYCLE_CREATE"
	InvoicingCycleUpdateV1 = "V1_INVOICING_CYCLE_UPDATE"
)

type InvoicingCycleCreateEvent struct {
	Tenant       string                   `json:"tenant" validate:"required"`
	Type         InvoicingCycleTypeString `json:"type"`
	CreatedAt    time.Time                `json:"createdAt"`
	SourceFields commonmodel.Source       `json:"sourceFields"`
}

func NewInvoicingCycleCreateEvent(aggregate eventstore.Aggregate, invoicingCycleType InvoicingCycleType, sourceFields commonmodel.Source) (eventstore.Event, error) {
	eventData := InvoicingCycleCreateEvent{
		Tenant:       aggregate.GetTenant(),
		Type:         invoicingCycleType.StringValue(),
		CreatedAt:    utils.Now(),
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
	Tenant       string                   `json:"tenant" validate:"required"`
	Type         InvoicingCycleTypeString `json:"type,omitempty"`
	UpdatedAt    time.Time                `json:"updatedAt"`
	SourceFields commonmodel.Source       `json:"sourceFields"`
}

func NewInvoicingCycleUpdateEvent(aggregate eventstore.Aggregate, invoicingCycleType InvoicingCycleType, sourceFields commonmodel.Source) (eventstore.Event, error) {
	eventData := InvoicingCycleUpdateEvent{
		Tenant:       aggregate.GetTenant(),
		UpdatedAt:    utils.Now(),
		Type:         invoicingCycleType.StringValue(),
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
