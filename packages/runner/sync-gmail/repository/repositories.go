package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/runner/sync-gmail/config"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	neo4jRepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type Repositories struct {
	Neo4jDriver *neo4j.DriverWithContext

	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jRepository.Repositories

	//pg repositories
	RawEmailRepository         RawEmailRepository
	RawCalendarEventRepository RawCalendarEventRepository

	//neo4j repositories
	EmailRepository            EmailRepository
	InteractionEventRepository InteractionEventRepository
	OrganizationRepository     OrganizationRepository
	ActionRepository           ActionRepository
	DomainRepository           DomainRepository
	MeetingRepository          MeetingRepository
}

func InitRepos(cfg *config.Config, driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB) *Repositories {
	repositories := Repositories{
		Neo4jDriver: driver,

		PostgresRepositories: postgresRepository.InitRepositories(postgresDB),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, cfg.Neo4jDb.Database),

		RawEmailRepository:         NewRawEmailRepository(postgresDB.AsyncGormDB),
		RawCalendarEventRepository: NewRawCalendarEventRepository(postgresDB.AsyncGormDB),

		EmailRepository:            NewEmailRepository(driver),
		InteractionEventRepository: NewInteractionEventRepository(driver),
		OrganizationRepository:     NewOrganizationRepository(driver),
		ActionRepository:           NewActionRepository(driver),
		DomainRepository:           NewDomainRepository(driver),
		MeetingRepository:          NewMeetingRepository(driver),
	}

	return &repositories
}
