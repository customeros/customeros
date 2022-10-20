package repository

import (
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
	"github.com.openline-ai.customer-os-analytics-api/repository/helper"
	"gorm.io/gorm"
)

type AppInfoRepo struct {
	db *gorm.DB
}

// FIXME alexb replace with authentication
var tenant = "openline"

type AppInfoRepository interface {
	FindAll() helper.QueryResult
	FindOneById(id string) helper.QueryResult
}

func NewAppInfoRepo(db *gorm.DB) *AppInfoRepo {
	return &AppInfoRepo{db: db}
}

func (r *AppInfoRepo) FindAll() helper.QueryResult {
	var applications entity.ApplicationEntities

	err := r.db.Where(&entity.ApplicationEntity{Tenant: tenant}).Find(&applications).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &applications}
}

func (r *AppInfoRepo) FindOneById(id string) helper.QueryResult {
	var application entity.ApplicationEntity

	err := r.db.Where(&entity.ApplicationEntity{ID: id, Tenant: tenant}).Take(&application).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &application}
}
