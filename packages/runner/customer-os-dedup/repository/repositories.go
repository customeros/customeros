package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"gorm.io/gorm"
)

type Dbs struct {
	Neo4jDriver *neo4j.DriverWithContext
}

type Repositories struct {
	Dbs Dbs

	CommonRepositories *commonRepository.Repositories

	TenantRepository       TenantRepository
	OrganizationRepository OrganizationRepository
}

func InitRepositories(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver: driver,
		},
		CommonRepositories:     commonRepository.InitRepositories(gormDb, driver),
		TenantRepository:       NewTenantRepository(driver),
		OrganizationRepository: NewOrganizationRepository(driver),
	}
	return &repositories
}
