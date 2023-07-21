package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

type Repositories struct {
	Drivers                Drivers
	ContactRepository      ContactRepository
	OrganizationRepository OrganizationRepository
	PhoneNumberRepository  PhoneNumberRepository
	EmailRepository        EmailRepository
	UserRepository         UserRepository
	LocationRepository     LocationRepository
	CountryRepository      CountryRepository
	JobRoleRepository      JobRoleRepository
	SocialRepository       SocialRepository
}

func InitRepos(driver *neo4j.DriverWithContext) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
		PhoneNumberRepository:  NewPhoneNumberRepository(driver),
		EmailRepository:        NewEmailRepository(driver),
		ContactRepository:      NewContactRepository(driver),
		OrganizationRepository: NewOrganizationRepository(driver),
		UserRepository:         NewUserRepository(driver),
		LocationRepository:     NewLocationRepository(driver),
		CountryRepository:      NewCountryRepository(driver),
		JobRoleRepository:      NewJobRoleRepository(driver),
		SocialRepository:       NewSocialRepository(driver),
	}
	return &repositories
}
