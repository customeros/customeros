package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"gorm.io/gorm"
)

type UserGmailImportPageTokenRepository interface {
	GetGmailImportPageTokenForUsername(tenantName, username string) (*string, error)
	UpdateGmailImportPageTokenForUsername(tenantName, username, historyId string) error
}

type userGmailImportPageTokenImpl struct {
	gormDb *gorm.DB
}

func NewUserGmailImportPageTokenRepository(gormDb *gorm.DB) UserGmailImportPageTokenRepository {
	return &userGmailImportPageTokenImpl{gormDb: gormDb}
}

func (repo *userGmailImportPageTokenImpl) GetGmailImportPageTokenForUsername(tenantName, username string) (*string, error) {
	result := entity.UserGmailImportPageToken{}
	err := repo.gormDb.First(&result, "tenant_name = ? AND username = ?", tenantName, username).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// handle record not found error
			return nil, nil
		} else {
			return nil, fmt.Errorf("GetGmailImportPageTokenForUsername: %s", err.Error())
		}
	}
	return &result.HistoryId, nil
}

func (repo *userGmailImportPageTokenImpl) UpdateGmailImportPageTokenForUsername(tenantName, username, historyId string) error {
	result := entity.UserGmailImportPageToken{}
	err := repo.gormDb.Find(&result, "tenant_name = ? AND username = ?", tenantName, username).Error

	if err != nil {
		return fmt.Errorf("GetGmailImportPageTokenForUsername - update: %s", err.Error())
	}

	if result.TenantName == "" {
		result.TenantName = tenantName
		result.Username = username
	}

	result.HistoryId = historyId
	err = repo.gormDb.Save(&result).Error
	if err != nil {
		return fmt.Errorf("GetGmailImportPageTokenForUsername - insert: %s", err.Error())
	}

	return nil
}
