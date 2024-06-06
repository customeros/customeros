package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type ContractTempAggregate struct {
	*aggregate.CommonTenantIdTempAggregate
}

func NewContractTempAggregateWithTenantAndID(tenant, id string) *ContractTempAggregate {
	contractTempAggregate := ContractTempAggregate{}
	contractTempAggregate.CommonTenantIdTempAggregate = aggregate.NewCommonTempAggregateWithTenantAndId(ContractAggregateType, tenant, id)
	contractTempAggregate.Tenant = tenant

	return &contractTempAggregate
}

func (a *ContractTempAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractTempAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *contractpb.RefreshContractStatusGrpcRequest:
		return nil, a.refreshContractStatus(ctx, r)
	case *contractpb.RefreshContractLtvGrpcRequest:
		return nil, a.refreshContractLtv(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContractTempAggregate) refreshContractStatus(ctx context.Context, request *contractpb.RefreshContractStatusGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.refreshContractStatus")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	updateEvent, err := event.NewContractRefreshStatusEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractRefreshStatusEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(updateEvent)
}

func (a *ContractTempAggregate) refreshContractLtv(ctx context.Context, r *contractpb.RefreshContractLtvGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.refreshContractLtv")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	updateEvent, err := event.NewContractRefreshLtvEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractRefreshLtvEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.GetAppSource(),
	})

	return a.Apply(updateEvent)
}
