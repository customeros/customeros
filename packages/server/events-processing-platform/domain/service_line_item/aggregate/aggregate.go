package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/service_line_item"
	events2 "github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
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
	case *servicelineitempb.DeleteServiceLineItemGrpcRequest:
		return nil, a.DeleteServiceLineItem(ctx, r)
	case *servicelineitempb.CreateServiceLineItemGrpcRequest:
		return nil, a.CreateServiceLineItem(ctx, r)
	case *servicelineitempb.UpdateServiceLineItemGrpcRequest:
		return nil, a.UpdateServiceLineItem(ctx, r)
	case *servicelineitempb.PauseServiceLineItemGrpcRequest:
		return nil, a.PauseServiceLineItem(ctx, r)
	case *servicelineitempb.ResumeServiceLineItemGrpcRequest:
		return nil, a.ResumeServiceLineItem(ctx, r)
	default:
		return nil, nil
	}
}

func (a *ServiceLineItemAggregate) CreateServiceLineItem(ctx context.Context, r *servicelineitempb.CreateServiceLineItemGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.createServiceLineItem")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	// fail if quantity is negative
	if r.Quantity < 0 {
		err := errors.New(events2.FieldValidation + ": quantity must not be negative")
		tracing.TraceErr(span, err)
		return err
	}

	// Adjust vat rate
	if r.VatRate < 0 {
		r.VatRate = 0
	}
	r.VatRate = utils.TruncateFloat64(r.VatRate, 2)

	sourceFields := common.Source{}
	sourceFields.FromGrpc(r.SourceFields)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.CreatedAt), utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.UpdatedAt), createdAtNotNil)
	startedAtNotNil := utils.ToDate(utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.StartedAt), utils.Now()))
	endedAtNillable := utils.ToDatePtr(utils.TimestampProtoToTimePtr(r.EndedAt))

	if endedAtNillable != nil && endedAtNillable.Before(startedAtNotNil) {
		err := errors.New(events2.FieldValidation + ": endedAt must be after startedAt")
		tracing.TraceErr(span, err)
		return err
	}

	dataFields := model.ServiceLineItemDataFields{
		Billed:     model.BilledType(r.Billed),
		Quantity:   r.Quantity,
		Price:      r.Price,
		Name:       r.Name,
		ContractId: r.ContractId,
		ParentId:   utils.StringFirstNonEmpty(r.ParentId, GetServiceLineItemObjectID(a.GetID(), a.GetTenant())),
		VatRate:    r.VatRate,
		Comments:   r.Comments,
	}

	createEvent, err := event.NewServiceLineItemCreateEvent(
		a,
		dataFields,
		sourceFields,
		createdAtNotNil,
		updatedAtNotNil,
		startedAtNotNil,
		endedAtNillable,
		"", // alexbalexb TODO: previousVersionId pass it from service
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewServiceLineItemCreateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&createEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *ServiceLineItemAggregate) UpdateServiceLineItem(ctx context.Context, r *servicelineitempb.UpdateServiceLineItemGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.UpdateServiceLineItem")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	// Do not allow updates on deleted or canceled service line items
	if a.ServiceLineItem.IsDeleted {
		err := errors.New(events2.Validate + ": cannot update a deleted service line item")
		tracing.TraceErr(span, err)
		return err
	}
	if a.ServiceLineItem.IsCanceled {
		err := errors.New(events2.Validate + ": cannot update a canceled service line item")
		tracing.TraceErr(span, err)
		return err
	}

	// fail if quantity is negative
	if r.Quantity < 0 {
		err := errors.New(events2.FieldValidation + ": quantity must not be negative")
		tracing.TraceErr(span, err)
		return err
	}

	billedType := model.BilledType(r.Billed)
	// do not allow changing billed type
	if a.ServiceLineItem.Billed != billedType.String() && a.ServiceLineItem.Billed != "" {
		err := errors.New(events2.Validate + ": cannot change billed type")
		tracing.TraceErr(span, err)
		return err
	}

	// Adjust vat rate
	if r.VatRate < 0 {
		r.VatRate = 0
	}
	r.VatRate = utils.TruncateFloat64(r.VatRate, 2)

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(r.UpdatedAt), utils.Now())

	source := common.Source{}
	source.FromGrpc(r.SourceFields)

	dataFields := model.ServiceLineItemDataFields{
		Billed:   billedType,
		Quantity: r.Quantity,
		Price:    r.Price,
		Name:     r.Name,
		Comments: r.Comments,
		VatRate:  r.VatRate,
	}

	// Prepare the data for the update event
	updateEvent, err := event.NewServiceLineItemUpdateEvent(
		a,
		dataFields,
		source,
		updatedAtNotNil,
		utils.TimestampProtoToTimePtr(r.StartedAt),
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewServiceLineItemUpdateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.LoggedInUserId,
		App:    source.AppSource,
	})

	return a.Apply(updateEvent)
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
	eventstore.EnrichEventWithMetadataExtended(&deleteEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.GetLoggedInUserId(),
		App:    r.GetAppSource(),
	})

	return a.Apply(deleteEvent)
}

func (a *ServiceLineItemAggregate) PauseServiceLineItem(ctx context.Context, r *servicelineitempb.PauseServiceLineItemGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.PauseServiceLineItem")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	pauseEvent, err := event.NewServiceLineItemPauseEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewServiceLineItemPauseEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&pauseEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.GetLoggedInUserId(),
		App:    r.GetAppSource(),
	})

	return a.Apply(pauseEvent)
}

func (a *ServiceLineItemAggregate) ResumeServiceLineItem(ctx context.Context, r *servicelineitempb.ResumeServiceLineItemGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ServiceLineItemAggregate.ResumeServiceLineItem")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", r)

	resumeEvent, err := event.NewServiceLineItemResumeEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewServiceLineItemResumeEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&resumeEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: r.GetLoggedInUserId(),
		App:    r.GetAppSource(),
	})

	return a.Apply(resumeEvent)
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
	case event.ServiceLineItemPauseV1,
		event.ServiceLineItemResumeV1:
		return nil
	default:
		if strings.HasPrefix(evt.GetEventType(), events2.EsInternalStreamPrefix) {
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
	if events2.SourceOpenline == eventData.Source.Source {
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
