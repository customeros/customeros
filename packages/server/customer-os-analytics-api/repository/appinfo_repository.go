package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/helper"
	"gorm.io/gorm"
)

type AppInfoRepo struct {
	db *gorm.DB
}

type AppInfoRepository interface {
	FindAll(ctx context.Context) helper.QueryResult
	FindOneById(ctx context.Context, id string) helper.QueryResult
}

func NewAppInfoRepo(db *gorm.DB) *AppInfoRepo {
	return &AppInfoRepo{db: db}
}

func (r *AppInfoRepo) FindAll(ctx context.Context) helper.QueryResult {
	var applications entity.ApplicationEntities

	err := r.db.Where(&entity.ApplicationEntity{Tenant: common.GetContext(ctx).Tenant}).Find(&applications).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &applications}
}

func (r *AppInfoRepo) FindOneById(ctx context.Context, id string) helper.QueryResult {
	var application entity.ApplicationEntity

	err := r.db.Where(&entity.ApplicationEntity{ID: id, Tenant: common.GetContext(ctx).Tenant}).Take(&application).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &application}
}
