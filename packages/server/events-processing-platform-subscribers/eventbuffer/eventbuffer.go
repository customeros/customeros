package eventbuffer

import (
	"context"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/reminder"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
)

type EventBufferWatcher struct {
	ebr           repository.EventBufferRepository
	logger        logger.Logger
	es            eventstore.AggregateStore
	signalChannel chan os.Signal
	ticker        *time.Ticker
}

func NewEventBufferWatcher(ebr repository.EventBufferRepository, logger logger.Logger, es eventstore.AggregateStore) *EventBufferWatcher {
	return &EventBufferWatcher{ebr: ebr, logger: logger, es: es}
}

// Start starts the EventBufferWatcher
func (eb *EventBufferWatcher) Start(ctx context.Context) {
	eb.logger.Info("EventBufferWatcher started")

	eb.ticker = time.NewTicker(time.Second * 30)
	eb.signalChannel = make(chan os.Signal, 1)
	signal.Notify(eb.signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				// Run dispatch logic every n seconds
				eb.logger.Info("EventBufferWatcher.Dispatch: dispatch buffered events")
				err := eb.Dispatch(ctx)
				if err != nil {
					eb.logger.Errorf("EventBufferWatcher.Dispatch: error dispatching events: %s", err.Error())
				}
			case <-eb.signalChannel:
				// Shutdown goroutine
				eb.logger.Info("EventBufferWatcher.Dispatch: Got signal, exiting...")
				runtime.Goexit()
			}
		}
	}(eb.ticker)
}

// Stop stops the EventBufferWatcher
func (eb *EventBufferWatcher) Stop() {
	eb.signalChannel <- syscall.SIGTERM // TODO get the signal from the caller
	eb.ticker.Stop()
	eb.logger.Info("EventBufferWatcher stopped")
	close(eb.signalChannel)
	eb.signalChannel = nil
}

// Dispatch dispatches all expired events from event_buffer table, and delete them after dispatching
func (eb *EventBufferWatcher) Dispatch(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferWatcher.Dispatch")
	defer span.Finish()
	now := time.Now().UTC()
	eventBuffers, err := eb.ebr.GetByExpired(now)
	if err != nil {
		return err
	}
	if len(eventBuffers) == 0 {
		return nil
	}
	tracing.LogObjectAsJson(span, "expiredEvents", eventBuffers)
	for _, eventBuffer := range eventBuffers {
		if err := eb.HandleEvent(ctx, eventBuffer); err != nil {
			tracing.TraceErr(span, err)
			eb.logger.Errorf("EventBufferWatcher.Dispatch: error handling event: %s", err.Error())
			continue
		}
		err = eb.ebr.Delete(eventBuffer)
		if err != nil {
			return err
		}
	}
	return err
}

// HandleEvent loads the event aggregate and applies the event to it and pushes it into event store
func (eb *EventBufferWatcher) HandleEvent(ctx context.Context, eventBuffer entity.EventBuffer) error {
	evt := eventstore.Event{
		EventID:       eventBuffer.EventID,
		EventType:     eventBuffer.EventType,
		Data:          eventBuffer.EventData,
		Timestamp:     eventBuffer.EventTimestamp.UTC(),
		AggregateType: eventstore.AggregateType(eventBuffer.EventAggregateType),
		AggregateID:   eventBuffer.EventAggregateID,
		Version:       eventBuffer.EventVersion,
		Metadata:      eventBuffer.EventMetadata,
	}
	return eb.handleEvent(ctx, evt)
}

func (eb *EventBufferWatcher) handleEvent(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferWatcher.handleEvent")
	defer span.Finish()
	switch evt.EventType {
	case orgevents.OrganizationUpdateOwnerNotificationV1:
		//var data orgevents.OrganizationOwnerUpdateEvent
		//if err := json.Unmarshal(evt.Data, &data); err != nil {
		//	tracing.TraceErr(span, err)
		//	return err
		//}
		//organizationAggregate, err := orgaggregate.LoadOrganizationAggregate(ctx, eb.es, data.Tenant, data.OrganizationId, eventstore.LoadAggregateOptions{})
		//if err != nil {
		//	tracing.TraceErr(span, err)
		//	return err
		//}
		//err = organizationAggregate.Apply(evt)
		//if err != nil {
		//	tracing.TraceErr(span, err)
		//	return err
		//}
		//// Persist the changes to the event store
		//if err = eb.es.Save(ctx, organizationAggregate); err != nil {
		//	tracing.TraceErr(span, err)
		//	return err
		//}
		//return err
		return nil
	case reminder.ReminderNotificationV1:
		//var data reminder.ReminderNotificationEvent
		//if err := json.Unmarshal(evt.Data, &data); err != nil {
		//	tracing.TraceErr(span, err)
		//	return err
		//}
		//reminderAggregate, err := reminder.LoadReminderAggregate(ctx, eb.es, data.Tenant, data.ReminderId, eventstore.LoadAggregateOptions{})
		//if err != nil {
		//	tracing.TraceErr(span, err)
		//	return err
		//}
		//err = reminderAggregate.Apply(evt)
		//if err != nil {
		//	tracing.TraceErr(span, err)
		//	return err
		//}
		//// Persist the changes to the event store
		//if err = eb.es.Save(ctx, reminderAggregate); err != nil {
		//	tracing.TraceErr(span, err)
		//	return err
		//}
		//return err
		return nil
	default:
		return nil
	}
}
