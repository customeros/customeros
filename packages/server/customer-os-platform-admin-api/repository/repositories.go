package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	postgresAuthRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Drivers       Drivers
	neo4jDatabase string

	PostgresRepositories     *postgresRepository.Repositories
	PostgresAuthRepositories *postgresAuthRepository.Repositories

	TenantRepository       TenantRepository
	OrganizationRepository OrganizationRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB, neo4jDatabase string) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
		neo4jDatabase:            neo4jDatabase,
		PostgresRepositories:     postgresRepository.InitRepositories(gormDb),
		PostgresAuthRepositories: postgresAuthRepository.InitRepositories(gormDb),
	}
	repositories.OrganizationRepository = NewOrganizationRepository(driver)
	repositories.TenantRepository = NewTenantRepository(driver)

	repositories.PostgresRepositories.Migration(gormDb)
	repositories.PostgresAuthRepositories.Migration(gormDb)

	return &repositories
}
