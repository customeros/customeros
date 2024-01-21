package invoice

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

const (
	InvoiceAggregateType eventstore.AggregateType = "invoice"
)

type InvoiceAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Invoice *Invoice
}

func GetInvoiceObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, InvoiceAggregateType)
}

func LoadInvoiceAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string, options eventstore.LoadAggregateOptions) (*InvoiceAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadInvoiceAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	invoiceAggregate := NewInvoiceAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, invoiceAggregate, options)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return invoiceAggregate, nil
}

func NewInvoiceAggregateWithTenantAndID(tenant, id string) *InvoiceAggregate {
	invoiceAggregate := InvoiceAggregate{}
	invoiceAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(InvoiceAggregateType, tenant, id)
	invoiceAggregate.SetWhen(invoiceAggregate.When)
	invoiceAggregate.Invoice = &Invoice{}
	invoiceAggregate.Tenant = tenant

	return &invoiceAggregate
}

func (a *InvoiceAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *invoicepb.NewInvoiceForContractRequest:
		return nil, a.CreateNewInvoiceForContract(ctx, r)
	case *invoicepb.FillInvoiceRequest:
		return nil, a.FillInvoice(ctx, r)
	case *invoicepb.GenerateInvoicePdfRequest:
		return nil, a.CreatePdfRequestedEvent(ctx, r)
	case *invoicepb.PdfGeneratedInvoiceRequest:
		return nil, a.CreatePdfGeneratedEvent(ctx, r)
	case *invoicepb.PayInvoiceRequest:
		return nil, a.PayInvoice(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *InvoiceAggregate) CreatePdfGeneratedEvent(ctx context.Context, request *invoicepb.PdfGeneratedInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.CreatePdfGeneratedEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())

	event, err := NewInvoicePdfGeneratedEvent(a, updatedAtNotNil, request.RepositoryFileId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInvoicePdfGeneratedEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *InvoiceAggregate) CreateNewInvoiceForContract(ctx context.Context, request *invoicepb.NewInvoiceForContractRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.CreateNewInvoiceForContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	// TODO add here logic to generate invoice number
	invoiceNumber := uuid.New().String()

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	periodStartDate := utils.TimestampProtoToTimePtr(request.InvoicePeriodStart)
	periodEndDate := utils.TimestampProtoToTimePtr(request.InvoicePeriodEnd)
	billingCycle := BillingCycle(request.BillingCycle)

	createEvent, err := NewInvoiceForContractCreateEvent(a, sourceFields, request.ContractId, request.Currency, invoiceNumber, billingCycle.String(), request.DryRun, createdAtNotNil, *periodStartDate, *periodEndDate)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *InvoiceAggregate) FillInvoice(ctx context.Context, request *invoicepb.FillInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.FillInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	if a.Invoice == nil {
		err := errors.New("invoice is nil")
		tracing.TraceErr(span, err)
		return err
	}

	// prepare invoice lines
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())
	var invoiceLines []InvoiceLineEvent

	for _, line := range request.InvoiceLines {
		invoiceLines = append(invoiceLines, InvoiceLineEvent{
			Id:        uuid.New().String(),
			CreatedAt: updatedAtNotNil,
			SourceFields: commonmodel.Source{
				Source:    constants.SourceOpenline,
				AppSource: request.AppSource,
			},
			Name:                    line.Name,
			Price:                   line.Price,
			Quantity:                line.Quantity,
			Amount:                  line.Amount,
			VAT:                     line.Vat,
			TotalAmount:             line.Total,
			ServiceLineItemId:       line.ServiceLineItemId,
			ServiceLineItemParentId: line.ServiceLineItemParentId,
			BilledType:              BilledType(line.BilledType).String(),
		})
	}

	fillEvent, err := NewInvoiceFillEvent(a, updatedAtNotNil, *a.Invoice, request.Amount, request.Vat, request.Total, invoiceLines)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceFillEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&fillEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(fillEvent)
}

func (a *InvoiceAggregate) CreatePdfRequestedEvent(ctx context.Context, r *invoicepb.GenerateInvoicePdfRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.CreatePdfRequestedEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := NewInvoicePdfRequestedEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInvoicePdfRequestedEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(event)
}

func (a *InvoiceAggregate) PayInvoice(ctx context.Context, request *invoicepb.PayInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.PayInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	payEvent, err := NewInvoicePayEvent(a, utils.TimestampProtoToTimePtr(request.UpdatedAt), sourceFields, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoicePayEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&payEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	return a.Apply(payEvent)
}

func (a *InvoiceAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case InvoiceCreateForContractV1:
		return a.onInvoiceCreateEvent(evt)
	case InvoiceFillV1:
		return a.onFillInvoice(evt)
	case InvoicePdfGeneratedV1:
		return a.onPdfGeneratedInvoice(evt)
	case InvoicePayV1:
		return a.onPayInvoice(evt)
	case InvoicePdfRequestedV1:
		return nil
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *InvoiceAggregate) onInvoiceCreateEvent(evt eventstore.Event) error {
	var eventData InvoiceForContractCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Invoice.ID = a.ID
	a.Invoice.CreatedAt = eventData.CreatedAt
	a.Invoice.UpdatedAt = eventData.CreatedAt
	a.Invoice.ContractId = eventData.ContractId
	a.Invoice.SourceFields = eventData.SourceFields
	a.Invoice.InvoiceNumber = eventData.InvoiceNumber
	a.Invoice.DryRun = eventData.DryRun
	a.Invoice.Currency = eventData.Currency
	a.Invoice.PeriodStartDate = eventData.PeriodStartDate
	a.Invoice.PeriodEndDate = eventData.PeriodEndDate
	a.Invoice.BillingCycle = eventData.BillingCycle

	return nil
}

func (a *InvoiceAggregate) onFillInvoice(evt eventstore.Event) error {
	var eventData InvoiceFillEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Invoice.UpdatedAt = eventData.UpdatedAt
	a.Invoice.Amount = eventData.Amount
	a.Invoice.VAT = eventData.VAT
	a.Invoice.TotalAmount = eventData.TotalAmount
	for _, line := range eventData.InvoiceLines {
		a.Invoice.InvoiceLines = append(a.Invoice.InvoiceLines, InvoiceLine{
			Name:                    line.Name,
			Price:                   line.Price,
			Quantity:                line.Quantity,
			Amount:                  line.Amount,
			VAT:                     line.VAT,
			TotalAmount:             line.TotalAmount,
			ServiceLineItemId:       line.ServiceLineItemId,
			ServiceLineItemParentId: line.ServiceLineItemParentId,
			CreatedAt:               line.CreatedAt,
			UpdatedAt:               line.CreatedAt,
			SourceFields:            line.SourceFields,
			BilledType:              line.BilledType,
		})
	}

	return nil
}

func (a *InvoiceAggregate) onPdfGeneratedInvoice(evt eventstore.Event) error {
	var eventData InvoicePdfGeneratedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Invoice.RepositoryFileId = eventData.RepositoryFileId

	return nil
}

func (a *InvoiceAggregate) onPayInvoice(evt eventstore.Event) error {
	return nil
}
