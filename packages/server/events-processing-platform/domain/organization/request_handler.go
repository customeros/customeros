package organization

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type OrganizationRequestHandler interface {
	HandleWithRetry(ctx context.Context, tenant, objectId string, request any) (any, error)
	HandleTemp(ctx context.Context, tenant, objectId string, request any) (any, error)
	HandleTempWithRetry(ctx context.Context, tenant, objectId string, request any) (any, error)
}

type organizationRequestHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
	cfg config.Utils
}

func NewOrganizationRequestHandler(log logger.Logger, es eventstore.AggregateStore, cfg config.Utils) OrganizationRequestHandler {
	return &organizationRequestHandler{log: log, es: es, cfg: cfg}
}

func (h *organizationRequestHandler) HandleWithRetry(ctx context.Context, tenant, objectId string, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRequestHandler.HandleWithRetry")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		organizationAggregate, err := aggregate.LoadOrganizationAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptionsWithRequired())
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		result, err := organizationAggregate.HandleRequest(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = h.es.Save(ctx, organizationAggregate)
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

func (h *organizationRequestHandler) HandleTemp(ctx context.Context, tenant, objectId string, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRequestHandler.Handle")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)

	orgTempAggregate, err := aggregate.LoadOrganizationTempAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptionsWithRequired().WithSkipLoadEvents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	result, err := orgTempAggregate.HandleRequest(ctx, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, h.es.Save(ctx, orgTempAggregate)
}

func (h *organizationRequestHandler) HandleTempWithRetry(ctx context.Context, tenant, objectId string, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRequestHandler.HandleTempWithRetry")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "request", request)

	for attempt := 0; attempt == 0 || attempt < h.cfg.RetriesOnOptimisticLockException; attempt++ {
		orgTempAggregate, err := aggregate.LoadOrganizationTempAggregate(ctx, h.es, tenant, objectId, *eventstore.NewLoadAggregateOptions().WithSkipLoadEvents())
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		result, err := orgTempAggregate.HandleRequest(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = h.es.Save(ctx, orgTempAggregate)
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
