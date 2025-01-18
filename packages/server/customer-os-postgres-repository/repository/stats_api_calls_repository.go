package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type StatsApiCallsRepository interface {
	Increment(ctx context.Context, tenant, api string) (*postgres_entity.StatsApiCalls, error)
}

type statsApiCallsRepository struct {
	db *gorm.DB
}

func NewStatsApiCallsRepository(gormDb *gorm.DB) StatsApiCallsRepository {
	return &statsApiCallsRepository{db: gormDb}
}

func (r statsApiCallsRepository) Increment(ctx context.Context, tenant, api string) (*postgres_entity.StatsApiCalls, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "StatsApiCallsRepository.Increment")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("api", api))

	// Check if the record already exists
	var stats postgres_entity.StatsApiCalls
	err := r.db.WithContext(ctx).
		Where("tenant = ? AND api = ? AND day = ?", tenant, api, utils.Today()).
		First(&stats).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new record
			newStats := &postgres_entity.StatsApiCalls{
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
	stats.UpdatedAt = utils.Now()
	err = r.db.WithContext(ctx).Save(&stats).Error
	return &stats, err
}
