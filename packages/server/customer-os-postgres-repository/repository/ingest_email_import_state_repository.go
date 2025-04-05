package postgres_repository

import (
	"errors"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type IngestEmailImportStateRepository interface {
	GetEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod) (*postgres_entity.IngestEmailImportState, error)
	CreateEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod, startDate, stopDate *time.Time, active bool, cursor string) (*postgres_entity.IngestEmailImportState, error)
	UpdateEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod, cursor string) (*postgres_entity.IngestEmailImportState, error)
	ActivateEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod) error
	DeactivateEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod) error
}

type ingestEmailImportStateImpl struct {
	gormDb *gorm.DB
}

func NewIngestEmailImportStateRepository(gormDb *gorm.DB) IngestEmailImportStateRepository {
	return &ingestEmailImportStateImpl{gormDb: gormDb}
}

func (repo *ingestEmailImportStateImpl) GetEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod) (*postgres_entity.IngestEmailImportState, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailImportStateRepository.GetEmailImportState")
	defer spans.Finish()
	spans.LogKV("provider", provider, "username", username, "period", period)

	result := postgres_entity.IngestEmailImportState{}
	err := repo.gormDb.First(&result, "tenant = ? AND provider = ? AND username = ? AND period = ?", tenantName, provider, username, period).Error

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

func (repo *ingestEmailImportStateImpl) CreateEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod, startDate, stopDate *time.Time, active bool, cursor string) (*postgres_entity.IngestEmailImportState, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailImportStateRepository.CreateEmailImportState")
	defer spans.Finish()
	spans.LogKV("provider", provider, "username", username, "period", period)

	result := postgres_entity.IngestEmailImportState{}
	err := repo.gormDb.Find(&result, "tenant = ? AND provider = ? AND username = ? AND period = ?", tenantName, provider, username, period).Error

	if err != nil {
		return nil, fmt.Errorf("GetEmailImportState - update: %s", err.Error())
	}

	if result.Tenant != "" {
		return nil, fmt.Errorf("GetEmailImportState - already exists: %s; %s; %s", tenantName, username, period)
	}

	result.Tenant = tenantName
	result.Username = username
	result.Provider = provider
	result.Period = period
	result.StartDate = startDate
	result.StopDate = stopDate
	result.Active = active

	result.Cursor = cursor
	err = repo.gormDb.Save(&result).Error
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("UpdateEmailImportState - insert: %s", err.Error())
	}

	return repo.GetEmailImportState(ctx, tenantName, provider, username, period)
}

func (repo *ingestEmailImportStateImpl) UpdateEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod, cursor string) (*postgres_entity.IngestEmailImportState, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailImportStateRepository.UpdateEmailImportState")
	defer spans.Finish()
	spans.LogKV("provider", provider, "username", username, "period", period)

	gmailImportState, err := repo.GetEmailImportState(ctx, tenantName, provider, username, period)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if gmailImportState == nil {
		return nil, fmt.Errorf("UpdateEmailImportState - not found: %s; %s; %s", tenantName, username, period)
	}

	gmailImportState.Cursor = cursor
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("UpdateEmailImportState - insert: %s", err.Error())
	}

	return repo.GetEmailImportState(ctx, tenantName, provider, username, period)
}

func (repo *ingestEmailImportStateImpl) ActivateEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailImportStateRepository.ActivateEmailImportState")
	defer spans.Finish()
	spans.LogKV("provider", provider, "username", username, "period", period)

	gmailImportState, err := repo.GetEmailImportState(ctx, tenantName, provider, username, period)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	if gmailImportState == nil {
		return nil
	}

	gmailImportState.Active = true
	err = repo.gormDb.Save(&gmailImportState).Error
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("DeactivateEmailImportState - update: %s", err.Error())
	}

	return nil
}

func (repo *ingestEmailImportStateImpl) DeactivateEmailImportState(ctx context.Context, tenantName, provider, username string, period postgres_entity.IngestEmailImportStatePeriod) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "IngestEmailImportStateRepository.DeactivateEmailImportState")
	defer spans.Finish()
	spans.LogKV("provider", provider, "username", username, "period", period)

	emailImportState, err := repo.GetEmailImportState(ctx, tenantName, provider, username, period)
	if emailImportState == nil {
		return nil
	}

	emailImportState.Active = false
	err = repo.gormDb.Save(&emailImportState).Error
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("DeactivateEmailImportState - update: %s", err.Error())
	}

	return nil
}
