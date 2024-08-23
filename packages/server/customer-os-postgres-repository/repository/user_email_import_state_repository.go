package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type UserEmailImportStateRepository interface {
	GetEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState) (*entity.UserEmailImportState, error)
	CreateEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState, startDate, stopDate *time.Time, active bool, cursor string) (*entity.UserEmailImportState, error)
	UpdateEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState, cursor string) (*entity.UserEmailImportState, error)
	ActivateEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState) error
	DeactivateEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState) error
}

type userEmailImportStateImpl struct {
	gormDb *gorm.DB
}

func NewUserEmailImportStateRepository(gormDb *gorm.DB) UserEmailImportStateRepository {
	return &userEmailImportStateImpl{gormDb: gormDb}
}

func (repo *userEmailImportStateImpl) GetEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState) (*entity.UserEmailImportState, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEmailImportStateRepository.GetEmailImportState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenantName)

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

func (repo *userEmailImportStateImpl) CreateEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState, startDate, stopDate *time.Time, active bool, cursor string) (*entity.UserEmailImportState, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEmailImportStateRepository.CreateEmailImportState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenantName)

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

	err = repo.InsertHistoryRecord(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("UpdateEmailImportState - insert history: %s", err.Error())
	}

	return repo.GetEmailImportState(ctx, tenantName, provider, username, state)
}

func (repo *userEmailImportStateImpl) UpdateEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState, cursor string) (*entity.UserEmailImportState, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEmailImportStateRepository.UpdateEmailImportState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenantName)

	gmailImportState, err := repo.GetEmailImportState(ctx, tenantName, provider, username, state)
	if gmailImportState == nil {
		return nil, fmt.Errorf("UpdateEmailImportState - not found: %s; %s; %s", tenantName, username, state)
	}

	gmailImportState.Cursor = cursor
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return nil, fmt.Errorf("UpdateEmailImportState - insert: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(ctx, gmailImportState)
	if err != nil {
		return nil, fmt.Errorf("UpdateEmailImportState - insert history: %s", err.Error())
	}

	return repo.GetEmailImportState(ctx, tenantName, provider, username, state)
}

func (repo *userEmailImportStateImpl) ActivateEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEmailImportStateRepository.ActivateEmailImportState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenantName)

	gmailImportState, err := repo.GetEmailImportState(ctx, tenantName, provider, username, state)
	if gmailImportState == nil {
		return nil
	}

	gmailImportState.Active = true
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return fmt.Errorf("DeactivateEmailImportState - update: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(ctx, gmailImportState)
	if err != nil {
		return fmt.Errorf("DeactivateEmailImportState - insert history: %s", err.Error())
	}

	return nil
}

func (repo *userEmailImportStateImpl) DeactivateEmailImportState(ctx context.Context, tenantName, provider, username string, state entity.EmailImportState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEmailImportStateRepository.DeactivateEmailImportState")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenantName)

	gmailImportState, err := repo.GetEmailImportState(ctx, tenantName, provider, username, state)
	if gmailImportState == nil {
		return nil
	}

	gmailImportState.Active = false
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		return fmt.Errorf("DeactivateEmailImportState - update: %s", err.Error())
	}

	err = repo.InsertHistoryRecord(ctx, gmailImportState)
	if err != nil {
		return fmt.Errorf("DeactivateEmailImportState - insert history: %s", err.Error())
	}

	return nil
}

func (repo *userEmailImportStateImpl) InsertHistoryRecord(ctx context.Context, record *entity.UserEmailImportState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserEmailImportStateRepository.InsertHistoryRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

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
