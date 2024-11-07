package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
	result := r.db.WithContext(ctx).Where("email = ?", cacheEmailValidation.Email).Order("created_at desc").First(&existingData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Record doesn't exist, create a new one
			now := utils.Now()
			cacheEmailValidation.CreatedAt = now
			cacheEmailValidation.UpdatedAt = now
			if err := r.db.WithContext(ctx).Save(&cacheEmailValidation).Error; err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to save cache email validation"))
				return nil, err
			}
		} else {
			// Some other error occurred
			tracing.TraceErr(span, result.Error)
			return nil, result.Error
		}
	} else {
		// Record exists, update it
		updates := map[string]interface{}{
			"updated_at":            utils.Now(),
			"deliverable":           cacheEmailValidation.Deliverable,
			"is_mailbox_full":       cacheEmailValidation.IsMailboxFull,
			"is_role_account":       cacheEmailValidation.IsRoleAccount,
			"is_system_generated":   cacheEmailValidation.IsSystemGenerated,
			"is_free_account":       cacheEmailValidation.IsFreeAccount,
			"smtp_success":          cacheEmailValidation.SmtpSuccess,
			"response_code":         cacheEmailValidation.ResponseCode,
			"error_code":            cacheEmailValidation.ErrorCode,
			"description":           cacheEmailValidation.Description,
			"retry_validation":      cacheEmailValidation.RetryValidation,
			"normalized_email":      cacheEmailValidation.NormalizedEmail,
			"username":              cacheEmailValidation.Username,
			"domain":                cacheEmailValidation.Domain,
			"tls_required":          cacheEmailValidation.TLSRequired,
			"health_is_greylisted":  cacheEmailValidation.HealthIsGreylisted,
			"health_is_blacklisted": cacheEmailValidation.HealthIsBlacklisted,
			"health_server_ip":      cacheEmailValidation.HealthServerIP,
			"health_from_email":     cacheEmailValidation.HealthFromEmail,
			"health_retry_after":    cacheEmailValidation.HealthRetryAfter,
			"alternate_email":       cacheEmailValidation.AlternateEmail,
			"error":                 cacheEmailValidation.Error,
			"data":                  cacheEmailValidation.Data,
		}
		if err := r.db.WithContext(ctx).Model(&existingData).Updates(updates).Error; err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to update cache email validation"))
			return nil, err
		}
		cacheEmailValidation = existingData
	}

	return &cacheEmailValidation, nil
}
