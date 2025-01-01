package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	ServiceLineItemAggregateType eventstore.AggregateType = "service_line_item"
)

type ServiceLineItemAggregate struct {
	*eventstore.CommonTenantIdAggregate
	ServiceLineItem *model.ServiceLineItem
}

func NewServiceLineItemAggregateWithTenantAndID(tenant, id string) *ServiceLineItemAggregate {
	serviceLineItemAggregate := ServiceLineItemAggregate{}
	serviceLineItemAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(ServiceLineItemAggregateType, tenant, id)
	serviceLineItemAggregate.SetWhen(serviceLineItemAggregate.When)
	serviceLineItemAggregate.ServiceLineItem = &model.ServiceLineItem{}
	serviceLineItemAggregate.Tenant = tenant

	return &serviceLineItemAggregate
}

// GetServiceLineItemObjectID generates the object ID for a service line item.
func GetServiceLineItemObjectID(aggregateID string, tenant string) string {
	return eventstore.GetAggregateObjectID(aggregateID, tenant, ServiceLineItemAggregateType)
}

func (a *ServiceLineItemAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *servicelineitempb.CloseServiceLineItemGrpcRequest:
		return nil, a.CloseServiceLineItem(ctx, r, params)
	default:
		return nil, nil
	}
}

func (a *ServiceLineItemAggregate) CloseServiceLineItem(ctx context.Context, r *servicelineitempb.CloseServiceLineItemGrpcRequest, params map[string]any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.CloseServiceLineItem")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	// if future version - produce delete event
	if a.ServiceLineItem.StartedAt.After(utils.Now()) {
		deleteEvent, err := event.NewServiceLineItemCloseEvent(a)

		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewServiceLineItemDeleteEvent")
		}

		eventstore.EnrichEventWithMetadataExtended(&deleteEvent, span, eventstore.EventMetadata{
			Tenant: a.Tenant,
			UserId: r.GetLoggedInUserId(),
			App:    r.GetAppSource(),
		})
		return a.Apply(deleteEvent)
	}

	// Create the close event
	updatedAtNotNil := utils.Now()
	endedAtNotNil := utils.ToDate(utils.IfNotNilTimeWithDefault(r.EndedAt, utils.Now()))

	cancelled := false
	val, ok := params[model.PARAM_CANCELLED]
	if ok {
		cancelled = val.(bool)
	}

	closeEvent, err := event.NewServiceLineItemCloseEvent(a, endedAtNotNil, updatedAtNotNil, cancelled)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewServiceLineItemCloseEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&closeEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.GetLoggedInUserId(),
		App:    r.GetAppSource(),
	})

	return a.Apply(closeEvent)
}

func (a *ServiceLineItemAggregate) When(evt eventstore.Event) error {
	return nil
}
