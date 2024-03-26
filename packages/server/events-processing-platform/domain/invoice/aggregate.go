package invoice

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
	"time"
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

func (a *InvoiceAggregate) HandleRequest(ctx context.Context, request any, params ...map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.HandleRequest")
	defer span.Finish()

	invoiceNumber := ""
	if params != nil {
		if _, ok := params[0][PARAM_INVOICE_NUMBER]; ok {
			invoiceNumber = params[0][PARAM_INVOICE_NUMBER].(string)
		}
	}

	switch r := request.(type) {
	case *invoicepb.NewInvoiceForContractRequest:
		return nil, a.CreateNewInvoiceForContract(ctx, r)
	case *invoicepb.FillInvoiceRequest:
		return nil, a.FillInvoice(ctx, r, invoiceNumber)
	case *invoicepb.GenerateInvoicePdfRequest:
		return nil, a.CreatePdfRequestedEvent(ctx, r)
	case *invoicepb.PdfGeneratedInvoiceRequest:
		return nil, a.CreatePdfGeneratedEvent(ctx, r)
	case *invoicepb.PayInvoiceRequest:
		return nil, a.PayInvoice(ctx, r)
	case *invoicepb.UpdateInvoiceRequest:
		return nil, a.UpdateInvoice(ctx, r)
	case *invoicepb.PayInvoiceNotificationRequest:
		return nil, a.CreatePayInvoiceNotificationEvent(ctx, r)
	case *invoicepb.RequestFillInvoiceRequest:
		return nil, a.CreateFillRequestedEvent(ctx, r)
	case *invoicepb.PermanentlyDeleteDraftInvoiceRequest:
		return nil, a.PermanentlyDeleteDraftInvoice(ctx, r)
	case *invoicepb.VoidInvoiceRequest:
		return nil, a.VoidInvoice(ctx, r)
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

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	periodStartDate := utils.TimestampProtoToTimePtr(request.InvoicePeriodStart)
	periodEndDate := utils.TimestampProtoToTimePtr(request.InvoicePeriodEnd)
	billingCycle := BillingCycle(request.BillingCycle)

	createEvent, err := NewInvoiceForContractCreateEvent(a, sourceFields, request.ContractId, request.Currency, billingCycle.String(), request.Note, request.DryRun, request.OffCycle, request.Postpaid, createdAtNotNil, *periodStartDate, *periodEndDate)
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

func (a *InvoiceAggregate) FillInvoice(ctx context.Context, request *invoicepb.FillInvoiceRequest, invoiceNumber string) error {
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

	invoiceNumberForEvent := invoiceNumber
	if a.Invoice.InvoiceNumber != "" {
		invoiceNumberForEvent = a.Invoice.InvoiceNumber
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

	invoiceStatus := InvoiceStatus(request.Status).String()
	fillEvent, err := NewInvoiceFillEvent(a, updatedAtNotNil, *a.Invoice,
		request.Customer.Name, request.Customer.AddressLine1, request.Customer.AddressLine2, request.Customer.Zip, request.Customer.Locality, request.Customer.Country, request.Customer.Region, request.Customer.Email,
		request.Provider.LogoRepositoryFileId, request.Provider.Name, request.Provider.Email, request.Provider.AddressLine1, request.Provider.AddressLine2, request.Provider.Zip, request.Provider.Locality, request.Provider.Country, request.Provider.Region,
		request.Note, invoiceStatus, invoiceNumberForEvent, request.Amount, request.Vat, request.Total, invoiceLines)
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

func (a *InvoiceAggregate) CreateFillRequestedEvent(ctx context.Context, r *invoicepb.RequestFillInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.CreateFillRequestedEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := NewInvoiceFillRequestedEvent(a, r.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInvoiceFillRequestedEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(event)
}

func (a *InvoiceAggregate) CreatePayInvoiceNotificationEvent(ctx context.Context, r *invoicepb.PayInvoiceNotificationRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.CreatePayInvoiceNotificationEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := NewInvoicePayNotificationEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInvoicePayNotificationEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(event)
}

func (a *InvoiceAggregate) UpdateInvoice(ctx context.Context, r *invoicepb.UpdateInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.UpdateInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.UpdatedAt), utils.Now())
	fieldsMask := extractFieldsMask(r.FieldsMask)
	status := InvoiceStatus(r.Status).String()

	events := []eventstore.Event{}
	updateEvent, err := NewInvoiceUpdateEvent(a, updatedAtNotNil, fieldsMask, status, r.PaymentLink)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUpdateInvoiceEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})
	events = append(events, updateEvent)

	// if status updated, and set from non-paid to paid
	if len(fieldsMask) == 0 || utils.Contains(fieldsMask, FieldMaskStatus) &&
		a.Invoice.Status != neo4jenum.InvoiceStatusPaid.String() &&
		status == neo4jenum.InvoiceStatusPaid.String() {
		paidEvent, err := NewInvoicePaidEvent(a)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewInvoicePaidEvent")
		}
		aggregate.EnrichEventWithMetadataExtended(&paidEvent, span, aggregate.EventMetadata{
			Tenant: r.Tenant,
			UserId: r.LoggedInUserId,
			App:    r.AppSource,
		})
		events = append(events, paidEvent)
	}

	return a.ApplyAll(events)
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

func (a *InvoiceAggregate) PermanentlyDeleteDraftInvoice(ctx context.Context, request *invoicepb.PermanentlyDeleteDraftInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.PermanentlyDeleteDraftInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	if a.Invoice == nil {
		err := errors.New("invoice is nil")
		tracing.TraceErr(span, err)
		return err
	}
	if a.Invoice.Status != neo4jenum.InvoiceStatusDraft.String() {
		err := errors.New("invoice status is not draft")
		tracing.TraceErr(span, err)
		return err
	}
	if len(a.Invoice.InvoiceLines) > 0 {
		err := errors.New("invoice has invoice lines")
		tracing.TraceErr(span, err)
		return err
	}
	deleteEvent, err := NewInvoiceDeleteEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoicePayEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&deleteEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	streamMetadata := esdb.StreamMetadata{}
	streamMetadata.SetMaxAge(time.Duration(constants.StreamMetadataMaxAgeSecondsExtended) * time.Second)
	a.SetStreamMetadata(&streamMetadata)

	return a.Apply(deleteEvent)
}

func (a *InvoiceAggregate) VoidInvoice(ctx context.Context, request *invoicepb.VoidInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.VoidInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	if a.Invoice == nil {
		err := errors.New("invoice is nil")
		tracing.TraceErr(span, err)
		return err
	}

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())

	voidEvent, err := NewInvoiceVoidEvent(a, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceVoidEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&voidEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	streamMetadata := esdb.StreamMetadata{}
	streamMetadata.SetMaxAge(time.Duration(constants.StreamMetadataMaxAgeSecondsExtended) * time.Second)
	a.SetStreamMetadata(&streamMetadata)

	return a.Apply(voidEvent)
}

func (a *InvoiceAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case InvoiceCreateForContractV1:
		return a.onInvoiceCreateEvent(evt)
	case InvoiceFillV1:
		return a.onFillInvoice(evt)
	case InvoicePdfGeneratedV1:
		return a.onPdfGeneratedInvoice(evt)
	case InvoiceUpdateV1:
		return a.onUpdateInvoice(evt)
	case InvoicePayV1,
		InvoicePdfRequestedV1,
		InvoiceFillRequestedV1,
		InvoicePaidV1,
		InvoicePayNotificationV1,
		InvoiceDeleteV1,
		InvoiceVoidV1:
		return nil
	default:
		if strings.HasPrefix(evt.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
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
	a.Invoice.DryRun = eventData.DryRun
	a.Invoice.OffCycle = eventData.OffCycle
	a.Invoice.Postpaid = eventData.Postpaid
	a.Invoice.Currency = eventData.Currency
	a.Invoice.PeriodStartDate = eventData.PeriodStartDate
	a.Invoice.PeriodEndDate = eventData.PeriodEndDate
	a.Invoice.BillingCycle = eventData.BillingCycle
	a.Invoice.Note = eventData.Note
	a.Invoice.Status = neo4jenum.InvoiceStatusDraft.String()

	return nil
}

func (a *InvoiceAggregate) onFillInvoice(evt eventstore.Event) error {
	var eventData InvoiceFillEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Invoice.InvoiceNumber = eventData.InvoiceNumber
	a.Invoice.UpdatedAt = eventData.UpdatedAt
	a.Invoice.Amount = eventData.Amount
	a.Invoice.VAT = eventData.VAT
	a.Invoice.TotalAmount = eventData.TotalAmount
	a.Invoice.Note = eventData.Note
	if eventData.Status != "" {
		a.Invoice.Status = eventData.Status
	}
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

func (a *InvoiceAggregate) onUpdateInvoice(evt eventstore.Event) error {
	var eventData InvoiceUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if eventData.UpdateStatus() {
		a.Invoice.Status = eventData.Status
	}
	if eventData.UpdatePaymentLink() {
		a.Invoice.PaymentLink = eventData.PaymentLink
	}

	return nil
}

func extractFieldsMask(requestFieldsMask []invoicepb.InvoiceFieldMask) []string {
	var fieldsMask []string
	for _, field := range requestFieldsMask {
		switch field {
		case invoicepb.InvoiceFieldMask_INVOICE_FIELD_STATUS:
			fieldsMask = append(fieldsMask, FieldMaskStatus)
		case invoicepb.InvoiceFieldMask_INVOICE_FIELD_PAYMENT_LINK:
			fieldsMask = append(fieldsMask, FieldMaskPaymentLink)
		}
	}
	fieldsMask = utils.RemoveDuplicates(fieldsMask)
	return fieldsMask
}
