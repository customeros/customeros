package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type RequestHandler interface {
	HandleGRPCRequest(ctx context.Context, initAggregate func() eventstore.Aggregate, aggregateOptions eventstore.LoadAggregateOptions, request any, params ...map[string]any) (any, error)
}

type requestHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewRequestHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) *requestHandler {
	return &requestHandler{log: log, es: es, cfg: cfg}
}

func (h *requestHandler) HandleGRPCRequest(ctx context.Context, initAggregate func() eventstore.Aggregate, aggregateOptions eventstore.LoadAggregateOptions, request any, params ...map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RequestHandler.HandleGRPCRequest")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "request", request)
	span.LogFields(log.Object("aggregateOptions", aggregateOptions))

	if params != nil && len(params) > 0 {
		span.LogFields(log.Object("params", params))
	}

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		agg := initAggregate()
		err := aggregate.LoadAggregate(ctx, h.es, agg, aggregateOptions)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if aggregateOptions.Required && eventstore.IsAggregateNotFound(agg) {
			tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
			return nil, eventstore.ErrAggregateNotFound
		}

		var requestParams map[string]any
		if params != nil && len(params) > 0 {
			requestParams = params[0]
		}
		result, err := agg.HandleGRPCRequest(ctx, request, requestParams)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = h.es.Save(ctx, agg)
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
