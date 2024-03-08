package event_handler

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type UpdateReminderHandler interface {
	Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *reminderpb.UpdateReminderGrpcRequest) error
}

type updateReminderHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewUpdateReminderHandler(log logger.Logger, es eventstore.AggregateStore) UpdateReminderHandler {
	return &updateReminderHandler{log: log, es: es}
}

func (h *updateReminderHandler) Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *reminderpb.UpdateReminderGrpcRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateReminderHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	// Load or initialize the org plan aggregate
	reminderAggregate, err := aggregate.LoadReminderAggregate(ctx, h.es, baseRequest.Tenant, baseRequest.ObjectID, eventstore.LoadAggregateOptions{})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if eventstore.IsAggregateNotFound(reminderAggregate) {
		tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
		return eventstore.ErrAggregateNotFound
	}

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(request.UpdatedAt, utils.Now())

	updateEvent, err := events.NewReminderUpdateEvent(
		reminderAggregate,
		baseRequest.ObjectID,
		baseRequest.Tenant,
		request.Content,
		utils.TimestampProtoToTime(request.DueDate),
		request.Dismissed,
		updatedAtNotNil,
		extractReminderFieldsMask(request.FieldsMask),
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewReminderUpdateEvent")
	}
	commonAggregate.EnrichEventWithMetadataExtended(&updateEvent, span, commonAggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: baseRequest.LoggedInUserId,
	})
	if err := h.es.Save(ctx, reminderAggregate); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "es.Save")
	}

	return nil
}

func extractReminderFieldsMask(fields []reminderpb.ReminderFieldMask) []string {
	fieldsMask := make([]string, 0)
	if len(fields) == 0 {
		return fieldsMask
	}
	for _, field := range fields {
		switch field {
		case reminderpb.ReminderFieldMask_REMINDER_PROPERTY_CONTENT:
			fieldsMask = append(fieldsMask, events.FieldMaskContent)
		case reminderpb.ReminderFieldMask_REMINDER_PROPERTY_DISMISSED:
			fieldsMask = append(fieldsMask, events.FieldMaskDismissed)
		case reminderpb.ReminderFieldMask_REMINDER_PROPERTY_DUE_DATE:
			fieldsMask = append(fieldsMask, events.FieldMaskDueDate)
		}
	}
	return utils.RemoveDuplicates(fieldsMask)
}
