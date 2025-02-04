package invoice

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"github.com/customeros/customeros/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	events2 "github.com/customeros/customeros/packages/server/events/constants"
	"github.com/customeros/customeros/packages/server/events/eventstore"
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

	switch r := request.(type) {
	case *invoicepb.PayInvoiceNotificationRequest:
		return nil, a.CreatePayInvoiceNotificationEvent(ctx, r)
	case *invoicepb.RemindInvoiceNotificationRequest:
		return nil, a.CreateRemindInvoiceNotificationEvent(ctx, r)
	case *invoicepb.PermanentlyDeleteInitializedInvoiceRequest:
		return nil, a.PermanentlyDeleteInitializedInvoice(ctx, r)
	case *invoicepb.VoidInvoiceRequest:
		return nil, a.VoidInvoice(ctx, r)
	default:
		return nil, nil
	}
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
	return nil
}
