package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type EmailValidationRequestBulkRepository interface {
	RegisterRequest(ctx context.Context, tenant, requestId, fileName string, verifyCatchAll bool, totalRecords int) (*postgres_entity.EmailValidationRequestBulk, error)
	GetByRequestID(ctx context.Context, requestId string) (*postgres_entity.EmailValidationRequestBulk, error)
	IncrementDeliverableEmails(ctx context.Context, requestID string) error
	IncrementUndeliverableEmails(ctx context.Context, requestID string) error
	MarkRequestAsCompleted(ctx context.Context, requestID, fileStoreId string) error
	GetOldestUncompletedRequests(ctx context.Context, limit int) ([]postgres_entity.EmailValidationRequestBulk, error)
}

type emailValidationRequestBulkRepository struct {
	db *gorm.DB
}

func NewEmailValidationRequestBulkRepository(gormDb *gorm.DB) EmailValidationRequestBulkRepository {
	return &emailValidationRequestBulkRepository{db: gormDb}
}

func (r emailValidationRequestBulkRepository) RegisterRequest(ctx context.Context, tenant, requestId, fileName string, verifyCatchAll bool, totalRecords int) (*postgres_entity.EmailValidationRequestBulk, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailValidationRequestBulkRepository.RegisterRequest")
	defer spans.Finish()

	spans.LogKV("requestId", requestId)
	spans.LogKV("fileName", fileName)
	spans.LogKV("totalRecords", totalRecords)
	spans.LogKV("verifyCatchAll", verifyCatchAll)

	// Create a new EmailValidationRequestBulk record
	record := postgres_entity.EmailValidationRequestBulk{
		RequestID:      requestId,
		Tenant:         tenant,
		FileName:       fileName,
		Status:         postgres_entity.EmailValidationRequestBulkStatusProcessing, // Initial status
		TotalEmails:    totalRecords,
		CreatedAt:      utils.Now(),
		Priority:       assignPriority(totalRecords),
		VerifyCatchAll: verifyCatchAll,
	}

	// Insert the new record into the database
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (r emailValidationRequestBulkRepository) GetByRequestID(ctx context.Context, requestID string) (*postgres_entity.EmailValidationRequestBulk, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailValidationRequestBulkRepository.GetByRequestID")
	defer spans.Finish()

	spans.LogKV("requestID", requestID)

	var record postgres_entity.EmailValidationRequestBulk

	// Query the database for the record with the given request ID
	if err := r.db.WithContext(ctx).Where("request_id = ?", requestID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &record, nil
}

func (r emailValidationRequestBulkRepository) IncrementDeliverableEmails(ctx context.Context, requestID string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailValidationRequestBulkRepository.IncrementDeliverableEmails")
	defer spans.Finish()

	spans.LogKV("requestID", requestID)

	// Increment deliverable emails count and update the updated_at timestamp
	if err := r.db.WithContext(ctx).
		Model(&postgres_entity.EmailValidationRequestBulk{}).
		Where("request_id = ?", requestID).
		UpdateColumns(map[string]interface{}{
			"deliverable_emails": gorm.Expr("deliverable_emails + ?", 1),
			"updated_at":         utils.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r emailValidationRequestBulkRepository) IncrementUndeliverableEmails(ctx context.Context, requestID string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailValidationRequestBulkRepository.IncrementUndeliverableEmails")
	defer spans.Finish()

	spans.LogKV("requestID", requestID)

	// Increment undeliverable emails count
	if err := r.db.WithContext(ctx).
		Model(&postgres_entity.EmailValidationRequestBulk{}).
		Where("request_id = ?", requestID).
		UpdateColumns(map[string]interface{}{
			"undeliverable_emails": gorm.Expr("undeliverable_emails + ?", 1),
			"updated_at":           utils.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r emailValidationRequestBulkRepository) MarkRequestAsCompleted(ctx context.Context, requestId, fileStoreId string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailValidationRequestBulkRepository.MarkRequestAsCompleted")
	defer spans.Finish()

	spans.LogKV("requestId", requestId)
	spans.LogKV("fileStoreId", fileStoreId)

	// Update the status to "completed" and set the updated_at field to the current time
	if err := r.db.WithContext(ctx).
		Model(&postgres_entity.EmailValidationRequestBulk{}).
		Where("request_id = ?", requestId).
		Updates(map[string]interface{}{
			"file_store_id": fileStoreId,
			"status":        postgres_entity.EmailValidationRequestBulkStatusCompleted,
			"updated_at":    utils.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r emailValidationRequestBulkRepository) GetOldestUncompletedRequests(ctx context.Context, limit int) ([]postgres_entity.EmailValidationRequestBulk, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EmailValidationRequestBulkRepository.GetOldestUncompletedRequests")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	var records []postgres_entity.EmailValidationRequestBulk

	// Query the database for the oldest uncompleted requests
	if err := r.db.WithContext(ctx).
		Where("status = ?", postgres_entity.EmailValidationRequestBulkStatusProcessing).
		Order("created_at ASC").
		Limit(limit).
		Find(&records).Error; err != nil {
		return nil, err
	}

	spans.LogKV("result.count", len(records))
	return records, nil
}
