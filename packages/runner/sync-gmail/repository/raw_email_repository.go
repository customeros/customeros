package repository

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RawEmailRepository interface {
	GetDistinctUsersForImport() ([]entity.RawEmail, error)
	GetEmailsIdsForSync(externalSystem, tenantName string) ([]entity.RawEmail, error)
	GetEmailsIdsForUserForSync(externalSystem, tenantName, userSource string) ([]entity.RawEmail, error)
	GetEmailForSync(id uuid.UUID) (*entity.RawEmail, error)
	GetEmailForSyncByMessageId(externalSystem, tenant, usernameSource, messageId string) (*entity.RawEmail, error)
	MarkSentToEventStore(id uuid.UUID, sentToEventStoreState entity.RawState, reason, error *string) error
}

type rawEmailRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawEmailRepository(gormDb *gorm.DB) RawEmailRepository {
	return &rawEmailRepositoryImpl{gormDb: gormDb}
}

func (repo *rawEmailRepositoryImpl) GetDistinctUsersForImport() ([]entity.RawEmail, error) {
	results := []entity.RawEmail{}

	err := repo.gormDb.Select("DISTINCT tenant, username").Where("sent_to_event_store_state = ?", "PENDING").Find(&results).Error

	if err != nil {
		logrus.Errorf("Failed getting distinct users for import")
		return nil, err
	}

	return results, nil
}

func (repo *rawEmailRepositoryImpl) GetEmailsIdsForSync(externalSystem, tenantName string) ([]entity.RawEmail, error) {
	result := []entity.RawEmail{}
	err := repo.gormDb.Order("sent_at desc").Select([]string{"id"}).Limit(25).Find(&result, "external_system = ? AND tenant = ? AND sent_to_event_store_state = 'PENDING'", externalSystem, tenantName).Error

	if err != nil {
		logrus.Errorf("Failed getting rawEmails: %s; %s", externalSystem, tenantName)
		return nil, err
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) GetEmailsIdsForUserForSync(externalSystem, tenantName, userSource string) ([]entity.RawEmail, error) {
	result := []entity.RawEmail{}
	err := repo.gormDb.Order("sent_at desc").Select([]string{"id"}).Limit(25).Find(&result, "external_system = ? AND tenant = ? AND username = ? AND sent_to_event_store_state = 'PENDING'", externalSystem, tenantName, userSource).Error

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

func (repo *rawEmailRepositoryImpl) GetEmailForSyncByMessageId(externalSystem, tenant, usernameSource, messageId string) (*entity.RawEmail, error) {
	var result entity.RawEmail
	err := repo.gormDb.Where("external_system = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, usernameSource, messageId).Find(&result).Error

	if err != nil {
		logrus.Errorf("GetEmailForSyncByMessageId - failed: %s; %s; %s", externalSystem, tenant, messageId)
		return nil, err
	}

	return &result, nil
}

func (repo *rawEmailRepositoryImpl) MarkSentToEventStore(id uuid.UUID, sentToEventStoreState entity.RawState, reason, error *string) error {
	tx := repo.gormDb.Model(&entity.RawEmail{}).Where("id = ?", id)

	tx.Update("sent_to_event_store_state", sentToEventStoreState)
	tx.Update("sent_to_event_store_reason", reason)
	tx.Update("sent_to_event_store_error", error)

	err := tx.Error

	if err != nil {
		logrus.Errorf("Failed marking email as sent to event store: %v", id)
		return err
	}

	return nil
}
