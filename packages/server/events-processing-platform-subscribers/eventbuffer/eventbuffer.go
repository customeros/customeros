package eventbuffer

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	orgaggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	orgevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	reminder "github.com/openline-ai/openline-customer-os/packages/server/events/event/reminder/event"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
)

// Deprecated: use events.EventBufferWatcher instead
type EventBufferWatcher struct {
	ebr           postgresRepository.EventBufferRepository
	logger        logger.Logger
	es            eventstore.AggregateStore
	signalChannel chan os.Signal
	ticker        *time.Ticker
}

// Deprecated: use events.EventBufferWatcher instead
func NewEventBufferWatcher(ebr postgresRepository.EventBufferRepository, logger logger.Logger, es eventstore.AggregateStore) *EventBufferWatcher {
	return &EventBufferWatcher{ebr: ebr, logger: logger, es: es}
}

// Deprecated: use events.EventBufferWatcher instead
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

// Deprecated: use events.EventBufferWatcher instead
func (eb *EventBufferWatcher) Stop() {
	eb.signalChannel <- syscall.SIGTERM // TODO get the signal from the caller
	eb.ticker.Stop()
	eb.logger.Info("EventBufferWatcher stopped")
	close(eb.signalChannel)
	eb.signalChannel = nil
}

// Deprecated: use events.EventBufferWatcher instead
func (eb *EventBufferWatcher) Dispatch(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferWatcher.Dispatch")
	defer span.Finish()
	now := utils.Now()
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
			tracing.TraceErr(span, errors.Wrap(err, "error handling event"))
			eb.logger.Errorf("EventBufferWatcher.Dispatch: error handling event: %s", err.Error())
			continue
		}
		err = eb.ebr.Delete(&eventBuffer)
		if err != nil {
			return err
		}
	}
	return err
}

// Deprecated: use events.EventBufferWatcher instead
func (eb *EventBufferWatcher) HandleEvent(ctx context.Context, eventBuffer postgresEntity.EventBuffer) error {
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

// Deprecated: use events.EventBufferWatcher instead
func (eb *EventBufferWatcher) handleEvent(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferWatcher.handleEvent")
	defer span.Finish()

	switch evt.EventType {
	case orgevents.OrganizationUpdateOwnerNotificationV1:
		var data orgevents.OrganizationOwnerUpdateEvent
		if err := json.Unmarshal(evt.Data, &data); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		organizationAggregate, err := orgaggregate.LoadOrganizationAggregate(ctx, eb.es, data.Tenant, data.OrganizationId, eventstore.LoadAggregateOptions{})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		err = organizationAggregate.Apply(evt)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		// Persist the changes to the event store
		if err = eb.es.Save(ctx, organizationAggregate); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return err
	case reminder.ReminderNotificationV1:
		// TODO: implement the logic for the reminder and generic events
		return nil
	default:
		return errors.New("Event type not supported")
	}
}
