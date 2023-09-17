package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"gorm.io/gorm"
)

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

type Repositories struct {
	Drivers Drivers

	CommonRepositories *commonRepository.Repositories

	ContactRepository          ContactRepository
	OrganizationRepository     OrganizationRepository
	PhoneNumberRepository      PhoneNumberRepository
	EmailRepository            EmailRepository
	UserRepository             UserRepository
	LocationRepository         LocationRepository
	CountryRepository          CountryRepository
	JobRoleRepository          JobRoleRepository
	SocialRepository           SocialRepository
	InteractionEventRepository InteractionEventRepository
	ActionRepository           ActionRepository
	LogEntryRepository         LogEntryRepository
	TagRepository              TagRepository
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
		CommonRepositories:         commonRepository.InitRepositories(gormDb, driver),
		PhoneNumberRepository:      NewPhoneNumberRepository(driver),
		EmailRepository:            NewEmailRepository(driver),
		ContactRepository:          NewContactRepository(driver),
		OrganizationRepository:     NewOrganizationRepository(driver),
		UserRepository:             NewUserRepository(driver),
		LocationRepository:         NewLocationRepository(driver),
		CountryRepository:          NewCountryRepository(driver),
		JobRoleRepository:          NewJobRoleRepository(driver),
		SocialRepository:           NewSocialRepository(driver),
		InteractionEventRepository: NewInteractionEventRepository(driver),
		ActionRepository:           NewActionRepository(driver),
		LogEntryRepository:         NewLogEntryRepository(driver),
		TagRepository:              NewTagRepository(driver),
	}
	return &repositories
}
