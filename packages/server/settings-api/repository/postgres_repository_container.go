package repository

import (
	"gorm.io/gorm"
)

type PostgresRepositories struct {
	TenantSettingsRepository TenantSettingsRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositories {
	p := &PostgresRepositories{
		TenantSettingsRepository: NewTenantSettingsRepository(db),
	}

	return p
}
