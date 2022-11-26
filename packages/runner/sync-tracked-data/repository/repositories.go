package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/gen"
)

type DbDrivers struct {
	Neo4jDriver *neo4j.Driver
	EntClient   *gen.Client
}

type Repositories struct {
	Drivers                     DbDrivers
	ContactActionItemRepository ContactActionItemRepository
	TrackedVisitorRepository    TrackedVisitorRepository
}

func InitRepos(driver *neo4j.Driver, client *gen.Client) *Repositories {
	container := Repositories{
		Drivers: DbDrivers{
			Neo4jDriver: driver,
			EntClient:   client,
		},
	}
	container.ContactActionItemRepository = NewContactActionItemRepository(driver, &container)
	container.TrackedVisitorRepository = NewTrackedVisitorRepository(client)
	return &container
}
