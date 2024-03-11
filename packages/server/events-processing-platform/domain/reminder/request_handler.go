package reminder

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type ReminderRequestHandler interface {
	Handle(ctx context.Context, tenant, objectId string, request any, params ...map[string]any) (any, error)
	HandleWithRetry(ctx context.Context, tenant, objectId string, aggregateRequired bool, request any, params ...map[string]any) (any, error)
}

type reminderRequestHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewReminderRequestHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) ReminderRequestHandler {
	return &reminderRequestHandler{log: log, es: es, cfg: cfg}
}

func (h *reminderRequestHandler) Handle(ctx context.Context, tenant, objectId string, request any, params ...map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderRequestHandler.Handle")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)
	if params != nil && len(params) > 0 {
		span.LogFields(log.Object("params", params))
	}

	reminderAggregate, err := aggregate.LoadReminderAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	var requestParams map[string]any
	if params != nil && len(params) > 0 {
		requestParams = params[0]
	}
	result, err := reminderAggregate.HandleRequest(ctx, request, requestParams)
	if err != nil {
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
	if params != nil && len(params) > 0 {
		span.LogFields(log.Object("params", params))
	}

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		reminderAggregate, err := aggregate.LoadReminderAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptions())
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if aggregateRequired && eventstore.IsAggregateNotFound(reminderAggregate) {
			tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
			return nil, eventstore.ErrAggregateNotFound
		}

		var requestParams map[string]any
		if params != nil && len(params) > 0 {
			requestParams = params[0]
		}
		result, err := reminderAggregate.HandleRequest(ctx, request, requestParams)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = h.es.Save(ctx, reminderAggregate)
		if err == nil {
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
