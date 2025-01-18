package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type Repositories struct {
	Drivers              Drivers
	Neo4jRepositories    *neo4jrepository.Repositories
	PostgresRepositories *postgresRepository.Repositories

	ActionItemRepository           ActionItemRepository
	CalendarRepository             CalendarRepository
	ContactRepository              ContactRepository
	CustomFieldRepository          CustomFieldRepository
	CustomFieldsTemplateRepository CustomFieldTemplateRepository
	DashboardRepository            DashboardRepository
	DashboardV2Repository          DashboardV2Repository
	EmailRepository                EmailRepository
	ExternalSystemRepository       ExternalSystemRepository
	IssueRepository                IssueRepository
	LocationRepository             LocationRepository
	MeetingRepository              MeetingRepository
	NoteRepository                 NoteRepository
	OrganizationRepository         OrganizationRepository
	SearchRepository               SearchRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

func InitRepos(driver *neo4j.DriverWithContext, database string, postgresDB *commonConfig.PostgresDB) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
	}
	repositories.Neo4jRepositories = neo4jrepository.InitNeo4jRepositories(driver, database)
	repositories.PostgresRepositories = postgresRepository.InitRepositories(postgresDB)

	repositories.ActionItemRepository = NewActionItemRepository(driver)
	repositories.CalendarRepository = NewCalendarRepository(driver)
	repositories.ContactRepository = NewContactRepository(driver, database)
	repositories.CustomFieldRepository = NewCustomFieldRepository(driver, database)
	repositories.CustomFieldsTemplateRepository = NewCustomFieldTemplateRepository(driver, database)
	repositories.DashboardRepository = NewDashboardRepository(driver)
	repositories.DashboardV2Repository = NewDashboardV2Repository(driver)
	repositories.EmailRepository = NewEmailRepository(driver, database)
	repositories.ExternalSystemRepository = NewExternalSystemRepository(driver)
	repositories.IssueRepository = NewIssueRepository(driver, database)
	repositories.LocationRepository = NewLocationRepository(driver)
	repositories.MeetingRepository = NewMeetingRepository(driver)
	repositories.NoteRepository = NewNoteRepository(driver)
	repositories.OrganizationRepository = NewOrganizationRepository(driver, database)
	repositories.SearchRepository = NewSearchRepository(driver)
	return &repositories
}
