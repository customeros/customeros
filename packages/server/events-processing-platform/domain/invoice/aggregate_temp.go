package invoice

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type InvoiceTempAggregate struct {
	*aggregate.CommonTenantIdTempAggregate
}

func LoadInvoiceTempAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string) (*InvoiceTempAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadInvoiceAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	invoiceAggregate := NewInvoiceTempAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, invoiceAggregate, eventstore.LoadAggregateOptions{
		SkipLoadEvents: true,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "LoadInvoiceAggregate")
	}

	return invoiceAggregate, nil
}

func NewInvoiceTempAggregateWithTenantAndID(tenant, id string) *InvoiceTempAggregate {
	invoiceAggregate := InvoiceTempAggregate{}
	invoiceAggregate.CommonTenantIdTempAggregate = aggregate.NewCommonTempAggregateWithTenantAndId(InvoiceAggregateType, tenant, id)
	invoiceAggregate.Tenant = tenant

	return &invoiceAggregate
}

func (a *InvoiceTempAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *invoicepb.NewInvoiceRequest:
		return nil, a.SimulateInvoice(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *InvoiceTempAggregate) SimulateInvoice(ctx context.Context, request *invoicepb.NewInvoiceRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceTempAggregate.SimulateInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createEvent, err := NewInvoiceNewEvent(a, sourceFields, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SimulateInvoice")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	return a.Apply(createEvent)
}
