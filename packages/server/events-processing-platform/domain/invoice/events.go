package invoice

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/pkg/errors"
	"time"
)

const (
	InvoiceCreateV1       = "V1_INVOICE_CREATE"
	InvoiceFillV1         = "V1_INVOICE_FILL"
	InvoicePdfGeneratedV1 = "V1_INVOICE_PDF_GENERATED"
	InvoicePayV1          = "V1_INVOICE_PAY"
)

type InvoiceCreateEvent struct {
	Tenant       string             `json:"tenant" validate:"required"`
	ContractId   string             `json:"organizationId" validate:"required"`
	CreatedAt    time.Time          `json:"createdAt"`
	SourceFields commonmodel.Source `json:"sourceFields"`

	DryRun      bool                    `json:"dryRun"`
	DryRunLines []DryRunServiceLineItem `json:"dryRunLines"`

	Number  string    `json:"number"`
	Date    time.Time `json:"date" validate:"required"`
	DueDate time.Time `json:"dueDate" validate:"required"`
}

func NewInvoiceCreateEvent(aggregate eventstore.Aggregate, sourceFields commonmodel.Source, request *invoicepb.NewOnCycleInvoiceForContractRequest) (eventstore.Event, error) {
	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	dateNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.InvoicePeriodStart), utils.Now())
	eventData := InvoiceCreateEvent{
		Tenant:       aggregate.GetTenant(),
		ContractId:   request.ContractId,
		CreatedAt:    createdAtNotNil,
		SourceFields: sourceFields,

		DryRun:  request.DryRun,
		Number:  uuid.New().String(), // todo logic for number generation
		Date:    dateNotNil,
		DueDate: dateNotNil,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceCreateEvent")
	}

	return event, nil
}

func SimulateInvoiceNewEvent(aggregate eventstore.Aggregate, sourceFields commonmodel.Source, request *invoicepb.SimulateInvoiceRequest) (eventstore.Event, error) {
	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	dateNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.Date), utils.Now())
	eventData := InvoiceCreateEvent{
		Tenant:       aggregate.GetTenant(),
		ContractId:   request.ContractId,
		CreatedAt:    createdAtNotNil,
		SourceFields: sourceFields,

		DryRun:  true,
		Number:  uuid.New().String(),
		Date:    dateNotNil,
		DueDate: dateNotNil,
	}

	eventData.DryRunLines = make([]DryRunServiceLineItem, len(request.DryRunServiceLineItems))
	for i, line := range request.DryRunServiceLineItems {
		eventData.DryRunLines[i] = DryRunServiceLineItem{
			ServiceLineItemId: line.ServiceLineItemId,
			Name:              line.Name,
			Billed:            line.Billed.String(),
			Price:             line.Price,
			Quantity:          line.Quantity,
		}
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceCreateEvent")
	}

	return event, nil
}

type InvoiceFillEvent struct {
	Tenant       string             `json:"tenant" validate:"required"`
	UpdatedAt    time.Time          `json:"createdAt"`
	SourceFields commonmodel.Source `json:"sourceFields"`

	Amount float64                `json:"amount" validate:"required"`
	VAT    float64                `json:"vat" validate:"required"`
	Total  float64                `json:"total" validate:"required"`
	Lines  []InvoiceLineFillEvent `json:"invoiceLines" validate:"required"`
}

type InvoiceLineFillEvent struct {
	Tenant   string  `json:"tenant" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Quantity int64   `json:"quantity" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	VAT      float64 `json:"vat" validate:"required"`
	Total    float64 `json:"total" validate:"required"`
}

func NewInvoiceFillEvent(aggregate eventstore.Aggregate, updatedAt *time.Time, sourceFields commonmodel.Source, request *invoicepb.FillInvoiceRequest) (eventstore.Event, error) {
	eventData := InvoiceFillEvent{
		Tenant:       aggregate.GetTenant(),
		UpdatedAt:    *updatedAt,
		SourceFields: sourceFields,

		Amount: request.Amount,
		VAT:    request.Vat,
		Total:  request.Total,
		Lines:  make([]InvoiceLineFillEvent, len(request.Lines)),
	}
	for i, line := range request.Lines {
		eventData.Lines[i] = InvoiceLineFillEvent{
			Name:     line.Name,
			Price:    line.Price,
			Quantity: line.Quantity,
			Amount:   line.Amount,
			VAT:      line.Vat,
			Total:    line.Total,
		}
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InvoiceFillEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InvoiceFillV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InvoiceFillEvent")
	}

	return event, nil
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
