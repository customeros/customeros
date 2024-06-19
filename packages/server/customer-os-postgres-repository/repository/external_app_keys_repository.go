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

type externalAppKeysRepository struct {
	gormDb *gorm.DB
}

type ExternalAppKeysRepository interface {
	IncrementUsageCount(ctx context.Context, id uint64) helper.QueryResult
	GetAppKeys(ctx context.Context, app, group string, usageLimit int) helper.QueryResult
}

func NewExternalAppKeysRepository(gormDb *gorm.DB) ExternalAppKeysRepository {
	return &externalAppKeysRepository{gormDb: gormDb}
}

func (e externalAppKeysRepository) GetAppKeys(ctx context.Context, app, group string, usageLimit int) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "ExternalAppKeysRepository.GetAppKeys")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.LogFields(log.String("app", app), log.String("group1", group))

	var appKeys []entity.ExternalAppKeys
	err := e.gormDb.
		Where("app = ? AND group1 = ? AND usage_count < ?", app, group, usageLimit).
		Find(&appKeys).Limit(10).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: appKeys}
}

func (e externalAppKeysRepository) IncrementUsageCount(ctx context.Context, id uint64) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "ExternalAppKeysRepository.IncrementUsageCount")
	defer span.Finish()
	span.SetTag(tracing.SpanTagComponent, "postgresRepository")
	span.LogFields(log.Uint64("id", id))

	// create entry if not exists
	appKey := entity.ExternalAppKeys{
		ID: id,
	}
	err := e.gormDb.
		Where("id = ?", id).
		First(&appKey).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	// increment usage_count
	err = e.gormDb.
		Model(&entity.ExternalAppKeys{}).
		Where("id = ?", id).
		UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).
		UpdateColumn("updated_at", gorm.Expr("current_timestamp")).
		Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: appKey}
}
