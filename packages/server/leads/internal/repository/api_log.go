package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/customeros/customeros/packages/server/leads/internal/database"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

type APICallLogRepository interface {
	Create(ctx context.Context, log *models.APICallLog) error
	GetByID(ctx context.Context, id string) (*models.APICallLog, error)
	FindByRequestID(ctx context.Context, requestID string) (*models.APICallLog, error)
	FindByVendor(ctx context.Context, vendor enum.APIVendor, limit, offset int) ([]*models.APICallLog, error)
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

func (r *apiCallLogRepository) Create(ctx context.Context, log *models.APICallLog) error {
	span, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.Create")
	defer span.Finish()

	// Generate ID if not provided
	if log.ID == "" {
		log.ID = utils.GenerateNanoIDWithPrefix("api", 21)
	}

	return r.write.WithContext(ctx).Create(log).Error
}

func (r *apiCallLogRepository) GetByID(ctx context.Context, id string) (*models.APICallLog, error) {
	span, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.GetByID")
	defer span.Finish()

	var log models.APICallLog
	err := r.read.WithContext(ctx).Where("id = ?", id).First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("API call log not found with ID: %s", id)
		}
		return nil, err
	}

	return &log, nil
}

func (r *apiCallLogRepository) FindByRequestID(ctx context.Context, requestID string) (*models.APICallLog, error) {
	span, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.FindByRequestID")
	defer span.Finish()

	var log models.APICallLog
	err := r.read.WithContext(ctx).Where("request_id = ?", requestID).First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil when not found
		}
		return nil, err
	}

	return &log, nil
}

func (r *apiCallLogRepository) FindByVendor(ctx context.Context, vendor enum.APIVendor, limit, offset int) ([]*models.APICallLog, error) {
	span, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.FindByVendor")
	defer span.Finish()

	var logs []*models.APICallLog
	err := r.read.WithContext(ctx).
		Where("vendor = ?", vendor).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *apiCallLogRepository) UpdateResponseData(ctx context.Context, id string, statusCode int, responseBody []byte, errorMessage *string) error {
	span, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.UpdateResponseData")
	defer span.Finish()

	updates := map[string]interface{}{
		"status_code": statusCode,
	}

	if responseBody != nil {
		updates["response_body"] = responseBody
	}

	if errorMessage != nil {
		updates["error_message"] = errorMessage
	}

	return r.write.WithContext(ctx).
		Model(&models.APICallLog{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *apiCallLogRepository) DeleteOlderThan(ctx context.Context, age time.Duration) (int64, error) {
	span, ctx := telemetry.StartPostgresSpan(ctx, "apiCallLogRepository.DeleteOlderThan")
	defer span.Finish()

	cutoffTime := time.Now().Add(-age)
	result := r.write.WithContext(ctx).
		Where("timestamp < ?", cutoffTime).
		Delete(&models.APICallLog{})

	return result.RowsAffected, result.Error
}

// SetupTimescaleDB initializes the TimescaleDB specifics for API call logs
func SetupTimescaleDB(db *gorm.DB) error {
	// Migrate the schema
	if err := db.AutoMigrate(&models.APICallLog{}); err != nil {
		return err
	}

	// Convert to hypertable
	if err := db.Exec(`SELECT create_hypertable('api_call_logs', 'timestamp', 
   	chunk_time_interval => INTERVAL '1 day',
   	if_not_exists => TRUE)`).Error; err != nil {
		return err
	}

	// Create indexes for common query patterns
	if err := db.Exec(`
   	CREATE INDEX IF NOT EXISTS idx_api_call_logs_vendor_timestamp ON api_call_logs (vendor, timestamp DESC);
   	CREATE INDEX IF NOT EXISTS idx_api_call_logs_request_id ON api_call_logs (request_id);
   	CREATE INDEX IF NOT EXISTS idx_api_call_logs_status_code ON api_call_logs (status_code);
   `).Error; err != nil {
		return err
	}

	return nil
}
