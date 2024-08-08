package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
	"time"
)

type StatsApiCallsRepository interface {
	Increment(ctx context.Context, tenant, api string) (*entity.StatsApiCalls, error)
}

type statsApiCallsRepository struct {
	db *gorm.DB
}

func NewStatsApiCallsRepository(gormDb *gorm.DB) StatsApiCallsRepository {
	return &statsApiCallsRepository{db: gormDb}
}

func (r statsApiCallsRepository) Increment(ctx context.Context, tenant, api string) (*entity.StatsApiCalls, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "StatsApiCallsRepository.Increment")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("api", api))

	// Check if the record already exists
	var stats entity.StatsApiCalls
	err := r.db.WithContext(ctx).
		Where("tenant = ? AND api = ? AND day = ?", tenant, api, utils.Today()).
		First(&stats).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new record
			newStats := &entity.StatsApiCalls{
				Tenant: tenant,
				Api:    api,
				Day:    utils.Today(),
				Calls:  1,
			}
			err = r.db.WithContext(ctx).Create(newStats).Error
			return newStats, err
		}
		return nil, err
	}

	// Increment the call count and update the record
	stats.Calls++
	stats.UpdatedAt = time.Now()
	err = r.db.WithContext(ctx).Save(&stats).Error
	return &stats, err
}
