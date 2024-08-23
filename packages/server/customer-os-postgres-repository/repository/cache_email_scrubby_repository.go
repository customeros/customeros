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

type CacheEmailScrubbyRepository interface {
	Save(ctx context.Context, cacheEmailScrubby entity.CacheEmailScrubby) (*entity.CacheEmailScrubby, error)
	GetAllByEmail(ctx context.Context, email string) ([]entity.CacheEmailScrubby, error)
	GetLatestByEmail(ctx context.Context, email string) (*entity.CacheEmailScrubby, error)
	SetStatus(ctx context.Context, id, status string) (*entity.CacheEmailScrubby, error)
	SetJustChecked(ctx context.Context, id string) (*entity.CacheEmailScrubby, error)
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

func (r *cacheEmailScrubbyRepository) SetStatus(ctx context.Context, id, status string) (*entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.SetStatus")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("id", id), log.String("status", status))

	var cacheEmailScrubby entity.CacheEmailScrubby
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&cacheEmailScrubby)

	if result.Error != nil {
		return nil, result.Error
	}

	cacheEmailScrubby.Status = status
	cacheEmailScrubby.CheckedAt = utils.Now()

	if err := r.db.WithContext(ctx).Save(&cacheEmailScrubby).Error; err != nil {
		return nil, err
	}

	return &cacheEmailScrubby, nil
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
