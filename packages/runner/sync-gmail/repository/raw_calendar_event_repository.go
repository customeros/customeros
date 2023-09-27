package repository

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RawCalendarEventRepository interface {
	GetCalendarEventsIdsForUserForSync(externalSystem, tenantName, userSource string) ([]entity.RawCalendarEvent, error)
	GetCalendarEventForSync(id uuid.UUID) (*entity.RawCalendarEvent, error)
	MarkSentToEventStore(id uuid.UUID, sentToEventStoreState entity.RawState, reason, error *string) error
}

type rawCalendarEventRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawCalendarEventRepository(gormDb *gorm.DB) RawCalendarEventRepository {
	return &rawCalendarEventRepositoryImpl{gormDb: gormDb}
}

func (repo *rawCalendarEventRepositoryImpl) GetCalendarEventsIdsForUserForSync(externalSystem, tenantName, userSource string) ([]entity.RawCalendarEvent, error) {
	result := []entity.RawCalendarEvent{}
	err := repo.gormDb.Select("id").Limit(25).Find(&result, "external_system = ? AND tenant_name = ? AND username_source = ? AND sent_to_event_store_state = 'PENDING'", externalSystem, tenantName, userSource).Error

	if err != nil {
		logrus.Errorf("Failed getting rawCalendarEvents: %s; %s", externalSystem, tenantName)
		return nil, err
	}

	return result, nil
}

func (repo *rawCalendarEventRepositoryImpl) GetCalendarEventForSync(id uuid.UUID) (*entity.RawCalendarEvent, error) {
	result := entity.RawCalendarEvent{}
	err := repo.gormDb.First(&result, id).Error

	if err != nil {
		logrus.Errorf("Failed getting rawCalendarEvent: %s", id)
		return nil, err
	}

	return &result, nil
}

func (repo *rawCalendarEventRepositoryImpl) MarkSentToEventStore(id uuid.UUID, sentToEventStoreState entity.RawState, reason, error *string) error {
	tx := repo.gormDb.Model(&entity.RawCalendarEvent{}).Where("id = ?", id)

	tx.Update("sent_to_event_store_state", sentToEventStoreState)
	tx.Update("sent_to_event_store_reason", reason)
	tx.Update("sent_to_event_store_error", error)

	err := tx.Error

	if err != nil {
		logrus.Errorf("Failed marking calendar event as sent to event store: %v", id)
		return err
	}

	return nil
}
