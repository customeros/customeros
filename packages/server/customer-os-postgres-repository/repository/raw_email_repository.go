package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type RawEmailRepository interface {
	CountForUsername(externalSystem, tenant, username string) (int64, error)
	EmailExistsByMessageId(externalSystem, tenant, username, messageId string) (bool, error)
	Store(externalSystem, tenant, username, providerMessageId, messageId, rawEmail string, sentAt time.Time, state entity.EmailImportState) error
}

type rawEmailRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawEmailRepository(gormDb *gorm.DB) RawEmailRepository {
	return &rawEmailRepositoryImpl{gormDb: gormDb}
}

func (repo *rawEmailRepositoryImpl) CountForUsername(externalSystem, tenant, username string) (int64, error) {
	var result int64
	err := repo.gormDb.Model(entity.RawEmail{}).Where("external_system = ? AND tenant = ? AND username = ?", externalSystem, tenant, username).Count(&result).Error

	if err != nil {
		logrus.Errorf("RawEmailRepository.CountForUsername - failed: %s; %s; %s", externalSystem, tenant, username)
		return 0, err
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) EmailExistsByMessageId(externalSystem, tenant, username, messageId string) (bool, error) {
	var result int64
	err := repo.gormDb.Model(entity.RawEmail{}).Where("external_system = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, username, messageId).Count(&result).Error

	if err != nil {
		logrus.Errorf("RawEmailRepository.EmailExistsByMessageId - failed: %s; %s; %s", externalSystem, tenant, messageId)
		return false, err
	}

	return result > 0, nil
}

func (repo *rawEmailRepositoryImpl) Store(externalSystem, tenant, username, providerMessageId, messageId, rawEmail string, sentAt time.Time, state entity.EmailImportState) error {
	result := entity.RawEmail{}
	err := repo.gormDb.Find(&result, "external_system = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, username, messageId).Error

	if err != nil {
		logrus.Errorf("RawEmailRepository.Store - failed retrieving rawEmail: %s; %s; %s; %s", externalSystem, tenant, username, messageId)
		return err
	}

	if result.Tenant != "" {
		logrus.Infof("RawEmailRepository.Store - already exists: %s; %s; %s; %s", externalSystem, tenant, username, messageId)
		return nil
	}

	result.ProviderMessageId = providerMessageId
	result.MessageId = messageId

	result.CreatedAt = utils.Now()
	result.SentAt = sentAt
	result.State = state
	result.ExternalSystem = externalSystem
	result.Tenant = tenant
	result.Username = username
	result.Data = rawEmail
	result.SentToEventStoreState = "PENDING"

	err = repo.gormDb.Save(&result).Error
	if err != nil {
		logrus.Errorf("RawEmailRepository.Store - failed storing rawEmail: %s; %s; %s; %s", externalSystem, tenant, username, messageId)
		return err
	}

	return nil
}
