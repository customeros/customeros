package contact

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type ContactTempAggregate struct {
	*eventstore.CommonTenantIdTempAggregate
}

func NewContactTempAggregateWithTenantAndID(tenant, id string) *ContactTempAggregate {
	contactTempAggregate := ContactTempAggregate{}
	contactTempAggregate.CommonTenantIdTempAggregate = eventstore.NewCommonTempAggregateWithTenantAndId(ContactAggregateType, tenant, id)
	contactTempAggregate.Tenant = tenant

	return &contactTempAggregate
}

func (a *ContactTempAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactTempAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *contactpb.EnrichContactGrpcRequest:
		return nil, a.requestEnrichContact(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContactTempAggregate) requestEnrichContact(ctx context.Context, request *contactpb.EnrichContactGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactTempAggregate.requestEnrichContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	enrichEvent, err := event.NewContactRequestEnrich(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactRequestEnrich")
	}
	eventstore.EnrichEventWithMetadataExtended(&enrichEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(enrichEvent)
}
