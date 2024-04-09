package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

const (
	ServiceLineItemAggregateType eventstore.AggregateType = "service_line_item"
)

type ServiceLineItemAggregate struct {
	*aggregate.CommonTenantIdAggregate
	ServiceLineItem *model.ServiceLineItem
}

func NewServiceLineItemAggregateWithTenantAndID(tenant, id string) *ServiceLineItemAggregate {
	serviceLineItemAggregate := ServiceLineItemAggregate{}
	serviceLineItemAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(ServiceLineItemAggregateType, tenant, id)
	serviceLineItemAggregate.SetWhen(serviceLineItemAggregate.When)
	serviceLineItemAggregate.ServiceLineItem = &model.ServiceLineItem{}
	serviceLineItemAggregate.Tenant = tenant

	return &serviceLineItemAggregate
}

func (a *ServiceLineItemAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *servicelineitempb.CloseServiceLineItemGrpcRequest:
		return nil, a.CloseServiceLineItem(ctx, r, params)
	case *servicelineitempb.DeleteServiceLineItemGrpcRequest:
		return nil, a.DeleteServiceLineItem(ctx, r)
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
		deleteEvent, err := event.NewServiceLineItemDeleteEvent(a)

		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewServiceLineItemDeleteEvent")
		}

		aggregate.EnrichEventWithMetadataExtended(&deleteEvent, span, aggregate.EventMetadata{
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
	aggregate.EnrichEventWithMetadataExtended(&closeEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.GetLoggedInUserId(),
		App:    r.GetAppSource(),
	})

	return a.Apply(closeEvent)
}

func (a *ServiceLineItemAggregate) DeleteServiceLineItem(ctx context.Context, r *servicelineitempb.DeleteServiceLineItemGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.deleteServiceLineItem")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	deleteEvent, err := event.NewServiceLineItemDeleteEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewServiceLineItemDeleteEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&deleteEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.GetLoggedInUserId(),
		App:    r.GetAppSource(),
	})

	return a.Apply(deleteEvent)
}

func (a *ServiceLineItemAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ServiceLineItemCreateV1:
		return a.onCreate(evt)
	case event.ServiceLineItemUpdateV1:
		return a.onUpdate(evt)
	case event.ServiceLineItemDeleteV1:
		return a.onDelete()
	case event.ServiceLineItemCloseV1:
		return a.onClose(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *ServiceLineItemAggregate) onCreate(evt eventstore.Event) error {
	var eventData event.ServiceLineItemCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.ServiceLineItem.ID = a.ID
	a.ServiceLineItem.ContractId = eventData.ContractId
	a.ServiceLineItem.ParentId = eventData.ParentId
	a.ServiceLineItem.Billed = eventData.Billed
	a.ServiceLineItem.Quantity = eventData.Quantity
	a.ServiceLineItem.Price = eventData.Price
	a.ServiceLineItem.Name = eventData.Name
	a.ServiceLineItem.CreatedAt = eventData.CreatedAt
	a.ServiceLineItem.UpdatedAt = eventData.UpdatedAt
	a.ServiceLineItem.Source = eventData.Source
	a.ServiceLineItem.StartedAt = eventData.StartedAt
	a.ServiceLineItem.EndedAt = eventData.EndedAt
	a.ServiceLineItem.Comments = eventData.Comments
	a.ServiceLineItem.VatRate = eventData.VatRate

	return nil
}

// onServiceLineItemUpdate handles the update event for a service line item.
func (a *ServiceLineItemAggregate) onUpdate(evt eventstore.Event) error {
	var eventData event.ServiceLineItemUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	// Apply the changes from the event to the service line item model
	a.ServiceLineItem.Name = eventData.Name
	a.ServiceLineItem.Price = eventData.Price
	a.ServiceLineItem.Quantity = eventData.Quantity
	a.ServiceLineItem.Billed = eventData.Billed
	a.ServiceLineItem.UpdatedAt = eventData.UpdatedAt
	if constants.SourceOpenline == eventData.Source.Source {
		a.ServiceLineItem.Source.SourceOfTruth = eventData.Source.Source
	}
	a.ServiceLineItem.Comments = eventData.Comments
	a.ServiceLineItem.VatRate = eventData.VatRate
	if eventData.StartedAt != nil {
		a.ServiceLineItem.StartedAt = *eventData.StartedAt
	}

	return nil
}

func (a *ServiceLineItemAggregate) onDelete() error {
	a.ServiceLineItem.IsDeleted = true
	return nil
}

func (a *ServiceLineItemAggregate) onClose(evt eventstore.Event) error {
	var eventData event.ServiceLineItemCloseEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.ServiceLineItem.EndedAt = &eventData.EndedAt
	a.ServiceLineItem.UpdatedAt = eventData.UpdatedAt
	a.ServiceLineItem.IsCanceled = eventData.IsCanceled

	return nil
}
