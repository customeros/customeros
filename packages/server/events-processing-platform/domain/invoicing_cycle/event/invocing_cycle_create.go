package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type InvoicingCycleCreateEvent struct {
	Tenant       string                         `json:"tenant" validate:"required"`
	Type         model.InvoicingCycleTypeString `json:"type"`
	CreatedAt    time.Time                      `json:"createdAt"`
	SourceFields commonmodel.Source             `json:"sourceFields"`
}

func NewInvoicingCycleCreateEvent(aggregate eventstore.Aggregate, invoicingCycleType model.InvoicingCycleType, sourceFields commonmodel.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := InvoicingCycleCreateEvent{
		Tenant:       aggregate.GetTenant(),
		Type:         invoicingCycleType.StringValue(),
		CreatedAt:    createdAt,
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
