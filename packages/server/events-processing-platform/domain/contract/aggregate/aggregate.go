package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	ContractAggregateType eventstore.AggregateType = "contract"
)

func GetContractObjectID(aggregateID, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, ContractAggregateType)
}

type ContractAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Contract *model.Contract
}

func NewContractAggregateWithTenantAndID(tenant, id string) *ContractAggregate {
	contractAggregate := ContractAggregate{}
	contractAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(ContractAggregateType, tenant, id)
	contractAggregate.SetWhen(contractAggregate.When)
	contractAggregate.Contract = &model.Contract{}
	contractAggregate.Tenant = tenant

	return &contractAggregate
}

func (a *ContractAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *contractpb.SoftDeleteContractGrpcRequest:
		return nil, a.softDeleteContract(ctx, r)
	case *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest:
		return nil, a.rolloutRenewalOpportunityOnExpiration(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContractAggregate) rolloutRenewalOpportunityOnExpiration(ctx context.Context, request *contractpb.RolloutRenewalOpportunityOnExpirationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.rolloutRenewalOpportunityOnExpiration")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	updateEvent, err := event.NewRolloutRenewalOpportunityEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewRolloutRenewalOpportunityEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.GetAppSource(),
	})

	return a.Apply(updateEvent)
}

func isUpdated(field string, fieldsMask []string) bool {
	return len(fieldsMask) == 0 || utils.Contains(fieldsMask, field)
}

func (a *ContractAggregate) softDeleteContract(ctx context.Context, r *contractpb.SoftDeleteContractGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContractAggregate.softDeleteContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	deleteEvent, err := event.NewContractDeleteEvent(a, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContractDeleteEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&deleteEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.LoggedInUserId,
		App:    r.AppSource,
	})

	return a.Apply(deleteEvent)
}

func (a *ContractAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ContractUpdateStatusV1:
		return a.onContractRefreshStatus(evt)
	case event.ContractRolloutRenewalOpportunityV1:
		return nil
	case event.ContractDeleteV1:
		return a.onContractDelete(evt)
	default:
		return nil
	}
}

func (a *ContractAggregate) onContractRefreshStatus(evt eventstore.Event) error {
	var eventData event.ContractUpdateStatusEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Contract.Status = eventData.Status
	return nil
}

func (a *ContractAggregate) onContractDelete(evt eventstore.Event) error {
	var eventData event.ContractDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Contract.Removed = true
	return nil
}
