package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

const chunkSize = 1000

type EmailValidationRecordRepository interface {
	BulkInsertRecords(ctx context.Context, tenant, requestId string, verifyCatchAll bool, emails []string) error
	UpdateEmailRecord(ctx context.Context, id uint64, newData string) error
	CountPendingRequests(ctx context.Context, priority int, createdBefore time.Time) (int64, error)
	CountPendingRequestsByRequestID(ctx context.Context, requestID string) (int64, error)
	GetUnprocessedEmailRecords(ctx context.Context, limit int) ([]entity.EmailValidationRecord, error)
	GetEmailRecordsInChunks(ctx context.Context, chunkSize int, offset int) ([]entity.EmailValidationRecord, error)
}

type emailValidationRecordRepository struct {
	db *gorm.DB
}

func NewEmailValidationRecordRepository(gormDb *gorm.DB) EmailValidationRecordRepository {
	return &emailValidationRecordRepository{db: gormDb}
}

func (r emailValidationRecordRepository) BulkInsertRecords(ctx context.Context, tenant, requestId string, verifyCatchAll bool, emails []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRecordRepository.BulkInsertRecords")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(
		log.String("requestId", requestId),
		log.Int("emailsCount", len(emails)),
		log.Bool("verifyCatchAll", verifyCatchAll))

	priority := assignPriority(len(emails))

	// Loop through emails and insert them in chunks
	for i := 0; i < len(emails); i += chunkSize {
		end := i + chunkSize
		if end > len(emails) {
			end = len(emails)
		}

		// Prepare the batch for the current chunk
		var records []entity.EmailValidationRecord
		for j := i; j < end; j++ {
			records = append(records, entity.EmailValidationRecord{
				RequestID:      requestId,
				Tenant:         tenant,
				Email:          emails[j],
				CreatedAt:      utils.Now(),
				Priority:       priority,
				Data:           "",
				VerifyCatchAll: verifyCatchAll,
			})
		}

		// Insert the chunk into the database
		if err := r.db.WithContext(ctx).Create(&records).Error; err != nil {
			// If any error occurs, return it (could enhance to log and proceed with next chunk)
			return err
		}
	}

	return nil
}

func assignPriority(recordCount int) int {
	switch {
	case recordCount <= 10:
		return 0 // Highest priority for small files
	case recordCount <= 50:
		return 1
	case recordCount <= 100:
		return 2
	case recordCount <= 500:
		return 3
	case recordCount <= 1000:
		return 4
	case recordCount <= 5000:
		return 5
	case recordCount <= 10000:
		return 6
	default:
		return 7
	}
}

func (r emailValidationRecordRepository) UpdateEmailRecord(ctx context.Context, id uint64, newData string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRecordRepository.UpdateEmailRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.Uint64("id", id), log.String("newData", newData))

	// Find the record by email and request_id and update the data field
	if err := r.db.WithContext(ctx).
		Model(&entity.EmailValidationRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"data":       newData,
			"updated_at": utils.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r emailValidationRecordRepository) CountPendingRequests(ctx context.Context, priority int, createdBefore time.Time) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRequestBulkRepository.CountPendingRequests")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.Int("priority", priority), log.Object("createdBefore", createdBefore))

	var count int64

	// Count records with same or lower priority, created on or before with given date, and where data is empty
	if err := r.db.WithContext(ctx).
		Model(&entity.EmailValidationRecord{}).
		Where("priority <= ? AND created_at <= ? AND data = ''", priority, createdBefore).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r emailValidationRecordRepository) GetUnprocessedEmailRecords(ctx context.Context, limit int) ([]entity.EmailValidationRecord, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRecordRepository.GetUnprocessedEmailRecords")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.Int("limit", limit))

	var records []entity.EmailValidationRecord

	// Query to return unprocessed email records ordered by priority and creation date
	if err := r.db.WithContext(ctx).
		Where("data = ''").      // Unprocessed records where data is empty
		Order("priority ASC").   // Order by priority (lowest first)
		Order("created_at ASC"). // Order by created date (oldest first)
		Limit(limit).            // Limit the number of records returned
		Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r emailValidationRecordRepository) CountPendingRequestsByRequestID(ctx context.Context, requestID string) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRecordRepository.CountPendingRequestsByRequestID")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var count int64

	// Count records where data is empty (unprocessed) for the specific requestID
	if err := r.db.WithContext(ctx).
		Model(&entity.EmailValidationRecord{}).
		Where("request_id = ? AND data = ''", requestID). // Filter for unprocessed records with the given requestID
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r emailValidationRecordRepository) GetEmailRecordsInChunks(ctx context.Context, chunkSize int, offset int) ([]entity.EmailValidationRecord, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRecordRepository.GetEmailRecordsInChunks")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var records []entity.EmailValidationRecord

	// Retrieve records in chunks
	if err := r.db.WithContext(ctx).
		Order("created_at ASC, id ASC"). // Ensure consistent ordering by created_at and id
		Limit(chunkSize).
		Offset(offset).
		Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}
