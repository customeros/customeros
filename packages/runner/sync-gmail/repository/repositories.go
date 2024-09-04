package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Neo4jDriver *neo4j.DriverWithContext

	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jRepository.Repositories

	//pg repositories
	RawEmailRepository         RawEmailRepository
	RawCalendarEventRepository RawCalendarEventRepository

	//neo4j repositories
	TenantRepository           TenantRepository
	UserRepository             UserRepository
	EmailRepository            EmailRepository
	InteractionEventRepository InteractionEventRepository
	OrganizationRepository     OrganizationRepository
	ActionRepository           ActionRepository
	ActionPointRepository      ActionPointRepository
	DomainRepository           DomainRepository
	MeetingRepository          MeetingRepository
}

func InitRepos(cfg *config.Config, driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Neo4jDriver: driver,

		PostgresRepositories: postgresRepository.InitRepositories(gormDb),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, cfg.Neo4jDb.Database),

		RawEmailRepository:         NewRawEmailRepository(gormDb),
		RawCalendarEventRepository: NewRawCalendarEventRepository(gormDb),

		TenantRepository:           NewTenantRepository(driver),
		UserRepository:             NewUserRepository(driver),
		EmailRepository:            NewEmailRepository(driver),
		InteractionEventRepository: NewInteractionEventRepository(driver),
		OrganizationRepository:     NewOrganizationRepository(driver),
		ActionRepository:           NewActionRepository(driver),
		ActionPointRepository:      NewActionPointRepository(driver),
		DomainRepository:           NewDomainRepository(driver),
		MeetingRepository:          NewMeetingRepository(driver),
	}

	return &repositories
}
