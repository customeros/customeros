package repository

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RawEmailRepository interface {
	GetEmailsIdsForSync(externalSystem, tenantName string) ([]entity.RawEmail, error)
	GetEmailForSync(id uuid.UUID) (*entity.RawEmail, error)
	MarkSentToEventStore(id uuid.UUID, sentToEventStore bool, error string) error
}

type rawEmailRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawEmailRepository(gormDb *gorm.DB) RawEmailRepository {
	return &rawEmailRepositoryImpl{gormDb: gormDb}
}

func (repo *rawEmailRepositoryImpl) GetEmailsIdsForSync(externalSystem, tenantName string) ([]entity.RawEmail, error) {
	result := []entity.RawEmail{}
	err := repo.gormDb.Select("id").Find(&result, "external_system = ? AND tenant_name = ? AND sent_to_event_store = false", externalSystem, tenantName).Limit(25).Error

	if err != nil {
		logrus.Errorf("Failed getting rawEmails: %s; %s", externalSystem, tenantName)
		return nil, err
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) GetEmailForSync(id uuid.UUID) (*entity.RawEmail, error) {
	result := entity.RawEmail{}
	err := repo.gormDb.First(&result, id).Error

	if err != nil {
		logrus.Errorf("Failed getting rawEmail: %s", id)
		return nil, err
	}

	return &result, nil
}

func (repo *rawEmailRepositoryImpl) MarkSentToEventStore(id uuid.UUID, sentToEventStore bool, error string) error {
	tx := repo.gormDb.Model(&entity.RawEmail{}).Where("id = ?", id)

	tx.Update("sent_to_event_store", sentToEventStore)

	if error != "" {
		tx.Update("sent_to_event_store_error", error)
	} else {
		tx.Update("sent_to_event_store_error", nil)
	}

	err := tx.Error

	if err != nil {
		logrus.Errorf("Failed marking email as sent to event store: %v", id)
		return err
	}

	return nil
}
