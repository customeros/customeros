package repository

import (
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/core-crm/internal/database"
)

type Repositories struct {
	CommonPostgres *postgres_repository.Repositories
	ICPDescription ICPDescriptionRepository
}

func InitRepos(openlineDB *database.DatabaseConnection) *Repositories {
	pg := &commonConfig.PostgresDB{
		GormDB: openlineDB.WriteDB,
	}

	return &Repositories{
		CommonPostgres: postgres_repository.InitRepositories(pg),
		ICPDescription: NewICPDescriptionRepository(openlineDB),
	}
}
