package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
	"time"
)

type UserEmailImportStateRepository interface {
	GetEmailImportState(tenantName, provider, username string, state entity.EmailImportState) (*entity.UserEmailImportState, error)
	CreateEmailImportState(tenantName, provider, username string, state entity.EmailImportState, startDate, stopDate *time.Time, active bool, cursor string) (*entity.UserEmailImportState, error)
	UpdateEmailImportState(tenantName, provider, username string, state entity.EmailImportState, cursor string) (*entity.UserEmailImportState, error)
	ActivateEmailImportState(tenantName, provider, username string, state entity.EmailImportState) error
	DeactivateEmailImportState(tenantName, provider, username string, state entity.EmailImportState) error
}

type userEmailImportStateImpl struct {
	gormDb *gorm.DB
}

func NewUserEmailImportStateRepository(gormDb *gorm.DB) UserEmailImportStateRepository {
	return &userEmailImportStateImpl{gormDb: gormDb}
}

func (repo *userEmailImportStateImpl) GetEmailImportState(tenantName, provider, username string, state entity.EmailImportState) (*entity.UserEmailImportState, error) {
	result := entity.UserEmailImportState{}
	err := repo.gormDb.First(&result, "tenant = ? AND provider = ? AND username = ? AND state = ?", tenantName, provider, username, state).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// handle record not found error
			return nil, nil
		} else {
			return nil, fmt.Errorf("GetEmailImportState: %s", err.Error())
		}
	}
	return &result, nil
}

func (repo *userEmailImportStateImpl) CreateEmailImportState(tenantName, provider, username string, state entity.EmailImportState, startDate, stopDate *time.Time, active bool, cursor string) (*entity.UserEmailImportState, error) {
	result := entity.UserEmailImportState{}
	err := repo.gormDb.Find(&result, "tenant = ? AND provider = ? AND username = ? AND state = ?", tenantName, provider, username, state).Error

	if err != nil {
		return nil, fmt.Errorf("GetEmailImportState - update: %s", err.Error())
	}

	if result.Tenant != "" {
		return nil, fmt.Errorf("GetEmailImportState - already exists: %s; %s; %s", tenantName, username, state)
	}

	result.Tenant = tenantName
	result.Username = username
	result.Provider = provider
	result.State = state
	result.StartDate = startDate
	result.StopDate = stopDate
	result.Active = active

	result.Cursor = cursor
	err = repo.gormDb.Save(&result).Error
	if err != nil {
		return nil, fmt.Errorf("UpdateEmailImportState - insert: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(&result)
	if err != nil {
		return nil, fmt.Errorf("UpdateEmailImportState - insert history: %s", err.Error())
	}

	return repo.GetEmailImportState(tenantName, provider, username, state)
}

func (repo *userEmailImportStateImpl) UpdateEmailImportState(tenantName, provider, username string, state entity.EmailImportState, cursor string) (*entity.UserEmailImportState, error) {
	gmailImportState, err := repo.GetEmailImportState(tenantName, provider, username, state)
	if gmailImportState == nil {
		return nil, fmt.Errorf("UpdateEmailImportState - not found: %s; %s; %s", tenantName, username, state)
	}

	gmailImportState.Cursor = cursor
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return nil, fmt.Errorf("UpdateEmailImportState - insert: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(gmailImportState)
	if err != nil {
		return nil, fmt.Errorf("UpdateEmailImportState - insert history: %s", err.Error())
	}

	return repo.GetEmailImportState(tenantName, provider, username, state)
}

func (repo *userEmailImportStateImpl) ActivateEmailImportState(tenantName, provider, username string, state entity.EmailImportState) error {
	gmailImportState, err := repo.GetEmailImportState(tenantName, provider, username, state)
	if gmailImportState == nil {
		return nil
	}

	gmailImportState.Active = true
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return fmt.Errorf("DeactivateEmailImportState - update: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(gmailImportState)
	if err != nil {
		return fmt.Errorf("DeactivateEmailImportState - insert history: %s", err.Error())
	}

	return nil
}

func (repo *userEmailImportStateImpl) DeactivateEmailImportState(tenantName, provider, username string, state entity.EmailImportState) error {
	gmailImportState, err := repo.GetEmailImportState(tenantName, provider, username, state)
	if gmailImportState == nil {
		return nil
	}

	gmailImportState.Active = false
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return fmt.Errorf("DeactivateEmailImportState - update: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(gmailImportState)
	if err != nil {
		return fmt.Errorf("DeactivateEmailImportState - insert history: %s", err.Error())
	}

	return nil
}

func (repo *userEmailImportStateImpl) InsertHistoryRecord(record *entity.UserEmailImportState) error {
	history := entity.UserEmailImportStateHistory{}
	history.CreatedAt = utils.Now()
	history.Tenant = record.Tenant
	history.Username = record.Username
	history.Provider = record.Provider
	history.State = record.State
	history.StartDate = record.StartDate
	history.StopDate = record.StopDate
	history.EntityId = record.ID
	history.Cursor = record.Cursor
	history.Active = record.Active

	return repo.gormDb.Create(&history).Error
}
