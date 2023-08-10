package repository

import (
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres"
	"gorm.io/gorm"
)

type PostgresRepositories struct {
	TenantSettingsRepository TenantSettingsRepository
	OAuthTokenRepository     repository.OAuthTokenRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositories {
	p := &PostgresRepositories{
		TenantSettingsRepository: NewTenantSettingsRepository(db),
		OAuthTokenRepository:     repository.NewOAuthTokenRepository(db),
	}

	return p
}
