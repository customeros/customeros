package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	repository "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"gorm.io/gorm"
)

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
	GormDb      *gorm.DB
}

type Repositories struct {
	Drivers Drivers

	Neo4jRepositories    *neo4jrepository.Repositories
	PostgresRepositories *postgresRepository.Repositories
	InvoiceRepository    repository.InvoiceRepository
}

func InitRepos(driver *neo4j.DriverWithContext, neo4jDatabase string, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
			GormDb:      gormDb,
		},
		Neo4jRepositories:    neo4jrepository.InitNeo4jRepositories(driver, neo4jDatabase),
		PostgresRepositories: postgresRepository.InitRepositories(gormDb),
		InvoiceRepository:    repository.NewInvoiceRepository(gormDb),
	}

	return &repositories
}

func Migration(db *gorm.DB) {

	err := db.AutoMigrate(&entity.InvoiceNumberEntity{})
	if err != nil {
		panic(err)
	}
}
