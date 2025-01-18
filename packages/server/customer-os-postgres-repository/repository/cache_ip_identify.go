package postgres_repository

import (
	"context"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type CacheIPIdentifyRepository interface {
	Create(ctx context.Context, snitcherData postgres_entity.CacheIPIdentify) error
	Find(ctx context.Context, snitcherData postgres_entity.CacheIPIdentify, cacheLookbackInDays int) (*postgres_entity.CacheIPIdentify, error)
}

type cacheIPIdentifyRepository struct {
	db *gorm.DB
}

func NewCacheIPIdentifyRepository(gormDb *gorm.DB) CacheIPIdentifyRepository {
	return &cacheIPIdentifyRepository{db: gormDb}
}

func (r *cacheIPIdentifyRepository) Create(ctx context.Context, ipData postgres_entity.CacheIPIdentify) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CacheSnitcherRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if ipData.IPAddress == "" || ipData.Domain == "" {
		span.LogFields(log.Object("ipData", ipData))
		err := errors.New("IP address or domain missing")
		tracing.TraceErr(span, err)
		return err
	}

	err := r.db.Create(&ipData).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *cacheIPIdentifyRepository) Find(ctx context.Context, ipData postgres_entity.CacheIPIdentify, cacheLookbackInDays int) (*postgres_entity.CacheIPIdentify, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CacheSnitcherRepository.Find")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	lookbackDate := time.Now().AddDate(0, 0, -cacheLookbackInDays)

	var result postgres_entity.CacheIPIdentify
	err := r.db.
		Where(&ipData).
		Where("created_at > ?", lookbackDate).
		Order("created_at DESC").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &result, nil
}
