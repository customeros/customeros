package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	cmn_repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	repository "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"gorm.io/gorm"
)

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

type Repositories struct {
	Drivers Drivers

	Neo4jRepositories       *neo4jrepository.Repositories
	CommonRepositories      *cmn_repository.Repositories
	CustomerOsIdsRepository repository.CustomerOsIdsRepository

	EmailRepository            EmailRepository
	ExternalSystemRepository   ExternalSystemRepository
	InteractionEventRepository InteractionEventRepository
	LocationRepository         LocationRepository
	OpportunityRepository      OpportunityRepository
	OrganizationRepository     OrganizationRepository
}

func InitRepos(driver *neo4j.DriverWithContext, neo4jDatabase string, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
		Neo4jRepositories:          neo4jrepository.InitNeo4jRepositories(driver, neo4jDatabase),
		CommonRepositories:         cmn_repository.InitRepositories(gormDb, driver),
		CustomerOsIdsRepository:    repository.NewCustomerOsIdsRepository(gormDb),
		EmailRepository:            NewEmailRepository(driver),
		OrganizationRepository:     NewOrganizationRepository(driver, neo4jDatabase),
		LocationRepository:         NewLocationRepository(driver),
		InteractionEventRepository: NewInteractionEventRepository(driver, neo4jDatabase),
		ExternalSystemRepository:   NewExternalSystemRepository(driver),
		OpportunityRepository:      NewOpportunityRepository(driver, neo4jDatabase),
	}

	err := gormDb.AutoMigrate(&entity.CustomerOsIds{})
	if err != nil {
		panic(err)
	}

	return &repositories
}
