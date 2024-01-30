package invoice

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/pkg/errors"
	"time"
)

const (
	InvoiceCreateForContractV1 = "V1_INVOICE_CREATE_FOR_CONTRACT"
	InvoiceFillV1              = "V1_INVOICE_FILL"
	InvoicePdfRequestedV1      = "V1_INVOICE_PDF_REQUESTED"
	InvoicePdfGeneratedV1      = "V1_INVOICE_PDF_GENERATED"
	InvoicePayV1               = "V1_INVOICE_PAY"
	InvoiceUpdateV1            = "V1_INVOICE_UPDATE"
	InvoicePaidV1              = "V1_INVOICE_PAID"
	InvoicePayNotificationV1   = "V1_INVOICE_PAY_NOTIFICATION"
)

const (
	FieldMaskStatus      = "status"
	FieldMaskPaymentLink = "paymentLink"
)

func SimulateInvoiceNewEvent(aggregate eventstore.Aggregate, sourceFields commonmodel.Source, request *invoicepb.SimulateInvoiceRequest) (eventstore.Event, error) {
	//createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	//dateNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.Date), utils.Now())
	//eventData := InvoiceCreateEvent{
	//	Tenant:       aggregate.GetTenant(),
	//	ContractId:   request.ContractId,
	//	CreatedAt:    createdAtNotNil,
	//	SourceFields: sourceFields,
	//
	//	DryRun: true,
	//	//Number:  uuid.New().String(),
	//	//Date:    dateNotNil,
	//	//DueDate: dateNotNil,
	//}

	//eventData.DryRunLines = make([]DryRunServiceLineItem, len(request.DryRunServiceLineItems))
	//for i, line := range request.DryRunServiceLineItems {
	//	eventData.DryRunLines[i] = DryRunServiceLineItem{
	//		ServiceLineItemId: line.ServiceLineItemId,
	//		Name:              line.Name,
	//		Billed:            line.Billed.String(),
	//		Price:             line.Price,
	//		Quantity:          line.Quantity,
	//	}
	//}

	//if err := validator.GetValidator().Struct(eventData); err != nil {
	//	return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceCreateEvent")
	//}
	//
	//event := eventstore.NewBaseEvent(aggregate, InvoiceCreateForContractV1)
	//if err := event.SetJsonData(&eventData); err != nil {
	//	return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceCreateEvent")
	//}

	//return event, nil
	return eventstore.Event{}, nil
}

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

type InvoicePayEvent struct {
	Tenant       string             `json:"tenant" validate:"required"`
	UpdatedAt    time.Time          `json:"createdAt"`
	SourceFields commonmodel.Source `json:"sourceFields"`
}

func NewInvoicePayEvent(aggregate eventstore.Aggregate, updatedAt *time.Time, sourceFields commonmodel.Source, request *invoicepb.PayInvoiceRequest) (eventstore.Event, error) {
	eventData := InvoicePayEvent{
		Tenant:       aggregate.GetTenant(),
		UpdatedAt:    *updatedAt,
		SourceFields: sourceFields,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoicePayEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoicePayV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoicePayEvent")
	}

	return event, nil
}
