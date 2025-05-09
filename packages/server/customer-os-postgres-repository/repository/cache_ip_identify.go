package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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
	spans, ctx := telemetry.StartPostgresSpan(ctx, "cacheIPIdentifyRepository.Create")
	defer spans.Finish()
	spans.LogObjectAsJson("ipData", ipData)

	if ipData.IPAddress == "" {
		err := errors.New("IP address missing")
		spans.TraceError(err)
		return err
	}

	err := r.db.Create(&ipData).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *cacheIPIdentifyRepository) FindByIP(ctx context.Context, ip string, cacheLookBackInDays int) (*postgres_entity.CacheIPIdentify, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "cacheIPIdentifyRepository.FindByIP")
	defer spans.Finish()
	spans.LogKV("ip", ip, "cacheLookBackInDays", cacheLookBackInDays)

	lookBackDate := utils.Now().AddDate(0, 0, -cacheLookBackInDays)

	var result postgres_entity.CacheIPIdentify
	err := r.db.
		Where("ip_address = ?", ip).
		Where("created_at > ?", lookBackDate).
		Order("created_at DESC").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", true)
	return &result, nil
}
