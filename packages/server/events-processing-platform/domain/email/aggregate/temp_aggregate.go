package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type EmailTempAggregate struct {
	*aggregate.CommonTenantIdTempAggregate
}

func NewEmailTempAggregateWithTenantAndID(tenant, id string) *EmailTempAggregate {
	emailTempAggregate := EmailTempAggregate{}
	emailTempAggregate.CommonTenantIdTempAggregate = aggregate.NewCommonTempAggregateWithTenantAndId(EmailAggregateType, tenant, id)
	emailTempAggregate.Tenant = tenant

	return &emailTempAggregate
}

func (a *EmailTempAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailTempAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *emailpb.RequestEmailValidationGrpcRequest:
		return nil, a.requestEmailValidation(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *EmailTempAggregate) requestEmailValidation(ctx context.Context, request *emailpb.RequestEmailValidationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailTempAggregate.requestEmailValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	updateEvent, err := events.NewEmailValidateEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "failed to create EmailValidateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(updateEvent)
}
