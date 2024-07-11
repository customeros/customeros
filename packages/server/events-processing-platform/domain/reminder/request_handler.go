package reminder

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ReminderRequestHandler interface {
	Handle(ctx context.Context, tenant, objectId string, request any, params ...map[string]any) (any, error)
	HandleWithRetry(ctx context.Context, tenant, objectId string, aggregateRequired bool, request any, params ...map[string]any) (any, error)
}

type reminderRequestHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
	ebs *eventstore.EventBufferService
}

func NewReminderRequestHandler(log logger.Logger, es eventstore.AggregateStore, ebs *eventstore.EventBufferService, cfg config.Utils) ReminderRequestHandler {
	return &reminderRequestHandler{log: log, es: es, ebs: ebs, cfg: cfg}
}

func (h *reminderRequestHandler) Handle(ctx context.Context, tenant, objectId string, request any, params ...map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderRequestHandler.Handle")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)
	if len(params) > 0 {
		span.LogFields(log.Object("params", params))
	}

	reminderAggregate, err := LoadReminderAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	var requestParams map[string]any
	if len(params) > 0 {
		requestParams = params[0]
	}
	result, err := reminderAggregate.HandleRequest(ctx, request, requestParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	req := request.(*reminderpb.CreateReminderGrpcRequest)

	if err := h.createAndParkReminderNotification(ctx, reminderAggregate, req.LoggedInUserId, req.OrganizationId, reminderAggregate.GetID(), req.SourceFields.AppSource, req.Content, req.DueDate.AsTime().UTC()); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, h.es.Save(ctx, reminderAggregate)
}

func (h *reminderRequestHandler) HandleWithRetry(ctx context.Context, tenant, objectId string, aggregateRequired bool, request any, params ...map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderRequestHandler.HandleWithRetry")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)
	span.LogFields(log.Bool("aggregateRequired", aggregateRequired))
	if len(params) > 0 {
		span.LogFields(log.Object("params", params))
	}

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		reminderAggregate, err := LoadReminderAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptions())
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if aggregateRequired && eventstore.IsAggregateNotFound(reminderAggregate) {
			tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
			return nil, eventstore.ErrAggregateNotFound
		}

		var requestParams map[string]any
		if len(params) > 0 {
			requestParams = params[0]
		}
		result, err := reminderAggregate.HandleRequest(ctx, request, requestParams)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = h.es.Save(ctx, reminderAggregate)
		if err == nil {

			if reminderAggregate.Reminder.Dismissed {
				err := h.deleteParkedReminderNotification(ctx, reminderAggregate.GetID())
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			} else {
				if err := h.updateParkedReminderNotification(ctx, reminderAggregate, reminderAggregate.GetID(), request); err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			}

			return result, nil // Save successful
		}

		if eventstore.IsEventStoreErrorCodeWrongExpectedVersion(err) {
			// Handle concurrency error
			if attempt == h.cfg.RetriesOnOptimisticLockException-1 {
				// If we have reached the maximum number of retries, return an error
				tracing.TraceErr(span, err)
				return nil, err
			}
			span.LogFields(log.Int("retryAttempt", attempt+1))
			time.Sleep(utils.BackOffExponentialDelay(attempt)) // backoffDelay is a function that increases the delay with each attempt
			continue                                           // Retry
		} else {
			// Some other error occurred
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	err := errors.New("reached maximum number of retries")
	tracing.TraceErr(span, err)
	return nil, err
}

func (h *reminderRequestHandler) createAndParkReminderNotification(ctx context.Context, agg *ReminderAggregate, userId, organizationId, reminderId, appSource, content string, createdAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderRequestHandler.createAndParkReminderNotification")
	defer span.Finish()

	if agg.Reminder.DueDate.Before(utils.Now()) {
		return nil
	}

	// create notification event and Park it in the event buffer for notifications
	event, err := createNotificationEvent(ctx, agg, userId, organizationId, appSource, content, createdAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	dueDate := agg.Reminder.DueDate.UTC() // when buffer should dispatch reminder notification event

	err = h.ebs.Park(*event, agg.GetTenant(), reminderId, dueDate)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (h *reminderRequestHandler) updateParkedReminderNotification(ctx context.Context, agg *ReminderAggregate, reminderId string, request any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderRequestHandler.updateParkedReminderNotification")
	defer span.Finish()

	req, ok := request.(*reminderpb.UpdateReminderGrpcRequest)
	if !ok {
		return nil
	}

	if agg.Reminder.DueDate.Before(utils.Now()) {
		return nil
	}

	parkedReminder, err := h.ebs.GetById(reminderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	//if reminder is missing, we recreate it
	if parkedReminder == nil {
		err := h.createAndParkReminderNotification(ctx, agg, req.LoggedInUserId, agg.Reminder.OrganizationID, reminderId, req.AppSource, req.Content, agg.Reminder.CreatedAt)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil
	}

	var parkedReminderEventData ReminderNotificationEvent
	if err := json.Unmarshal(parkedReminder.EventData, &parkedReminderEventData); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	parkedReminderEventData.Content = req.Content

	parkedReminder.ExpiryTimestamp = agg.Reminder.DueDate.UTC()
	parkedReminder.EventData, err = json.Marshal(parkedReminderEventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = h.ebs.Update(parkedReminder)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (h *reminderRequestHandler) deleteParkedReminderNotification(ctx context.Context, reminderId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderRequestHandler.deleteParkedReminderNotification")
	defer span.Finish()

	parkedReminder, err := h.ebs.GetById(reminderId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if parkedReminder == nil {
		return nil
	}

	err = h.ebs.Delete(parkedReminder)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func createNotificationEvent(ctx context.Context, agg *ReminderAggregate, loggedInUserId, organizationId, appSource, content string, reminderCreatedAt time.Time) (*eventstore.Event, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createNotificationEvent")
	defer span.Finish()
	tenant := agg.GetTenant()
	tracing.SetCommandHandlerSpanTags(ctx, span, tenant, loggedInUserId)

	event, err := NewReminderNotificationEvent(
		agg,
		loggedInUserId,
		organizationId,
		content,
		reminderCreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "NewOrganizationOwnerUpdateNotificationEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: tenant,
		UserId: loggedInUserId,
		App:    appSource,
	})

	return &event, nil
}
