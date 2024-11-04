package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	InvoiceCreateForContractV1 = "V1_INVOICE_CREATE_FOR_CONTRACT"
	InvoiceFillRequestedV1     = "V1_INVOICE_FILL_REQUESTED"
	InvoiceFillV1              = "V1_INVOICE_FILL"
	InvoicePdfRequestedV1      = "V1_INVOICE_PDF_REQUESTED"
	InvoicePdfGeneratedV1      = "V1_INVOICE_PDF_GENERATED"
	// Deprecated
	InvoiceUpdateV1             = "V1_INVOICE_UPDATE"
	InvoicePaidV1               = "V1_INVOICE_PAID"
	InvoicePayNotificationV1    = "V1_INVOICE_PAY_NOTIFICATION"
	InvoiceRemindNotificationV1 = "V1_INVOICE_REMIND_NOTIFICATION"
	InvoiceDeleteV1             = "V1_INVOICE_DELETE"
	InvoiceVoidV1               = "V1_INVOICE_VOID"
	// Deprecated
	InvoicePayV1 = "V1_INVOICE_PAY"
)

const (
	FieldMaskStatus      = "status"
	FieldMaskPaymentLink = "paymentLink"
)

type InvoicePdfGeneratedEvent struct {
	Tenant           string    `json:"tenant" validate:"required"`
	UpdatedAt        time.Time `json:"updatedAt"`
	RepositoryFileId string    `json:"repositoryFileId" validate:"required"`
}

func NewInvoicePdfGeneratedEvent(aggregate eventstore.Aggregate, updatedAt time.Time, repositoryFileId string) (eventstore.Event, error) {
	eventData := InvoicePdfGeneratedEvent{
		Tenant:           aggregate.GetTenant(),
		UpdatedAt:        updatedAt,
		RepositoryFileId: repositoryFileId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoicePdfGeneratedEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoicePdfGeneratedV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoicePdfGeneratedEvent")
	}

	return event, nil
}
