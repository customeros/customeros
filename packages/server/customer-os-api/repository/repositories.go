package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Drivers              Drivers
	Neo4jRepositories    *neo4jrepository.Repositories
	PostgresRepositories *postgresRepository.Repositories

	//Deprecated
	OrganizationRepository OrganizationRepository
	//Deprecated
	ContactRepository ContactRepository
	//Deprecated
	CustomFieldTemplateRepository CustomFieldTemplateRepository
	//Deprecated
	CustomFieldRepository CustomFieldRepository
	//Deprecated
	EntityTemplateRepository EntityTemplateRepository
	//Deprecated
	UserRepository UserRepository
	//Deprecated
	ExternalSystemRepository ExternalSystemRepository
	//Deprecated
	NoteRepository NoteRepository
	//Deprecated
	CalendarRepository CalendarRepository
	LocationRepository LocationRepository
	//Deprecated
	EmailRepository EmailRepository
	//Deprecated
	PhoneNumberRepository PhoneNumberRepository
	//Deprecated
	TagRepository TagRepository
	//Deprecated
	SearchRepository SearchRepository
	//Deprecated
	DashboardRepository DashboardRepository
	//Deprecated
	IssueRepository IssueRepository
	//Deprecated
	MeetingRepository MeetingRepository
	//Deprecated
	ActionRepository ActionRepository
	//Deprecated
	ActionItemRepository ActionItemRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

func InitRepos(driver *neo4j.DriverWithContext, database string, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
	}
	repositories.Neo4jRepositories = neo4jrepository.InitNeo4jRepositories(driver, database)
	repositories.PostgresRepositories = postgresRepository.InitRepositories(gormDb)

	repositories.OrganizationRepository = NewOrganizationRepository(driver, database)
	repositories.ContactRepository = NewContactRepository(driver, database)
	repositories.CustomFieldTemplateRepository = NewCustomFieldTemplateRepository(driver, database)
	repositories.CustomFieldRepository = NewCustomFieldRepository(driver, database)
	repositories.EntityTemplateRepository = NewEntityTemplateRepository(driver, &repositories)
	repositories.UserRepository = NewUserRepository(driver, database)
	repositories.ExternalSystemRepository = NewExternalSystemRepository(driver)
	repositories.NoteRepository = NewNoteRepository(driver)
	repositories.CalendarRepository = NewCalendarRepository(driver)
	repositories.LocationRepository = NewLocationRepository(driver)
	repositories.EmailRepository = NewEmailRepository(driver, database)
	repositories.PhoneNumberRepository = NewPhoneNumberRepository(driver)
	repositories.TagRepository = NewTagRepository(driver)
	repositories.SearchRepository = NewSearchRepository(driver)
	repositories.DashboardRepository = NewDashboardRepository(driver)
	repositories.IssueRepository = NewIssueRepository(driver, database)
	repositories.MeetingRepository = NewMeetingRepository(driver)
	repositories.ActionRepository = NewActionRepository(driver)
	repositories.ActionItemRepository = NewActionItemRepository(driver)
	return &repositories
}
