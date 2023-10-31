package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"gorm.io/gorm"
	"time"
)

type UserGmailImportStateRepository interface {
	GetGmailImportState(tenantName, username string, state entity.GmailImportState) (*entity.UserGmailImportState, error)
	CreateGmailImportState(tenantName, username string, state entity.GmailImportState, startDate, stopDate *time.Time, active bool, cursor string) (*entity.UserGmailImportState, error)
	UpdateGmailImportState(tenantName, username string, state entity.GmailImportState, cursor string) (*entity.UserGmailImportState, error)
	ActivateGmailImportState(tenantName, username string, state entity.GmailImportState) error
	DeactivateGmailImportState(tenantName, username string, state entity.GmailImportState) error
}

type userGmailImportStateImpl struct {
	gormDb *gorm.DB
}

func NewUserGmailImportStateRepository(gormDb *gorm.DB) UserGmailImportStateRepository {
	return &userGmailImportStateImpl{gormDb: gormDb}
}

func (repo *userGmailImportStateImpl) GetGmailImportState(tenantName, username string, state entity.GmailImportState) (*entity.UserGmailImportState, error) {
	result := entity.UserGmailImportState{}
	err := repo.gormDb.First(&result, "tenant = ? AND username = ? AND state = ?", tenantName, username, state).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// handle record not found error
			return nil, nil
		} else {
			return nil, fmt.Errorf("GetGmailImportState: %s", err.Error())
		}
	}
	return &result, nil
}

func (repo *userGmailImportStateImpl) CreateGmailImportState(tenantName, username string, state entity.GmailImportState, startDate, stopDate *time.Time, active bool, cursor string) (*entity.UserGmailImportState, error) {
	result := entity.UserGmailImportState{}
	err := repo.gormDb.Find(&result, "tenant = ? AND username = ? AND state = ?", tenantName, username, state).Error

	if err != nil {
		return nil, fmt.Errorf("GetGmailImportState - update: %s", err.Error())
	}

	if result.Tenant != "" {
		return nil, fmt.Errorf("GetGmailImportState - already exists: %s; %s; %s", tenantName, username, state)
	}

	result.Tenant = tenantName
	result.Username = username
	result.State = state
	result.StartDate = startDate
	result.StopDate = stopDate
	result.Active = active

	result.Cursor = cursor
	err = repo.gormDb.Save(&result).Error
	if err != nil {
		return nil, fmt.Errorf("UpdateGmailImportState - insert: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(&result)
	if err != nil {
		return nil, fmt.Errorf("UpdateGmailImportState - insert history: %s", err.Error())
	}

	return repo.GetGmailImportState(tenantName, username, state)
}

func (repo *userGmailImportStateImpl) UpdateGmailImportState(tenantName, username string, state entity.GmailImportState, cursor string) (*entity.UserGmailImportState, error) {
	gmailImportState, err := repo.GetGmailImportState(tenantName, username, state)
	if gmailImportState == nil {
		return nil, fmt.Errorf("UpdateGmailImportState - not found: %s; %s; %s", tenantName, username, state)
	}

	gmailImportState.Cursor = cursor
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return nil, fmt.Errorf("UpdateGmailImportState - insert: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(gmailImportState)
	if err != nil {
		return nil, fmt.Errorf("UpdateGmailImportState - insert history: %s", err.Error())
	}

	return repo.GetGmailImportState(tenantName, username, state)
}

func (repo *userGmailImportStateImpl) ActivateGmailImportState(tenantName, username string, state entity.GmailImportState) error {
	gmailImportState, err := repo.GetGmailImportState(tenantName, username, state)
	if gmailImportState == nil {
		return nil
	}

	gmailImportState.Active = true
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return fmt.Errorf("DeactivateGmailImportState - update: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(gmailImportState)
	if err != nil {
		return fmt.Errorf("DeactivateGmailImportState - insert history: %s", err.Error())
	}

	return nil
}

func (repo *userGmailImportStateImpl) DeactivateGmailImportState(tenantName, username string, state entity.GmailImportState) error {
	gmailImportState, err := repo.GetGmailImportState(tenantName, username, state)
	if gmailImportState == nil {
		return nil
	}

	gmailImportState.Active = false
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return fmt.Errorf("DeactivateGmailImportState - update: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(gmailImportState)
	if err != nil {
		return fmt.Errorf("DeactivateGmailImportState - insert history: %s", err.Error())
	}

	return nil
}

func (repo *userGmailImportStateImpl) InsertHistoryRecord(record *entity.UserGmailImportState) error {
	history := entity.UserGmailImportStateHistory{}
	history.CreatedAt = time.Now().UTC()
	history.Tenant = record.Tenant
	history.Username = record.Username
	history.State = record.State
	history.StartDate = record.StartDate
	history.StopDate = record.StopDate
	history.EntityId = record.ID
	history.Cursor = record.Cursor
	history.Active = record.Active

	return repo.gormDb.Create(&history).Error
}
