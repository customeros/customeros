package events

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	event "github.com/openline-ai/openline-customer-os/packages/server/events/events/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

const RetriesOnOptimisticLockException = 5

type RequestHandler interface {
	Store(ctx context.Context, request events.BaseEvent, aggregateOptions eventstore.LoadAggregateOptions) (any, error)
}

type eventStoreService struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewEventStoreService(log logger.Logger, es eventstore.AggregateStore) *eventStoreService {
	return &eventStoreService{log: log, es: es}
}

func (h *eventStoreService) Store(ctx context.Context, request events.BaseEvent, aggregateOptions eventstore.LoadAggregateOptions) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventStoreService.Store")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "request", request)
	span.LogFields(log.Object("aggregateOptions", aggregateOptions))

	for attempt := 0; attempt == 0 || attempt < RetriesOnOptimisticLockException; attempt++ {
		agg := InitAggregate(request)

		if agg == nil {
			err := errors.New("aggregate not initialized")
			tracing.TraceErr(span, err)
			return nil, err
		}

		err := LoadAggregate(ctx, h.es, agg, aggregateOptions)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if aggregateOptions.Required && eventstore.IsAggregateNotFound(agg) {
			tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
			return nil, eventstore.ErrAggregateNotFound
		}

		//todo validate
		//if err := validator.GetValidator().Struct(eventData); err != nil {
		//	return eventstore.Event{}, errors.Wrap(err, "failed to validate UserUpdateEvent")
		//}

		evt := eventstore.NewBaseEvent(agg, request.EventName)

		if err := evt.SetJsonData(&request); err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = agg.Apply(evt)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		err = h.es.Save(ctx, agg)
		if err == nil {
			return agg, nil
		}

		if eventstore.IsEventStoreErrorCodeWrongExpectedVersion(err) {
			// Handle concurrency error
			if attempt == RetriesOnOptimisticLockException-1 {
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

func InitAggregate(request events.BaseEvent) eventstore.Aggregate {
	switch request.EntityType {
	case CONTACT:
		return event.NewContactAggregateWithTenantAndID(request.Tenant, request.EntityId)
	}

	return nil
}

func LoadAggregate(ctx context.Context, eventStore eventstore.AggregateStore, agg eventstore.Aggregate, options eventstore.LoadAggregateOptions) error {
	err := eventStore.Exists(ctx, agg.GetID())
	if err != nil {
		if !errors.Is(err, eventstore.ErrAggregateNotFound) {
			return err
		} else {
			return nil
		}
	}

	if options.SkipLoadEvents {
		if err = eventStore.LoadVersion(ctx, agg); err != nil {
			return err
		}
	} else {
		if err = eventStore.Load(ctx, agg); err != nil {
			return err
		}
	}

	return nil
}
