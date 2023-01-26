package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type UserToTenantRepo struct {
	db *gorm.DB
}

type UserToTenantRepository interface {
	FindTenantByUsername(username string) helper.QueryResult
}

func NewUserToTenantRepo(db *gorm.DB) *UserToTenantRepo {
	return &UserToTenantRepo{db: db}
}

func (r *UserToTenantRepo) FindTenantByUsername(username string) helper.QueryResult {
	var e entity.UserToTenant

	err := r.db.
		Where("username = ?", username).
		First(&e).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: e.Tenant}
}
