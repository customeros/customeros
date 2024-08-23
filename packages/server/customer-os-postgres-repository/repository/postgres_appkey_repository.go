package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type appKeyRepository struct {
	db *gorm.DB
}

type AppKeyRepository interface {
	FindByKey(ctx context.Context, app string, key string) helper.QueryResult
}

func NewAppKeyRepo(db *gorm.DB) AppKeyRepository {
	return &appKeyRepository{db: db}
}

func (r *appKeyRepository) FindByKey(ctx context.Context, app string, key string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "AppKeyRepo.FindByKey")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("app", app))

	var appKey entity.AppKey

	err := r.db.
		Where("app_id = ?", app).
		Where("key = ?", key).
		Where("active = ?", true).
		First(&appKey).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &appKey}
}
