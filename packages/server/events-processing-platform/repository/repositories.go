package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	repository "github.com/customeros/customeros/packages/server/events-processing-platform/repository/postgres"
	"github.com/customeros/customeros/packages/server/events-processing-platform/repository/postgres/entity"
	"gorm.io/gorm"
)

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
	PostgresDB  *commonConfig.PostgresDB
}

type Repositories struct {
	Drivers Drivers

	Neo4jRepositories    *neo4jrepository.Repositories
	PostgresRepositories *postgresRepository.Repositories
	InvoiceRepository    repository.InvoiceRepository
}

func InitRepos(driver *neo4j.DriverWithContext, neo4jDatabase string, postgresDB *commonConfig.PostgresDB) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
			PostgresDB:  postgresDB,
		},
		Neo4jRepositories:    neo4jrepository.InitNeo4jRepositories(driver, neo4jDatabase),
		PostgresRepositories: postgresRepository.InitRepositories(postgresDB),
		InvoiceRepository:    repository.NewInvoiceRepository(postgresDB.GormDB),
	}

	return &repositories
}

func Migration(db *gorm.DB) {

	err := db.AutoMigrate(&entity.InvoiceNumberEntity{})
	if err != nil {
		panic(err)
	}
}
