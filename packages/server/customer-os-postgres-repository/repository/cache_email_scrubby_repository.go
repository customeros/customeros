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
	"time"
)

type CacheEmailScrubbyRepository interface {
	Save(ctx context.Context, cacheEmailScrubby entity.CacheEmailScrubby) (*entity.CacheEmailScrubby, error)
	GetAllByEmail(ctx context.Context, email string) ([]entity.CacheEmailScrubby, error)
	GetLatestByEmail(ctx context.Context, email string) (*entity.CacheEmailScrubby, error)
	SetStatus(ctx context.Context, email, status string) error
	SetJustChecked(ctx context.Context, id string) (*entity.CacheEmailScrubby, error)
	GetToCheck(ctx context.Context, delayFromPreviousCheckInHours, limit int) ([]entity.CacheEmailScrubby, error)
}

type cacheEmailScrubbyRepository struct {
	db *gorm.DB
}

func NewCacheEmailScrubbyRepository(gormDb *gorm.DB) CacheEmailScrubbyRepository {
	return &cacheEmailScrubbyRepository{db: gormDb}
}

func (r *cacheEmailScrubbyRepository) Save(ctx context.Context, cacheEmailScrubby entity.CacheEmailScrubby) (*entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.Save")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "cacheEmailScrubby", cacheEmailScrubby)

	now := utils.Now()
	cacheEmailScrubby.CreatedAt = now
	cacheEmailScrubby.CheckedAt = now

	if err := r.db.WithContext(ctx).Save(&cacheEmailScrubby).Error; err != nil {
		return nil, err
	}

	return &cacheEmailScrubby, nil
}

func (r *cacheEmailScrubbyRepository) GetAllByEmail(ctx context.Context, email string) ([]entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.GetAllByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("email", email))

	var cacheEmailScrubbys []entity.CacheEmailScrubby
	result := r.db.WithContext(ctx).Where("email = ?", email).Order("created_at desc").Find(&cacheEmailScrubbys)

	if result.Error != nil {
		return nil, result.Error
	}

	return cacheEmailScrubbys, nil
}

func (r *cacheEmailScrubbyRepository) GetLatestByEmail(ctx context.Context, email string) (*entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.GetLatestByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("email", email))

	var cacheEmailScrubby entity.CacheEmailScrubby
	result := r.db.WithContext(ctx).Where("email = ?", email).Order("created_at desc").First(&cacheEmailScrubby)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &cacheEmailScrubby, nil
}

func (r *cacheEmailScrubbyRepository) SetStatus(ctx context.Context, email, status string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.SetStatus")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("email", email), log.String("status", status))

	result := r.db.WithContext(ctx).
		Model(&entity.CacheEmailScrubby{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"status":     status,
			"checked_at": utils.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	// If no records were affected, it's not considered an error
	if result.RowsAffected == 0 {
		span.LogFields(log.String("message", "No records found for the given email"))
	}

	return nil
}

func (r *cacheEmailScrubbyRepository) SetJustChecked(ctx context.Context, id string) (*entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.SetJustChecked")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("id", id))

	var cacheEmailScrubby entity.CacheEmailScrubby
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&cacheEmailScrubby)

	if result.Error != nil {
		return nil, result.Error
	}

	cacheEmailScrubby.CheckedAt = utils.Now()

	if err := r.db.WithContext(ctx).Save(&cacheEmailScrubby).Error; err != nil {
		return nil, err
	}

	return &cacheEmailScrubby, nil
}

func (r *cacheEmailScrubbyRepository) GetToCheck(ctx context.Context, delayFromPreviousCheckInHours, limit int) ([]entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.GetToCheck")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.Int("delayFromPreviousCheckInHours", delayFromPreviousCheckInHours), log.Int("limit", limit))

	var cacheEmailScrubbys []entity.CacheEmailScrubby

	// Calculate the cutoff time
	cutoffTime := time.Now().Add(-time.Duration(delayFromPreviousCheckInHours) * time.Hour)

	result := r.db.WithContext(ctx).
		Where("status = ? AND checked_at < ?", string(entity.ScrubbyStatusPending), cutoffTime).
		Order("checked_at asc").
		Limit(limit).
		Find(&cacheEmailScrubbys)

	if result.Error != nil {
		return nil, result.Error
	}

	return cacheEmailScrubbys, nil
}
