package invoice

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

const (
	InvoiceAggregateType eventstore.AggregateType = "invoice"
)

type InvoiceAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Invoice *Invoice
}

func (a *InvoiceAggregate) IsTemporal() bool {
	if a.Invoice.DryRun {
		return true
	}
	return false
}

func (a *InvoiceAggregate) PrepareStreamMetadata() esdb.StreamMetadata {
	streamMetadata := esdb.StreamMetadata{}
	// set duration for 1 year
	streamMetadata.SetMaxAge(time.Duration(int64(365*24)) * time.Hour)
	return streamMetadata
}

func GetInvoiceObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, InvoiceAggregateType)
}

func NewInvoiceAggregateWithTenantAndID(tenant, id string) *InvoiceAggregate {
	invoiceAggregate := InvoiceAggregate{}
	invoiceAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(InvoiceAggregateType, tenant, id)
	invoiceAggregate.SetWhen(invoiceAggregate.When)
	invoiceAggregate.Invoice = &Invoice{}
	invoiceAggregate.Tenant = tenant

	return &invoiceAggregate
}

func (a *InvoiceAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.HandleGRPCRequest")
	defer span.Finish()

	invoiceNumber := ""
	if params != nil {
		if _, ok := params[PARAM_INVOICE_NUMBER]; ok {
			invoiceNumber = params[PARAM_INVOICE_NUMBER].(string)
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
	case *invoicepb.PayInvoiceNotificationRequest:
		return nil, a.CreatePayInvoiceNotificationEvent(ctx, r)
	case *invoicepb.RemindInvoiceNotificationRequest:
		return nil, a.CreateRemindInvoiceNotificationEvent(ctx, r)
	case *invoicepb.RequestFillInvoiceRequest:
		return nil, a.CreateFillRequestedEvent(ctx, r)
	case *invoicepb.PermanentlyDeleteInitializedInvoiceRequest:
		return nil, a.PermanentlyDeleteInitializedInvoice(ctx, r)
	case *invoicepb.VoidInvoiceRequest:
		return nil, a.VoidInvoice(ctx, r)
	default:
		return nil, nil
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

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
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

	sourceFields := common.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	periodStartDate := utils.TimestampProtoToTimePtr(request.InvoicePeriodStart)
	periodEndDate := utils.TimestampProtoToTimePtr(request.InvoicePeriodEnd)

	createEvent, err := NewInvoiceForContractCreateEvent(a, sourceFields, request.ContractId, request.Currency, request.Note, request.BillingCycleInMonths, request.DryRun, request.OffCycle, request.Postpaid, request.Preview, createdAtNotNil, *periodStartDate, *periodEndDate)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceCreateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&createEvent, span, eventstore.EventMetadata{
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

	// if invoice fill event already applied, return
	for _, appliedEvent := range a.AppliedEvents {
		if appliedEvent.GetEventType() == InvoiceFillV1 {
			err := errors.New("invoice already filled")
			tracing.TraceErr(span, err)
			return nil
		}
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
			SourceFields: common.Source{
				Source:    events2.SourceOpenline,
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

	eventstore.EnrichEventWithMetadataExtended(&fillEvent, span, eventstore.EventMetadata{
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

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
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

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
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

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(event)
}

func (a *InvoiceAggregate) CreateRemindInvoiceNotificationEvent(ctx context.Context, r *invoicepb.RemindInvoiceNotificationRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.CreateRemindInvoiceNotificationEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := NewInvoiceRemindNotificationEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInvoiceRemindNotificationEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: r.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(event)
}

func (a *InvoiceAggregate) PermanentlyDeleteInitializedInvoice(ctx context.Context, request *invoicepb.PermanentlyDeleteInitializedInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.PermanentlyDeleteInitializedInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	if a.Invoice == nil {
		err := errors.New("invoice is nil")
		tracing.TraceErr(span, err)
		return err
	}
	if a.Invoice.Status != neo4jenum.InvoiceStatusInitialized.String() {
		err := errors.New("invoice status is not initialized")
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

	eventstore.EnrichEventWithMetadataExtended(&deleteEvent, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	streamMetadata := esdb.StreamMetadata{}
	streamMetadata.SetMaxAge(time.Duration(events2.StreamMetadataMaxAgeSecondsExtended) * time.Second)
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

	eventstore.EnrichEventWithMetadataExtended(&voidEvent, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	streamMetadata := esdb.StreamMetadata{}
	streamMetadata.SetMaxAge(time.Duration(events2.StreamMetadataMaxAgeSecondsExtended) * time.Second)
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
	case InvoiceUpdateV1,
		InvoicePayV1,
		InvoicePdfRequestedV1,
		InvoiceFillRequestedV1,
		InvoicePaidV1,
		InvoicePayNotificationV1,
		InvoiceRemindNotificationV1,
		InvoiceDeleteV1,
		InvoiceVoidV1:
		return nil
	default:
		return nil
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
	a.Invoice.Preview = eventData.Preview
	a.Invoice.OffCycle = eventData.OffCycle
	a.Invoice.Postpaid = eventData.Postpaid
	a.Invoice.Currency = eventData.Currency
	a.Invoice.PeriodStartDate = eventData.PeriodStartDate
	a.Invoice.PeriodEndDate = eventData.PeriodEndDate
	a.Invoice.BillingCycleInMonths = eventData.BillingCycleInMonths
	a.Invoice.Note = eventData.Note
	a.Invoice.Status = neo4jenum.InvoiceStatusInitialized.String()
	if eventData.BillingCycle != "" {
		switch eventData.BillingCycle {
		case "MONTHLY":
			a.Invoice.BillingCycleInMonths = 1
		case "QUARTERLY":
			a.Invoice.BillingCycleInMonths = 3
		case "ANNUALLY":
			a.Invoice.BillingCycleInMonths = 12
		}
	}

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
