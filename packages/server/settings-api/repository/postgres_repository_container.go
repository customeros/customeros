package repository

import (
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"gorm.io/gorm"
)

type PostgresRepositories struct {
	PersonalIntegrationRepository commonRepository.PersonalIntegrationRepository
	TenantSettingsRepository      TenantSettingsRepository
	OAuthTokenRepository          repository.OAuthTokenRepository
}

func InitRepositories(db *gorm.DB) *PostgresRepositories {
	p := &PostgresRepositories{
		TenantSettingsRepository:      NewTenantSettingsRepository(db),
		OAuthTokenRepository:          repository.NewOAuthTokenRepository(db),
		PersonalIntegrationRepository: commonRepository.NewPersonalIntegrationsRepo(db),
	}

	return p
}
