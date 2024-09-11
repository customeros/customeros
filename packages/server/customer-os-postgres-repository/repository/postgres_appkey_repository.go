package repository

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type appKeyRepository struct {
	db *gorm.DB
}

type AppKeyRepository interface {
	FindByKey(ctx context.Context, app string, key string) (*entity.AppKey, error)
}

func NewAppKeyRepo(db *gorm.DB) AppKeyRepository {
	return &appKeyRepository{db: db}
}

func (r *appKeyRepository) FindByKey(ctx context.Context, app string, key string) (*entity.AppKey, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AppKeyRepo.FindByKey")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("app", app)

	var appKey entity.AppKey

	err := r.db.
		Where("app_id = ?", app).
		Where("key = ?", key).
		Where("active = ?", true).
		First(&appKey).Error

	// check if the app key is not found
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogFields(log.Bool("result.found", false))
			return nil, nil
		}
		return nil, err
	}

	span.LogFields(log.Bool("result.found", true))
	return &appKey, nil
}
