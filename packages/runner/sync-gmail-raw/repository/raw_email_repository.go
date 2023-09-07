package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type RawEmailRepository interface {
	EmailExistsByMessageId(externalSystem, tenant, usernameSource, messageId string) (bool, error)
	Store(externalSystem, tenantName, usernameSource, providerMessageId, messageId, rawEmail string) error
}

type rawEmailRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawEmailRepository(gormDb *gorm.DB) RawEmailRepository {
	return &rawEmailRepositoryImpl{gormDb: gormDb}
}

func (repo *rawEmailRepositoryImpl) EmailExistsByMessageId(externalSystem, tenant, usernameSource, messageId string) (bool, error) {
	var result int64
	err := repo.gormDb.Model(entity.RawEmail{}).Where("external_system = ? AND tenant_name = ? AND username_source = ? AND message_id = ?", externalSystem, tenant, usernameSource, messageId).Count(&result).Error

	if err != nil {
		logrus.Errorf("Failed getting rawEmail: %s; %s; %s", externalSystem, tenant, messageId)
		return false, err
	}

	return result > 0, nil
}

func (repo *rawEmailRepositoryImpl) Store(externalSystem, tenantName, usernameSource, providerMessageId, messageId, rawEmail string) error {
	result := entity.RawEmail{}
	err := repo.gormDb.Find(&result, "external_system = ? AND tenant_name = ? AND username_source = ? AND message_id = ?", externalSystem, tenantName, usernameSource, messageId).Error

	if err != nil {
		logrus.Errorf("Failed retrieving rawEmail: %s; %s; %s; %s", externalSystem, tenantName, usernameSource, messageId)
		return err
	}

	if result.TenantName != "" {
		logrus.Infof("RawEmailRepository.Store - already exists: %s; %s; %s; %s", externalSystem, tenantName, usernameSource, messageId)
		return nil
	}

	result.ProviderMessageId = providerMessageId
	result.MessageId = messageId

	result.CreatedAt = time.Now().UTC()
	result.ExternalSystem = externalSystem
	result.TenantName = tenantName
	result.UsernameSource = usernameSource
	result.Data = rawEmail
	result.SentToEventStoreState = "PENDING"

	err = repo.gormDb.Save(&result).Error
	if err != nil {
		logrus.Errorf("Failed storing rawEmail: %s; %s; %s; %s", externalSystem, tenantName, usernameSource, messageId)
		return err
	}

	return nil
}
