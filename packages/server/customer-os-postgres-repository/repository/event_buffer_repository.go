package postgres_repository

import (
	"errors"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"time"

	"gorm.io/gorm"
)

type eventBufferRepository struct {
	gormDb *gorm.DB
}

type EventBufferRepository interface {
	Upsert(eventBuffer *postgres_entity.EventBuffer) error
	GetByExpired(now time.Time) ([]postgres_entity.EventBuffer, error)
	GetByUUID(uuid string) (*postgres_entity.EventBuffer, error)
	Delete(eventBuffer *postgres_entity.EventBuffer) error
}

func NewEventBufferRepository(gormDb *gorm.DB) EventBufferRepository {
	return &eventBufferRepository{gormDb: gormDb}
}

func (repo *eventBufferRepository) Upsert(eventBuffer *postgres_entity.EventBuffer) error {
	return repo.gormDb.Save(eventBuffer).Error
}

func (repo *eventBufferRepository) GetByExpired(now time.Time) ([]postgres_entity.EventBuffer, error) {
	var eventBuffers []postgres_entity.EventBuffer
	err := repo.gormDb.Where("expiry_timestamp < ?", now).Find(&eventBuffers).Error
	return eventBuffers, err
}

func (repo *eventBufferRepository) GetByUUID(uuid string) (*postgres_entity.EventBuffer, error) {
	var eventBuffer postgres_entity.EventBuffer
	err := repo.gormDb.Where("uuid = ?", uuid).First(&eventBuffer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &eventBuffer, err
}

func (repo *eventBufferRepository) Delete(eventBuffer *postgres_entity.EventBuffer) error {
	return repo.gormDb.Delete(eventBuffer).Error
}
