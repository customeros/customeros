package contract

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type ContractRequestHandler interface {
	Handle(ctx context.Context, tenant, objectId string, request any) (any, error)
	HandleWithRetry(ctx context.Context, tenant, objectId string, aggregateRequired bool, request any) (any, error)
	HandleTemp(ctx context.Context, tenant, objectId string, request any) (any, error)
}

type contractRequestHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewContractRequestHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) ContractRequestHandler {
	return &contractRequestHandler{log: log, es: es, cfg: cfg}
}

func (h *contractRequestHandler) Handle(ctx context.Context, tenant, objectId string, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRequestHandler.Handle")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)

	contractAggregate, err := aggregate.LoadContractAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptions())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	result, err := contractAggregate.HandleRequest(ctx, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, h.es.Save(ctx, contractAggregate)
}

func (h *contractRequestHandler) HandleWithRetry(ctx context.Context, tenant, objectId string, aggregateRequired bool, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRequestHandler.HandleWithRetry")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)
	span.LogFields(log.Bool("aggregateRequired", aggregateRequired))

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		contractAggregate, err := aggregate.LoadContractAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptions())
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if aggregateRequired && eventstore.IsAggregateNotFound(contractAggregate) {
			tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
			return nil, eventstore.ErrAggregateNotFound
		}

		result, err := contractAggregate.HandleRequest(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = h.es.Save(ctx, contractAggregate)
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

func (h *contractRequestHandler) HandleTemp(ctx context.Context, tenant, objectId string, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRequestHandler.Handle")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)

	contractTempAggregate, err := aggregate.LoadContractTempAggregate(ctx, h.es, tenant, objectId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	result, err := contractTempAggregate.HandleRequest(ctx, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, h.es.Save(ctx, contractTempAggregate)
}
