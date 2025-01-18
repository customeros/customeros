package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type Repositories struct {
	PostgresRepositories *postgresRepository.Repositories

	TenantRepository       TenantRepository
	OrganizationRepository OrganizationRepository
}

func InitRepositories(driver *neo4j.DriverWithContext, neo4jDatabase string, postgresDB *commonConfig.PostgresDB) *Repositories {
	repositories := Repositories{
		PostgresRepositories:   postgresRepository.InitRepositories(postgresDB),
		TenantRepository:       NewTenantRepository(driver, neo4jDatabase),
		OrganizationRepository: NewOrganizationRepository(driver, neo4jDatabase),
	}
	return &repositories
}
