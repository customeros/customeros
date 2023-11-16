package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Drivers       Drivers
	neo4jDatabase string

	CommonRepositories *commonRepository.Repositories

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
		neo4jDatabase:      neo4jDatabase,
		CommonRepositories: commonRepository.InitRepositories(gormDb, driver),
	}
	repositories.OrganizationRepository = NewOrganizationRepository(driver)
	repositories.TenantRepository = NewTenantRepository(driver)

	return &repositories
}
