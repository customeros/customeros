package repository

import (
	"github.com.openline-ai.customer-os-analytics-api/repository/entity"
	"gorm.io/gorm"
)

type AppInfoRepo struct {
	db *gorm.DB
}

// FIXME alexb replace with authentication
var tenant = "openline"

type AppInfoRepository interface {
	FindAll() RepositoryResult
	FindOneById(id string) RepositoryResult
}

func NewAppInfoRepo(db *gorm.DB) *AppInfoRepo {
	return &AppInfoRepo{db: db}
}

func (r *AppInfoRepo) FindAll() RepositoryResult {
	var applications entity.ApplicationEntities

	err := r.db.Where(&entity.ApplicationEntity{Tenant: tenant}).Find(&applications).Error

	if err != nil {
		return RepositoryResult{Error: err}
	}

	return RepositoryResult{Result: &applications}
}

func (r *AppInfoRepo) FindOneById(id string) RepositoryResult {
	var application entity.ApplicationEntity

	err := r.db.Where(&entity.ApplicationEntity{ID: id, Tenant: tenant}).Take(&application).Error

	if err != nil {
		return RepositoryResult{Error: err}
	}

	return RepositoryResult{Result: &application}
}
