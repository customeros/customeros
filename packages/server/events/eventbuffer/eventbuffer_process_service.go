package eventbuffer

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
)

type EventBufferProcessService struct {
	eventBufferRepository postgresRepository.EventBufferRepository
	logger                logger.Logger
	grpc_clients          *grpc_client.Clients
	signalChannel         chan os.Signal
	ticker                *time.Ticker
}

func NewEventBufferService(ebr postgresRepository.EventBufferRepository, logger logger.Logger, grpc_clients *grpc_client.Clients) *EventBufferProcessService {
	return &EventBufferProcessService{eventBufferRepository: ebr, logger: logger, grpc_clients: grpc_clients}
}

func (eb *EventBufferProcessService) Start(ctx context.Context) {
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
func (eb *EventBufferProcessService) Stop() {
	eb.signalChannel <- syscall.SIGTERM // TODO get the signal from the caller
	eb.ticker.Stop()
	eb.logger.Info("EventBufferWatcher stopped")
	close(eb.signalChannel)
	eb.signalChannel = nil
}

// Dispatch dispatches all expired events from event_buffer table, and delete them after dispatching
func (eb *EventBufferProcessService) Dispatch(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferWatcher.Dispatch")
	defer span.Finish()
	now := time.Now().UTC()
	eventBuffers, err := eb.eventBufferRepository.GetByExpired(now)
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
		err = eb.eventBufferRepository.Delete(&eventBuffer)
		if err != nil {
			return err
		}
	}
	return err
}

// HandleEvent loads the event aggregate and applies the event to it and pushes it into event store
func (eb *EventBufferProcessService) HandleEvent(ctx context.Context, eventBuffer postgresEntity.EventBuffer) error {
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

func (eb *EventBufferProcessService) handleEvent(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferWatcher.handleEvent")
	defer span.Finish()

	dataBytes, err := json.Marshal(evt)
	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
		return err
	}

	//skip these 2 events that are handled by subscribers until we migrate and test them
	if evt.EventType == "V1_ORGANIZATION_UPDATE_OWNER_NOTIFICATION" || evt.EventType == "V1_REMINDER_NOTIFICATION" {
		return errors.New("Event type not supported")
	}

	_, err = eb.grpc_clients.EventStoreClient.StoreEvent(context.Background(), &eventstorepb.StoreEventGrpcRequest{
		EventData: string(dataBytes),
	})

	if err != nil {
		log.Fatalf("Failed: %v", err.Error())
		return err
	}

	return nil
}
