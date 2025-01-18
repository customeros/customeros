package repository

import (
	"github.com/google/uuid"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/runner/sync-gmail/config"
	"github.com/customeros/customeros/packages/runner/sync-gmail/entity"
)

type RawEmailRepository interface {
	GetDistinctUsersForImport() ([]entity.RawEmail, error)
	GetEmailsIdsForSync(externalSystem, tenantName string) ([]entity.RawEmail, error)
	GetEmailsIdsForUserForSync(tenantName, userSource string) ([]entity.RawEmail, error)
	GetEmailForSync(id uuid.UUID) (*entity.RawEmail, error)
	GetEmailForSyncByMessageId(tenant, usernameSource, messageId string) (*entity.RawEmail, error)
	MarkSentToEventStore(id uuid.UUID, sentToEventStoreState postgresentity.RawState, reason, error *string) error
}

type rawEmailRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawEmailRepository(gormDb *gorm.DB) RawEmailRepository {
	return &rawEmailRepositoryImpl{gormDb: gormDb}
}

func (repo *rawEmailRepositoryImpl) GetDistinctUsersForImport() ([]entity.RawEmail, error) {
	results := []entity.RawEmail{}

	err := repo.gormDb.Select("DISTINCT tenant, username").Where("status = ?", "PENDING").Find(&results).Error
	if err != nil {
		logrus.Errorf("Failed getting distinct users for import")
		return nil, err
	}

	return results, nil
}

func (repo *rawEmailRepositoryImpl) GetEmailsIdsForSync(externalSystem, tenantName string) ([]entity.RawEmail, error) {
	result := []entity.RawEmail{}
	err := repo.gormDb.Order("sent_at desc").Select([]string{"id"}).Limit(config.MAX_EMAILS_PER_RUN).Find(&result, "external_system = ? AND tenant = ? AND status = 'PENDING'", externalSystem, tenantName).Error
	if err != nil {
		logrus.Errorf("Failed getting rawEmails: %s; %s", externalSystem, tenantName)
		return nil, err
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) GetEmailsIdsForUserForSync(tenantName, userSource string) ([]entity.RawEmail, error) {
	result := []entity.RawEmail{}
	err := repo.gormDb.Order("sent_at desc").Select([]string{"id", "external_system"}).Limit(config.MAX_EMAILS_PER_RUN).Find(&result, "tenant = ? AND username = ? AND status = 'PENDING'", tenantName, userSource).Error
	if err != nil {
		logrus.Errorf("Failed getting rawEmails: %s; %s", tenantName, userSource)
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

func (repo *rawEmailRepositoryImpl) GetEmailForSyncByMessageId(tenant, usernameSource, messageId string) (*entity.RawEmail, error) {
	var result entity.RawEmail
	err := repo.gormDb.Where("tenant = ? AND username = ? AND message_id = ?", tenant, usernameSource, messageId).Find(&result).Error
	if err != nil {
		logrus.Errorf("GetEmailForSyncByMessageId - failed: %s; %s; %s", tenant, usernameSource, messageId)
		return nil, err
	}

	return &result, nil
}

func (repo *rawEmailRepositoryImpl) MarkSentToEventStore(id uuid.UUID, sentToEventStoreState postgresentity.RawState, reason, error *string) error {
	tx := repo.gormDb.Model(&entity.RawEmail{}).Where("id = ?", id)

	tx.Update("status", sentToEventStoreState)
	tx.Update("reason", reason)
	tx.Update("error", error)

	err := tx.Error
	if err != nil {
		logrus.Errorf("Failed marking email as sent to event store: %v", id)
		return err
	}

	return nil
}
