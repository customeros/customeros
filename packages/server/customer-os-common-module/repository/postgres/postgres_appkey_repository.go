package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type AppKeyRepo struct {
	db *gorm.DB
}

type AppKeyRepository interface {
	FindByKey(ctx context.Context, key string) helper.QueryResult
}

func NewAppKeyRepo(db *gorm.DB) *AppKeyRepo {
	return &AppKeyRepo{db: db}
}

func (r *AppKeyRepo) FindByKey(ctx context.Context, key string) helper.QueryResult {
	var appKey entity.AppKeyEntity

	err := r.db.
		Where("app_id = ?", "customer-os-api").
		Where("key = ?", key).
		Where("active = ?", true).
		First(&appKey).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &appKey}
}
