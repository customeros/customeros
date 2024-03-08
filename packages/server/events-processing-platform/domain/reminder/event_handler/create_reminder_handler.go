package event_handler

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type CreateReminderHandler interface {
	Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *reminderpb.CreateReminderGrpcRequest) error
}

type createReminderHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCreateReminderHandler(log logger.Logger, es eventstore.AggregateStore) CreateReminderHandler {
	return &createReminderHandler{log: log, es: es}
}

func (h *createReminderHandler) Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *reminderpb.CreateReminderGrpcRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createReminderHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	// Load or initialize the org plan aggregate
	reminderAggregate, err := aggregate.LoadReminderAggregate(ctx, h.es, baseRequest.Tenant, baseRequest.ObjectID, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if reminderAggregate == nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewReminderCreateEvent:AGGREGATE_IS_NIL")
	}

	created := utils.TimestampProtoToTime(request.CreatedAt)
	createdAtNotNil := utils.IfNotNilTimeWithDefault(created, utils.Now())
	due := utils.TimestampProtoToTime(request.DueDate)
	dueNotNil := utils.IfNotNilTimeWithDefault(due, utils.Now())

	createEvent, err := events.NewReminderCreateEvent(
		reminderAggregate,
		baseRequest.Tenant,
		request.Content,
		baseRequest.ObjectID,
		request.UserId,
		request.OrganizationId,
		request.Dismissed,
		createdAtNotNil,
		dueNotNil,
		cmnmod.SourceFromGrpc(request.SourceFields),
	)

	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewReminderCreateEvent")
	}
	commonAggregate.EnrichEventWithMetadataExtended(&createEvent, span, commonAggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.UserId,
	})

	if err := h.es.Save(ctx, reminderAggregate); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "es.Save")
	}

	return nil
}
