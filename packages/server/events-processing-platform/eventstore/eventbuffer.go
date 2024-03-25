package eventstore

import (
	"context"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"time"
)

type EventBufferService struct {
	eventBufferRepository repository.EventBufferRepository
}

func NewEventBufferService(eventBufferRepository repository.EventBufferRepository) *EventBufferService {
	return &EventBufferService{eventBufferRepository: eventBufferRepository}
}

func (eb *EventBufferService) Park(
	ctx context.Context,
	evt Event,
	tenant string,
	uuid string,
	expiryTimestamp time.Time,
) error {
	eventBuffer := entity.EventBuffer{
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
	err := eb.eventBufferRepository.Upsert(eventBuffer)
	if err != nil {
		return err
	}
	return nil
}
