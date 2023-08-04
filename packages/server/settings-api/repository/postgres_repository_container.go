package repository

import (
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"gorm.io/gorm"
)

type PostgresRepositories struct {
	TenantSettingsRepository TenantSettingsRepository
	UserSettingRepository    repository.UserSettingsRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositories {
	p := &PostgresRepositories{
		TenantSettingsRepository: NewTenantSettingsRepository(db),
		UserSettingRepository:    repository.NewUserSettingsRepository(db),
	}

	return p
}
