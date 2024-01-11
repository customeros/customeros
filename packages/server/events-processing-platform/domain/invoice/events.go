package invoice

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

const (
	InvoiceNewV1  = "V1_INVOICE_NEW"
	InvoiceFillV1 = "V1_INVOICE_FILL"
	InvoicePayV1  = "V1_INVOICE_PAY"
)

type InvoiceNewEvent struct {
	Tenant         string             `json:"tenant" validate:"required"`
	OrganizationId string             `json:"organizationId"`
	CreatedAt      time.Time          `json:"createdAt"`
	SourceFields   commonmodel.Source `json:"sourceFields"`
}

func NewInvoiceNewEvent(aggregate eventstore.Aggregate, organizationId string, createdAt *time.Time, sourceFields commonmodel.Source) (eventstore.Event, error) {
	eventData := InvoiceNewEvent{
		Tenant:         aggregate.GetTenant(),
		OrganizationId: organizationId,
		CreatedAt:      *createdAt,
		SourceFields:   sourceFields,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceNewEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceNewV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceNewEvent")
	}

	return event, nil
}

type EventHandlers struct {
	InvoiceNew InvoiceNewHandler
}

func NewEventHandlers(log logger.Logger, es eventstore.AggregateStore) *EventHandlers {
	return &EventHandlers{
		InvoiceNew: NewInvoiceNewHandler(log, es),
	}
}
