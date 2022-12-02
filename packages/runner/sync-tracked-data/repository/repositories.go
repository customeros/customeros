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
	Drivers            DbDrivers
	ContactRepository  ContactRepository
	PageViewRepository PageViewRepository
	ActionRepository   ActionRepository
}

func InitRepos(driver *neo4j.Driver, db *gorm.DB) *Repositories {
	repositories := Repositories{
		Drivers: DbDrivers{
			Neo4jDriver: driver,
			GormDb:      db,
		},
	}
	repositories.ContactRepository = NewContactRepository(driver)
	repositories.PageViewRepository = NewPageViewRepository(db)
	repositories.ActionRepository = NewActionRepository(driver)
	return &repositories
}
