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
	Add(ctx context.Context, email string) (*entity.CacheEmailScrubby, error)
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

func (c cacheEmailScrubbyRepository) Add(ctx context.Context, email string) (*entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.Add")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("email", email))

	cacheEmailScrubby := entity.CacheEmailScrubby{
		Email: email,
	}
	if err := c.db.WithContext(ctx).Create(&cacheEmailScrubby).Error; err != nil {
		return nil, err
	}

	return &cacheEmailScrubby, nil
}

func (c cacheEmailScrubbyRepository) GetAllByEmail(ctx context.Context, email string) ([]entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.GetAllByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("email", email))

	var cacheEmailScrubbys []entity.CacheEmailScrubby
	result := c.db.WithContext(ctx).Where("email = ?", email).Find(&cacheEmailScrubbys)

	if result.Error != nil {
		return nil, result.Error
	}

	return cacheEmailScrubbys, nil
}

func (c cacheEmailScrubbyRepository) GetLatestByEmail(ctx context.Context, email string) (*entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.GetLatestByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("email", email))

	var cacheEmailScrubby entity.CacheEmailScrubby
	result := c.db.WithContext(ctx).Where("email = ?", email).Order("created_at desc").First(&cacheEmailScrubby)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &cacheEmailScrubby, nil
}

func (c cacheEmailScrubbyRepository) SetStatus(ctx context.Context, id, status string) (*entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.SetStatus")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("id", id), log.String("status", status))

	var cacheEmailScrubby entity.CacheEmailScrubby
	result := c.db.WithContext(ctx).Where("id = ?", id).First(&cacheEmailScrubby)

	if result.Error != nil {
		return nil, result.Error
	}

	cacheEmailScrubby.Status = status
	cacheEmailScrubby.CheckedAt = utils.Now()

	if err := c.db.WithContext(ctx).Save(&cacheEmailScrubby).Error; err != nil {
		return nil, err
	}

	return &cacheEmailScrubby, nil
}

func (c cacheEmailScrubbyRepository) SetJustChecked(ctx context.Context, id string) (*entity.CacheEmailScrubby, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailScrubbyRepository.SetJustChecked")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("id", id))

	var cacheEmailScrubby entity.CacheEmailScrubby
	result := c.db.WithContext(ctx).Where("id = ?", id).First(&cacheEmailScrubby)

	if result.Error != nil {
		return nil, result.Error
	}

	cacheEmailScrubby.CheckedAt = utils.Now()

	if err := c.db.WithContext(ctx).Save(&cacheEmailScrubby).Error; err != nil {
		return nil, err
	}

	return &cacheEmailScrubby, nil
}
