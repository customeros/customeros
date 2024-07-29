package email

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type EmailTempAggregate struct {
	*eventstore.CommonTenantIdTempAggregate
}

func NewEmailTempAggregateWithTenantAndID(tenant, id string) *EmailTempAggregate {
	emailTempAggregate := EmailTempAggregate{}
	emailTempAggregate.CommonTenantIdTempAggregate = eventstore.NewCommonTempAggregateWithTenantAndId(EmailAggregateType, tenant, id)
	emailTempAggregate.Tenant = tenant

	return &emailTempAggregate
}

func (a *EmailTempAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailTempAggregate.HandleGRPCRequest")
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

	updateEvent, err := event.NewEmailValidateEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "failed to create EmailValidateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(updateEvent)
}
