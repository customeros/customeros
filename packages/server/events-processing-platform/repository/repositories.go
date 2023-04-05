package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

type Repositories struct {
	Drivers               Drivers
	PhoneNumberRepository PhoneNumberRepository
	ContactRepository     ContactRepository
}

func InitRepos(driver *neo4j.DriverWithContext) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
		PhoneNumberRepository: NewPhoneNumberRepository(driver),
		ContactRepository:     NewContactRepository(driver),
	}
	return &repositories
}
