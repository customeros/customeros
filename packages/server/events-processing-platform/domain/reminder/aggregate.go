package reminder

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

const ReminderAggregateType = "reminder"

type ReminderAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Reminder *Reminder
}

func (a *ReminderAggregate) HandleRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *reminderpb.CreateReminderGrpcRequest:
		return nil, a.CreateReminder(ctx, r)
	case *reminderpb.UpdateReminderGrpcRequest:
		return nil, a.UpdateReminder(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func NewReminderAggregateWithTenantAndID(tenant, id string) *ReminderAggregate {
	reminderAggregate := ReminderAggregate{}
	reminderAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(ReminderAggregateType, tenant, id)
	reminderAggregate.SetWhen(reminderAggregate.When)
	reminderAggregate.Reminder = &Reminder{}
	reminderAggregate.Tenant = tenant

	return &reminderAggregate
}

func (a *ReminderAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case ReminderCreateV1:
		return a.whenReminderCreate(event)
	case ReminderUpdateV1:
		return a.whenReminderUpdate(event)
	case ReminderNotificationV1:
		return nil
	default:
		if strings.HasPrefix(event.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *ReminderAggregate) whenReminderCreate(evt eventstore.Event) error {
	var reminderCreateEvent ReminderCreateEvent

	if err := evt.GetJsonData(&reminderCreateEvent); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Reminder = &Reminder{
		Content:        reminderCreateEvent.Content,
		DueDate:        reminderCreateEvent.DueDate,
		Dismissed:      reminderCreateEvent.Dismissed,
		CreatedAt:      reminderCreateEvent.CreatedAt,
		UserID:         reminderCreateEvent.UserId,
		OrganizationID: reminderCreateEvent.OrganizationId,
	}
	return nil
}

func (a *ReminderAggregate) whenReminderUpdate(evt eventstore.Event) error {
	var reminderUpdateEvent ReminderUpdateEvent

	if err := evt.GetJsonData(&reminderUpdateEvent); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if reminderUpdateEvent.UpdateContent() {
		a.Reminder.Content = reminderUpdateEvent.Content
	}
	if reminderUpdateEvent.UpdateDueDate() {
		a.Reminder.DueDate = reminderUpdateEvent.DueDate
	}
	if reminderUpdateEvent.UpdateDismissed() {
		a.Reminder.Dismissed = reminderUpdateEvent.Dismissed
	}
	return nil
}

func (a *ReminderAggregate) CreateReminder(ctx context.Context, request *reminderpb.CreateReminderGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ReminderAggregate.CreateReminder")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	dueDateNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.DueDate), utils.Now())
	sourceFields := events.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createEvent, err := NewReminderCreateEvent(
		a,
		request.Content,
		request.LoggedInUserId,
		request.OrganizationId,
		request.Dismissed,
		createdAtNotNil,
		dueDateNotNil,
		sourceFields,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewReminderCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *ReminderAggregate) UpdateReminder(ctx context.Context, request *reminderpb.UpdateReminderGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ReminderAggregate.UpdateReminder")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	dueDateNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.DueDate), utils.Now())

	updateEvent, err := NewReminderUpdateEvent(
		a,
		request.Content,
		dueDateNotNil,
		request.Dismissed,
		request.UpdatedAt.AsTime(),
		extractReminderFieldsMask(request.FieldsMask),
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewReminderUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})
	return a.Apply(updateEvent)
}

func extractReminderFieldsMask(fields []reminderpb.ReminderFieldMask) []string {
	fieldsMask := make([]string, 0)
	if len(fields) == 0 {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case reminderpb.ReminderFieldMask_REMINDER_PROPERTY_CONTENT:
			fieldsMask = append(fieldsMask, FieldMaskContent)
		case reminderpb.ReminderFieldMask_REMINDER_PROPERTY_DISMISSED:
			fieldsMask = append(fieldsMask, FieldMaskDismissed)
		case reminderpb.ReminderFieldMask_REMINDER_PROPERTY_DUE_DATE:
			fieldsMask = append(fieldsMask, FieldMaskDueDate)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}
