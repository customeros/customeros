package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Neo4jDriver *neo4j.DriverWithContext

	PostgresRepositories *postgresRepository.Repositories

	//pg repositories
	RawEmailRepository         RawEmailRepository
	RawCalendarEventRepository RawCalendarEventRepository

	//neo4j repositories
	TenantRepository           TenantRepository
	UserRepository             UserRepository
	EmailRepository            EmailRepository
	ExternalSystemRepository   ExternalSystemRepository
	InteractionEventRepository InteractionEventRepository
	OrganizationRepository     OrganizationRepository
	AnalysisRepository         AnalysisRepository
	ActionRepository           ActionRepository
	ActionPointRepository      ActionPointRepository
	DomainRepository           DomainRepository
	MeetingRepository          MeetingRepository
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Neo4jDriver: driver,

		PostgresRepositories: postgresRepository.InitRepositories(gormDb),

		RawEmailRepository:         NewRawEmailRepository(gormDb),
		RawCalendarEventRepository: NewRawCalendarEventRepository(gormDb),

		TenantRepository:           NewTenantRepository(driver),
		UserRepository:             NewUserRepository(driver),
		EmailRepository:            NewEmailRepository(driver),
		ExternalSystemRepository:   NewExternalSystemRepository(driver),
		InteractionEventRepository: NewInteractionEventRepository(driver),
		OrganizationRepository:     NewOrganizationRepository(driver),
		AnalysisRepository:         NewAnalysisRepository(driver),
		ActionRepository:           NewActionRepository(driver),
		ActionPointRepository:      NewActionPointRepository(driver),
		DomainRepository:           NewDomainRepository(driver),
		MeetingRepository:          NewMeetingRepository(driver),
	}

	return &repositories
}
