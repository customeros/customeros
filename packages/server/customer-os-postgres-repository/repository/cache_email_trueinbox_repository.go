package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type CacheEmailTrueinboxRepository interface {
	Create(ctx context.Context, record entity.CacheEmailTrueinbox) (*entity.CacheEmailTrueinbox, error)
	GetLatestByEmail(ctx context.Context, email string) (*entity.CacheEmailTrueinbox, error)
	GetAllByEmail(ctx context.Context, email string) ([]entity.CacheEmailTrueinbox, error)
}

type cacheEmailTrueinboxRepository struct {
	db *gorm.DB
}

func NewCacheEmailTrueinboxRepository(gormDb *gorm.DB) CacheEmailTrueinboxRepository {
	return &cacheEmailTrueinboxRepository{db: gormDb}
}

func (r cacheEmailTrueinboxRepository) GetAllByEmail(ctx context.Context, email string) ([]entity.CacheEmailTrueinbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailTrueinboxRepository.GetAllByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("email", email))

	var records []entity.CacheEmailTrueinbox
	err := r.db.Where("email = ?", email).Order("created_at desc").Find(&records).Error
	if err != nil {
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(records)))

	return records, nil
}

func (r cacheEmailTrueinboxRepository) GetLatestByEmail(ctx context.Context, email string) (*entity.CacheEmailTrueinbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailTrueinboxRepository.GetLatestByEmail")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.String("email", email))

	var data entity.CacheEmailTrueinbox
	err := r.db.Where("email = ?", email).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r cacheEmailTrueinboxRepository) Create(ctx context.Context, record entity.CacheEmailTrueinbox) (*entity.CacheEmailTrueinbox, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CacheEmailTrueinboxRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "record", record)

	record.CreatedAt = utils.Now()
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}
