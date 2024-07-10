package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type PhoneNumberTempAggregate struct {
	*aggregate.CommonTenantIdTempAggregate
}

func NewPhoneNumberTempAggregateWithTenantAndID(tenant, id string) *PhoneNumberTempAggregate {
	phoneNumberTempAggregate := PhoneNumberTempAggregate{}
	phoneNumberTempAggregate.CommonTenantIdTempAggregate = aggregate.NewCommonTempAggregateWithTenantAndId(PhoneNumberAggregateType, tenant, id)
	phoneNumberTempAggregate.Tenant = tenant

	return &phoneNumberTempAggregate
}

func (a *PhoneNumberTempAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberTempAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *phonenumberpb.RequestPhoneNumberValidationGrpcRequest:
		return nil, a.requestPhoneNumberValidation(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *PhoneNumberTempAggregate) requestPhoneNumberValidation(ctx context.Context, request *phonenumberpb.RequestPhoneNumberValidationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberTempAggregate.requestPhoneNumberValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	updateEvent, err := events.NewPhoneNumberValidateEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "failed to create PhoneNumberValidateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(updateEvent)
}
