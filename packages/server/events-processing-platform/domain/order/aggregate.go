package order

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	orderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/order"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

const (
	OrderAggregateType eventstore.AggregateType = "order"
)

type OrderAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Order *Order
}

func GetOrderObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, OrderAggregateType)
}

func LoadOrderAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string, options eventstore.LoadAggregateOptions) (*OrderAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadOrderAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	orderAggregate := NewOrderAggregateWithTenantAndID(tenant, objectID)

	err := eventstore.LoadAggregate(ctx, eventStore, orderAggregate, options)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return orderAggregate, nil
}

func NewOrderAggregateWithTenantAndID(tenant, id string) *OrderAggregate {
	orderAggregate := OrderAggregate{}
	orderAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(OrderAggregateType, tenant, id)
	orderAggregate.SetWhen(orderAggregate.When)
	orderAggregate.Order = &Order{}
	orderAggregate.Tenant = tenant

	return &orderAggregate
}

func (a *OrderAggregate) HandleRequest(ctx context.Context, request any, params ...map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *orderpb.UpsertOrderGrpcRequest:
		return nil, a.UpsertOrderRequest(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *OrderAggregate) UpsertOrderRequest(ctx context.Context, request *orderpb.UpsertOrderGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrderAggregate.UpsertOrderRequest")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())

	confirmedAtPtr := utils.TimestampProtoToTimePtr(request.ConfirmedAt)
	paidAtPtr := utils.TimestampProtoToTimePtr(request.PaidAt)
	fulfilledAtPtr := utils.TimestampProtoToTimePtr(request.FulfilledAt)
	canceledAtPtr := utils.TimestampProtoToTimePtr(request.CanceledAt)

	sourceFields := events.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	externalSystem := commonmodel.ExternalSystem{}
	externalSystem.FromGrpc(request.ExternalSystemFields)

	event, err := NewOrderUpsertEvent(a, sourceFields, externalSystem, request.OrganizationId, createdAtNotNil, updatedAtNotNil, confirmedAtPtr, paidAtPtr, fulfilledAtPtr, canceledAtPtr)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrderUpsertEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	return a.Apply(event)
}

func (a *OrderAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case OrderUpsertV1:
		return a.onOrderUpsertEvent(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *OrderAggregate) onOrderUpsertEvent(evt eventstore.Event) error {
	var eventData OrderUpsertEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Order.ID = a.ID
	a.Order.CreatedAt = eventData.CreatedAt
	a.Order.UpdatedAt = eventData.CreatedAt

	if eventData.ConfirmedAt != nil {
		a.Order.ConfirmedAt = *eventData.ConfirmedAt
	}
	if eventData.PaidAt != nil {
		a.Order.PaidAt = *eventData.PaidAt
	}
	if eventData.FulfilledAt != nil {
		a.Order.FulfilledAt = *eventData.FulfilledAt
	}
	if eventData.CanceledAt != nil {
		a.Order.CanceledAt = *eventData.CanceledAt
	}

	return nil
}
