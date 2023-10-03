package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"gorm.io/gorm"
)

type UserGmailImportStateRepository interface {
	GetGmailImportState(tenantName, username string) (*string, error)
	UpdateGmailImportState(tenantName, username, historyId string) error
}

type userGmailImportStateImpl struct {
	gormDb *gorm.DB
}

func NewUserGmailImportStateRepository(gormDb *gorm.DB) UserGmailImportStateRepository {
	return &userGmailImportStateImpl{gormDb: gormDb}
}

func (repo *userGmailImportStateImpl) GetGmailImportState(tenantName, username string) (*string, error) {
	result := entity.UserGmailImportState{}
	err := repo.gormDb.First(&result, "tenant_name = ? AND username = ?", tenantName, username).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// handle record not found error
			return nil, nil
		} else {
			return nil, fmt.Errorf("GetGmailImportState: %s", err.Error())
		}
	}
	return &result.HistoryId, nil
}

func (repo *userGmailImportStateImpl) UpdateGmailImportState(tenantName, username, historyId string) error {
	result := entity.UserGmailImportState{}
	err := repo.gormDb.Find(&result, "tenant_name = ? AND username = ?", tenantName, username).Error

	if err != nil {
		return fmt.Errorf("GetGmailImportState - update: %s", err.Error())
	}

	if result.TenantName == "" {
		result.TenantName = tenantName
		result.Username = username
	}

	result.HistoryId = historyId
	err = repo.gormDb.Save(&result).Error
	if err != nil {
		return fmt.Errorf("GetGmailImportState - insert: %s", err.Error())
	}

	history := entity.UserGmailImportStateHistory{}
	history.EntityId = result.ID
	history.TenantName = tenantName
	history.Username = username
	history.HistoryId = historyId
	err = repo.gormDb.Create(&history).Error

	if err != nil {
		return fmt.Errorf("GetGmailImportState - insert history: %s", err.Error())
	}

	return nil
}
