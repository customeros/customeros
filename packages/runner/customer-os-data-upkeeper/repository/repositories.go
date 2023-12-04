package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository/neo4j"
	commrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"gorm.io/gorm"
)

type Dbs struct {
	Neo4jDriver *neo4j.DriverWithContext
	GormDb      *gorm.DB
}

type Repositories struct {
	Dbs Dbs

	CommonRepositories *commrepo.Repositories

	ContractRepository neo4jrepo.ContractRepository
}

func InitRepositories(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver: driver,
			GormDb:      gormDb,
		},
		CommonRepositories: commrepo.InitRepositories(gormDb, driver),

		ContractRepository: neo4jrepo.NewContractRepository(driver),
	}
	return &repositories
}
