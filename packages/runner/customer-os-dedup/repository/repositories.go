package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Dbs struct {
	Neo4jDriver *neo4j.DriverWithContext
}

type Repositories struct {
	Dbs Dbs

	PostgresRepositories *postgresRepository.Repositories

	TenantRepository       TenantRepository
	OrganizationRepository OrganizationRepository
}

func InitRepositories(driver *neo4j.DriverWithContext, gormDb *gorm.DB, neo4jDatabase string) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver: driver,
		},
		PostgresRepositories:   postgresRepository.InitRepositories(gormDb),
		TenantRepository:       NewTenantRepository(driver, neo4jDatabase),
		OrganizationRepository: NewOrganizationRepository(driver, neo4jDatabase),
	}
	return &repositories
}
