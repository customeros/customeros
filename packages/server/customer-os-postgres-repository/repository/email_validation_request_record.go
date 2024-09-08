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
	BulkInsertRecords(ctx context.Context, tenant, requestId string, emails []string) error
	UpdateEmailRecord(ctx context.Context, requestId, email, newData string) error
	CountPendingRequests(ctx context.Context, priority int, createdBefore time.Time) (int64, error)
}

type emailValidationRecordRepository struct {
	db *gorm.DB
}

func NewEmailValidationRecordRepository(gormDb *gorm.DB) EmailValidationRecordRepository {
	return &emailValidationRecordRepository{db: gormDb}
}

func (r emailValidationRecordRepository) BulkInsertRecords(ctx context.Context, tenant, requestId string, emails []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRecordRepository.BulkInsertRecords")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("requestId", requestId), log.Int("emailsCount", len(emails)))

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
				RequestID: requestId,
				Tenant:    tenant,
				Email:     emails[j],
				CreatedAt: utils.Now(),
				Priority:  priority,
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

func (r emailValidationRecordRepository) UpdateEmailRecord(ctx context.Context, requestId, email, newData string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRecordRepository.UpdateEmailRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("requestId", requestId)
	span.LogKV("email", email)

	// Find the record by email and request_id and update the data field
	if err := r.db.WithContext(ctx).
		Model(&entity.EmailValidationRecord{}).
		Where("request_id = ? AND email = ?", requestId, email).
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
