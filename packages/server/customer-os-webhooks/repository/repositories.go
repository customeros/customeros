package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	Drivers                  Drivers
	ExternalSystemRepository ExternalSystemRepository
	UserRepository           UserRepository
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
	return &repositories
}
