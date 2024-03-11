package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const ReminderAggregateType = "reminder"

type ReminderAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Reminder *model.Reminder
}

func (a *ReminderAggregate) HandleRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *reminderpb.CreateReminderGrpcRequest:
		return nil, a.CreateReminder(ctx, r)
	case *reminderpb.UpdateReminderGrpcRequest:
		//return nil, a.UpdateReminder(ctx, r) // TODO implement update accordingly
		return nil, nil
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func NewReminderAggregateWithTenantAndID(tenant, id string) *ReminderAggregate {
	reminderAggregate := ReminderAggregate{}
	reminderAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(ReminderAggregateType, tenant, id)
	reminderAggregate.SetWhen(reminderAggregate.When)
	reminderAggregate.Reminder = &model.Reminder{}
	reminderAggregate.Tenant = tenant

	return &reminderAggregate
}

func (a *ReminderAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case events.ReminderCreateV1:
		return a.whenReminderCreate(event)
	case events.ReminderUpdateV1:
		return a.whenReminderUpdate(event)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *ReminderAggregate) whenReminderCreate(evt eventstore.Event) error {
	var reminderCreateEvent events.ReminderCreateEvent

	if err := evt.GetJsonData(&reminderCreateEvent); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Reminder = &model.Reminder{
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
	var reminderUpdateEvent events.ReminderUpdateEvent

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderAggregate.CreateReminder")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTime(request.CreatedAt), utils.Now())
	dueDateNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTime(request.DueDate), utils.Now())
	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createEvent, err := events.NewReminderCreateEvent(
		a,
		request.Content,
		request.UserId,
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
		UserId: request.UserId, // TODO - should be baseRequest.LoggedInUserId, add logged in user id to proto
		App:    request.SourceFields.AppSource,
	})

	return a.Apply(createEvent)
}
