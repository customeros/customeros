package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"gorm.io/gorm"
)

type DbDrivers struct {
	Neo4jDriver *neo4j.Driver
	GormDb      *gorm.DB
}

type Repositories struct {
	Drivers                     DbDrivers
	ContactActionItemRepository ContactActionItemRepository
	TrackedVisitorRepository    TrackedVisitorRepository
}

func InitRepos(driver *neo4j.Driver, gormDb *gorm.DB) *Repositories {
	container := Repositories{
		Drivers: DbDrivers{
			Neo4jDriver: driver,
			GormDb:      gormDb,
		},
	}
	container.ContactActionItemRepository = NewContactActionItemRepository(driver, &container)
	container.TrackedVisitorRepository = NewTrackedVisitorRepository(gormDb)
	return &container
}
