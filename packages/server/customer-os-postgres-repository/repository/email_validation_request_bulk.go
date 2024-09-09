package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type EmailValidationRequestBulkRepository interface {
	RegisterRequest(ctx context.Context, tenant, requestId, fileName string, verifyCatchAll bool, totalRecords int) (*entity.EmailValidationRequestBulk, error)
	GetByRequestID(ctx context.Context, requestId string) (*entity.EmailValidationRequestBulk, error)
	IncrementDeliverableEmails(ctx context.Context, requestID string) error
	IncrementUndeliverableEmails(ctx context.Context, requestID string) error
	MarkRequestAsCompleted(ctx context.Context, requestID string) error
}

type emailValidationRequestBulkRepository struct {
	db *gorm.DB
}

func NewEmailValidationRequestBulkRepository(gormDb *gorm.DB) EmailValidationRequestBulkRepository {
	return &emailValidationRequestBulkRepository{db: gormDb}
}

func (r emailValidationRequestBulkRepository) RegisterRequest(ctx context.Context, tenant, requestId, fileName string, verifyCatchAll bool, totalRecords int) (*entity.EmailValidationRequestBulk, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRequestBulkRepository.RegisterRequest")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(
		log.String("requestId", requestId),
		log.String("fileName", fileName),
		log.Int("totalRecords", totalRecords),
		log.Bool("verifyCatchAll", verifyCatchAll))

	// Create a new EmailValidationRequestBulk record
	record := entity.EmailValidationRequestBulk{
		RequestID:      requestId,
		Tenant:         tenant,
		FileName:       fileName,
		Status:         entity.EmailValidationRequestBulkStatusProcessing, // Initial status
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

func (r emailValidationRequestBulkRepository) GetByRequestID(ctx context.Context, requestID string) (*entity.EmailValidationRequestBulk, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRequestBulkRepository.GetByRequestID")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("requestID", requestID))

	var record entity.EmailValidationRequestBulk

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
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRequestBulkRepository.IncrementDeliverableEmails")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("requestID", requestID))

	// Increment deliverable emails count and update the updated_at timestamp
	if err := r.db.WithContext(ctx).
		Model(&entity.EmailValidationRequestBulk{}).
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
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRequestBulkRepository.IncrementUndeliverableEmails")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("requestID", requestID))

	// Increment undeliverable emails count
	if err := r.db.WithContext(ctx).
		Model(&entity.EmailValidationRequestBulk{}).
		Where("request_id = ?", requestID).
		UpdateColumns(map[string]interface{}{
			"deliverable_emails": gorm.Expr("undeliverable_emails + ?", 1),
			"updated_at":         utils.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r emailValidationRequestBulkRepository) MarkRequestAsCompleted(ctx context.Context, requestID string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailValidationRequestBulkRepository.MarkRequestAsCompleted")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("requestID", requestID)

	// Update the status to "completed" and set the updated_at field to the current time
	if err := r.db.WithContext(ctx).
		Model(&entity.EmailValidationRequestBulk{}).
		Where("request_id = ?", requestID).
		Updates(map[string]interface{}{
			"status":     entity.EmailValidationRequestBulkStatusCompleted,
			"updated_at": utils.Now(),
		}).Error; err != nil {
		return err
	}

	return nil
}
