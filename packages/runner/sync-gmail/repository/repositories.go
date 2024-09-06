package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Neo4jDriver    *neo4j.DriverWithContext
	PostgresDriver *gorm.DB

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

func InitRepos(cfg *config.Config, driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Neo4jDriver:    driver,
		PostgresDriver: gormDb,

		PostgresRepositories: postgresRepository.InitRepositories(gormDb),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, cfg.Neo4jDb.Database),

		RawEmailRepository:         NewRawEmailRepository(gormDb),
		RawCalendarEventRepository: NewRawCalendarEventRepository(gormDb),

		EmailRepository:            NewEmailRepository(driver),
		InteractionEventRepository: NewInteractionEventRepository(driver),
		OrganizationRepository:     NewOrganizationRepository(driver),
		ActionRepository:           NewActionRepository(driver),
		DomainRepository:           NewDomainRepository(driver),
		MeetingRepository:          NewMeetingRepository(driver),
	}

	return &repositories
}
