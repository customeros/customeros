package postgres_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/database"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"time"

	"gorm.io/gorm"
)

type APICallLogRepository interface {
	Create(ctx context.Context, log *postgres_entity.APICallLog) error
	GetByID(ctx context.Context, id string) (*postgres_entity.APICallLog, error)
	FindByRequestID(ctx context.Context, requestID string) (*postgres_entity.APICallLog, error)
	FindByVendor(ctx context.Context, vendor enum.APIVendor, limit, offset int) ([]*postgres_entity.APICallLog, error)
	UpdateResponseData(ctx context.Context, id string, statusCode int, responseBody []byte, errorMessage *string) error
	DeleteOlderThan(ctx context.Context, age time.Duration) (int64, error)
}

type apiCallLogRepository struct {
	read  *gorm.DB
	write *gorm.DB
}

func NewAPICallLogRepository(db *database.DbConnections) APICallLogRepository {
	return &apiCallLogRepository{
		read:  db.ReadDB,
		write: db.WriteDB,
	}
}

func (r *apiCallLogRepository) Create(ctx context.Context, log *postgres_entity.APICallLog) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.Create")
	defer spans.Finish()

	// Generate ID if not provided
	if log.ID == "" {
		log.ID = utils.GenerateNanoIdWithPrefix("api", 21)
	}

	err := r.write.WithContext(ctx).Create(log).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (r *apiCallLogRepository) GetByID(ctx context.Context, id string) (*postgres_entity.APICallLog, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.GetByID")
	defer spans.Finish()

	var log postgres_entity.APICallLog
	err := r.read.WithContext(ctx).Where("id = ?", id).First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("API call log not found with ID: %s", id)
		}
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &log, nil
}

func (r *apiCallLogRepository) FindByRequestID(ctx context.Context, requestID string) (*postgres_entity.APICallLog, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.FindByRequestID")
	defer spans.Finish()

	var log postgres_entity.APICallLog
	err := r.read.WithContext(ctx).Where("request_id = ?", requestID).First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			spans.LogKV("result.found", false)
			return nil, nil // Return nil, nil when not found
		}
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return &log, nil
}

func (r *apiCallLogRepository) FindByVendor(ctx context.Context, vendor enum.APIVendor, limit, offset int) ([]*postgres_entity.APICallLog, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.FindByVendor")
	defer spans.Finish()

	var logs []*postgres_entity.APICallLog
	err := r.read.WithContext(ctx).
		Where("vendor = ?", vendor).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return logs, nil
}

func (r *apiCallLogRepository) UpdateResponseData(ctx context.Context, id string, statusCode int, responseBody []byte, errorMessage *string) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.UpdateResponseData")
	defer spans.Finish()

	updates := map[string]interface{}{
		"status_code": statusCode,
	}

	if responseBody != nil {
		updates["response_body"] = responseBody
	}

	if errorMessage != nil {
		updates["error_message"] = errorMessage
	}

	err := r.write.WithContext(ctx).
		Model(&postgres_entity.APICallLog{}).
		Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *apiCallLogRepository) DeleteOlderThan(ctx context.Context, age time.Duration) (int64, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.DeleteOlderThan")
	defer spans.Finish()

	cutoffTime := time.Now().Add(-age)
	result := r.write.WithContext(ctx).
		Where("timestamp < ?", cutoffTime).
		Delete(&postgres_entity.APICallLog{})

	return result.RowsAffected, result.Error
}
