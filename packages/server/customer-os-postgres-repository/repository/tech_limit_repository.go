package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type techLimitRepository struct {
	gormDb *gorm.DB
}

type TechLimitRepository interface {
	GetTechLimit(ctx context.Context, key string) helper.QueryResult
	IncrementTechLimit(ctx context.Context, key string) helper.QueryResult
}

func NewTechLimitRepository(gormDb *gorm.DB) TechLimitRepository {
	return &techLimitRepository{gormDb: gormDb}
}

func (t techLimitRepository) GetTechLimit(ctx context.Context, key string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TechLimitRepository.GetTechLimit")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.LogFields(log.String("key", key))

	var techLimits []entity.TechLimit
	err := t.gormDb.
		Where("key = ?", key).
		Find(&techLimits).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	if len(techLimits) == 0 {
		return helper.QueryResult{Result: entity.TechLimit{
			Key:        key,
			UsageCount: 0,
		}}
	}
	return helper.QueryResult{Result: techLimits[0]}
}

func (t techLimitRepository) IncrementTechLimit(ctx context.Context, key string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "TechLimitRepository.IncrementTechLimit")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.LogFields(log.String("key", key))

	// create entry if not exists
	techLimit := entity.TechLimit{
		Key: key,
	}
	err := t.gormDb.
		Where("key = ?", key).
		FirstOrCreate(&techLimit).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	// increment usage_count
	err = t.gormDb.
		Model(&entity.TechLimit{}).
		Where("key = ?", key).
		UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).
		UpdateColumn("updated_at", gorm.Expr("current_timestamp")).
		Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: techLimit}
}
