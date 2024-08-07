package repository

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type CacheEmailValidationRepository interface {
	Get(ctx context.Context, email string) (*entity.CacheEmailValidation, error)
	Save(ctx context.Context, cacheEmailValidation entity.CacheEmailValidation) (*entity.CacheEmailValidation, error)
}

type cacheEmailValidationRepository struct {
	db *gorm.DB
}

func NewCacheEmailValidationRepository(gormDb *gorm.DB) CacheEmailValidationRepository {
	return &cacheEmailValidationRepository{db: gormDb}
}

func (r cacheEmailValidationRepository) Get(ctx context.Context, email string) (*entity.CacheEmailValidation, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailValidationRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("email", email))

	var cacheEmailValidation entity.CacheEmailValidation
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&cacheEmailValidation)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &cacheEmailValidation, nil
}

func (r cacheEmailValidationRepository) Save(ctx context.Context, cacheEmailValidation entity.CacheEmailValidation) (*entity.CacheEmailValidation, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailValidationRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "cacheEmailValidation", cacheEmailValidation)

	var existingData entity.CacheEmailValidation
	result := r.db.WithContext(ctx).Where("email = ?", cacheEmailValidation.Email).First(&existingData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Record doesn't exist, create a new one
			cacheEmailValidation.CreatedAt = utils.Now()
			cacheEmailValidation.UpdatedAt = utils.Now()
			if err := r.db.WithContext(ctx).Create(&cacheEmailValidation).Error; err != nil {
				return nil, err
			}
		} else {
			// Some other error occurred
			return nil, result.Error
		}
	} else {
		// Record exists, update it
		updates := map[string]interface{}{
			"updated_at":       utils.Now(),
			"is_deliverable":   cacheEmailValidation.IsDeliverable,
			"is_mailbox_full":  cacheEmailValidation.IsMailboxFull,
			"is_role_account":  cacheEmailValidation.IsRoleAccount,
			"is_free_account":  cacheEmailValidation.IsFreeAccount,
			"smtp_success":     cacheEmailValidation.SmtpSuccess,
			"response_code":    cacheEmailValidation.ResponseCode,
			"error_code":       cacheEmailValidation.ErrorCode,
			"description":      cacheEmailValidation.Description,
			"retry_validation": cacheEmailValidation.RetryValidation,
			"smtp_response":    cacheEmailValidation.SmtpResponse,
		}
		if err := r.db.WithContext(ctx).Model(&existingData).Updates(updates).Error; err != nil {
			return nil, err
		}
		cacheEmailValidation = existingData
	}

	return &cacheEmailValidation, nil
}
