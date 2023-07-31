package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RawEmailRepository interface {
	Store(externalSystem, tenantName, usernameSource, messageId, rawEmail string) error
}

type rawEmailRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawEmailRepository(gormDb *gorm.DB) RawEmailRepository {
	return &rawEmailRepositoryImpl{gormDb: gormDb}
}

func (repo *rawEmailRepositoryImpl) Store(externalSystem, tenantName, usernameSource, messageId, rawEmail string) error {
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

	result.ExternalSystem = externalSystem
	result.TenantName = tenantName
	result.UsernameSource = usernameSource
	result.MessageId = messageId
	result.Data = rawEmail

	err = repo.gormDb.Save(&result).Error
	if err != nil {
		logrus.Errorf("Failed storing rawEmail: %s; %s; %s; %s", externalSystem, tenantName, usernameSource, messageId)
		return err
	}

	return nil
}
