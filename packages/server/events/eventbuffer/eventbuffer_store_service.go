package eventbuffer

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"os"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type EventBufferStoreService struct {
	eventBufferRepository postgresRepository.EventBufferRepository
	logger                logger.Logger
	signalChannel         chan os.Signal
	ticker                *time.Ticker
}

func NewEventBufferStoreService(ebr postgresRepository.EventBufferRepository, logger logger.Logger) *EventBufferStoreService {
	return &EventBufferStoreService{eventBufferRepository: ebr, logger: logger}
}

// Deprecated
func (eb *EventBufferStoreService) Park(
	evt eventstore.Event,
	tenant string,
	uuid string,
	expiryTimestamp time.Time,
) error {
	eventBuffer := postgresEntity.EventBuffer{
		Tenant:             tenant,
		UUID:               uuid,
		ExpiryTimestamp:    expiryTimestamp.UTC(),
		EventID:            evt.EventID,
		EventType:          evt.EventType,
		EventData:          evt.Data,
		EventTimestamp:     evt.Timestamp.UTC(),
		EventAggregateID:   evt.AggregateID,
		EventAggregateType: string(evt.AggregateType),
		EventVersion:       evt.Version,
		EventMetadata:      evt.Metadata,
	}
	err := eb.eventBufferRepository.Upsert(&eventBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (eb *EventBufferStoreService) ParkBaseEvent(
	ctx context.Context,
	evt interface{},
	tenant string,
	expiryTimestamp time.Time,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventBufferStoreService.ParkBaseEvent")
	defer span.Finish()

	span.LogFields(log.Object("evt", evt))

	eventName := evt.(events.BaseEventAccessor).GetBaseEvent().EventName

	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	eventBuffer := postgresEntity.EventBuffer{
		UUID:            uuid.New().String(), // todo move to database generation
		Tenant:          tenant,
		ExpiryTimestamp: expiryTimestamp.UTC(),
		EventType:       eventName,
		EventData:       data,
	}
	err = eb.eventBufferRepository.Upsert(&eventBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (eb *EventBufferStoreService) GetById(uuid string) (*postgresEntity.EventBuffer, error) {
	return eb.eventBufferRepository.GetByUUID(uuid)
}

func (eb *EventBufferStoreService) Update(eventBuffer *postgresEntity.EventBuffer) error {
	return eb.eventBufferRepository.Upsert(eventBuffer)
}

func (eb *EventBufferStoreService) Delete(eventBuffer *postgresEntity.EventBuffer) error {
	return eb.eventBufferRepository.Delete(eventBuffer)
}
