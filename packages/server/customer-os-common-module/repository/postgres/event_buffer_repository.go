package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"time"

	"gorm.io/gorm"
)

type eventBufferRepository struct {
	gormDb *gorm.DB
}

type EventBufferRepository interface {
	Upsert(eventBuffer *entity.EventBuffer) error
	GetByExpired(now time.Time) ([]entity.EventBuffer, error)
	GetByUUID(uuid string) (*entity.EventBuffer, error)
	Delete(eventBuffer *entity.EventBuffer) error
}

func NewEventBufferRepository(gormDb *gorm.DB) EventBufferRepository {
	return &eventBufferRepository{gormDb: gormDb}
}

func (repo *eventBufferRepository) Upsert(eventBuffer *entity.EventBuffer) error {
	return repo.gormDb.Save(eventBuffer).Error
}

func (repo *eventBufferRepository) GetByExpired(now time.Time) ([]entity.EventBuffer, error) {
	var eventBuffers []entity.EventBuffer
	err := repo.gormDb.Where("expiry_timestamp < ?", now).Find(&eventBuffers).Error
	return eventBuffers, err
}

func (repo *eventBufferRepository) GetByUUID(uuid string) (*entity.EventBuffer, error) {
	var eventBuffer *entity.EventBuffer
	err := repo.gormDb.Where("uuid = ?", uuid).Find(&eventBuffer).Error
	return eventBuffer, err
}

func (repo *eventBufferRepository) Delete(eventBuffer *entity.EventBuffer) error {
	return repo.gormDb.Delete(eventBuffer).Error
}
