package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository/neo4j"
	commrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"gorm.io/gorm"
)

type Dbs struct {
	Neo4jDriver *neo4j.DriverWithContext
	GormDb      *gorm.DB
}

type Repositories struct {
	Dbs Dbs

	CommonRepositories *commrepo.Repositories

	Neo4jRepositories *neo4jRepository.Repositories

	OrganizationRepository neo4jrepo.OrganizationRepository
}

func InitRepositories(cfg *config.Config, driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver: driver,
			GormDb:      gormDb,
		},
		CommonRepositories: commrepo.InitRepositories(gormDb, driver),
		Neo4jRepositories:  neo4jRepository.InitNeo4jRepositories(driver, cfg.Neo4j.Database),

		OrganizationRepository: neo4jrepo.NewOrganizationRepository(driver),
	}
	return &repositories
}
