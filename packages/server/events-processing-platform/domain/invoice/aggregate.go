package invoice

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/customeros/customeros/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
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

func (a *InvoiceAggregate) When(evt eventstore.Event) error {
	return nil
}
