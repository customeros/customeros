package postgres_repository

import (
	"context"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type CacheIPIdentifyRepository interface {
	Create(ctx context.Context, snitcherData postgres_entity.CacheIPIdentify) error
	FindByIP(ctx context.Context, ip string, cacheLookbackInDays int) (*postgres_entity.CacheIPIdentify, error)
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
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "ipData", ipData)

	if ipData.IPAddress == "" {
		err := errors.New("IP address missing")
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

func (r *cacheIPIdentifyRepository) FindByIP(ctx context.Context, ip string, cacheLookBackInDays int) (*postgres_entity.CacheIPIdentify, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CacheSnitcherRepository.FindByIP")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogKV("ip", ip, "cacheLookBackInDays", cacheLookBackInDays)

	lookBackDate := time.Now().AddDate(0, 0, -cacheLookBackInDays)

	var result postgres_entity.CacheIPIdentify
	err := r.db.
		Where("ip_address = ?", ip).
		Where("created_at > ?", lookBackDate).
		Order("created_at DESC").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogKV("result.found", false)
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogKV("result.found", true)
	return &result, nil
}
