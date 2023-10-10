package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	authRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"gorm.io/gorm"
)

type PostgresRepositories struct {
	CommonRepositories *commonRepository.Repositories
	AuthRepositories   *authRepository.Repositories

	TenantSettingsRepository TenantSettingsRepository
}

func InitRepositories(db *gorm.DB, driver *neo4j.DriverWithContext) *PostgresRepositories {
	p := &PostgresRepositories{
		CommonRepositories:       commonRepository.InitRepositories(db, driver),
		AuthRepositories:         authRepository.InitRepositories(db),
		TenantSettingsRepository: NewTenantSettingsRepository(db),
	}

	return p
}
