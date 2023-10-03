package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	Drivers                  Drivers
	ExternalSystemRepository ExternalSystemRepository
	UserRepository           UserRepository
	TenantRepository         TenantRepository
	EmailRepository          EmailRepository
	PhoneNumberRepository    PhoneNumberRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

func InitRepos(driver *neo4j.DriverWithContext) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
	}
	repositories.ExternalSystemRepository = NewExternalSystemRepository(driver)
	repositories.UserRepository = NewUserRepository(driver)
	repositories.TenantRepository = NewTenantRepository(driver)
	repositories.EmailRepository = NewEmailRepository(driver)
	repositories.PhoneNumberRepository = NewPhoneNumberRepository(driver)
	return &repositories
}
