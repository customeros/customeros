package eventstore

import (
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
	err := eb.eventBufferRepository.Upsert(&eventBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (eb *EventBufferService) GetById(uuid string) (*entity.EventBuffer, error) {
	return eb.eventBufferRepository.GetByUUID(uuid)
}

func (eb *EventBufferService) Update(eventBuffer *entity.EventBuffer) error {
	return eb.eventBufferRepository.Upsert(eventBuffer)
}

func (eb *EventBufferService) Delete(eventBuffer *entity.EventBuffer) error {
	return eb.eventBufferRepository.Delete(eventBuffer)
}
