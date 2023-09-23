package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"gorm.io/gorm"
	"time"
)

type UserGCalImportStateRepository interface {
	GetGCalImportStateUsername(tenantName, username, calendarId string) (*entity.UserGCalImportState, error)
	UpdateGCalImportStateForUsername(tenantName, username, calendarId, syncToken, pageToken string, maxResults int64, timeMin, timeMax time.Time) error
}

type userGCalImportStateImpl struct {
	gormDb *gorm.DB
}

func NewUserGCalImportStateRepository(gormDb *gorm.DB) UserGCalImportStateRepository {
	return &userGCalImportStateImpl{gormDb: gormDb}
}

func (repo *userGCalImportStateImpl) GetGCalImportStateUsername(tenantName, username, calendarId string) (*entity.UserGCalImportState, error) {
	result := entity.UserGCalImportState{}
	err := repo.gormDb.First(&result, "tenant_name = ? AND username = ? AND calendar_id = ?", tenantName, username, calendarId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// handle record not found error
			return nil, nil
		} else {
			return nil, fmt.Errorf("GetGCalImportPageTokenForUsername: %s", err.Error())
		}
	}
	return &result, nil
}

func (repo *userGCalImportStateImpl) UpdateGCalImportStateForUsername(tenantName, username, calendarId, syncToken, pageToken string, maxResults int64, timeMin, timeMax time.Time) error {
	result := entity.UserGCalImportState{}
	err := repo.gormDb.Find(&result, "tenant_name = ? AND username = ? AND calendar_id = ?", tenantName, username, calendarId).Error

	if err != nil {
		return fmt.Errorf("UpdateGCalImportPageTokenForUsername - find: %s", err.Error())
	}

	if result.TenantName == "" {
		result.TenantName = tenantName
		result.Username = username
		result.CalendarId = calendarId
	}

	result.SyncToken = syncToken
	result.PageToken = pageToken
	result.MaxResults = maxResults
	result.TimeMin = timeMin
	result.TimeMax = timeMax

	err = repo.gormDb.Save(&result).Error
	if err != nil {
		return fmt.Errorf("UpdateGCalImportPageTokenForUsername - insert: %s", err.Error())
	}

	return nil
}
