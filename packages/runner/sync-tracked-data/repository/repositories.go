package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"gorm.io/gorm"
)

type DbDrivers struct {
	Neo4jDriver    *neo4j.Driver
	GormDb         *gorm.DB
	GormTrackingDb *gorm.DB
}

type Repositories struct {
	Drivers            DbDrivers
	ContactRepository  ContactRepository
	PageViewRepository PageViewRepository
	ActionRepository   ActionRepository
	SyncRunRepository  SyncRunRepository
}

func InitRepos(driver *neo4j.Driver, gormDb *gorm.DB, gormTrackingDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Drivers: DbDrivers{
			Neo4jDriver:    driver,
			GormDb:         gormDb,
			GormTrackingDb: gormTrackingDb,
		},
	}
	repositories.ContactRepository = NewContactRepository(driver)
	repositories.ActionRepository = NewActionRepository(driver)
	repositories.PageViewRepository = NewPageViewRepository(gormTrackingDb)
	repositories.SyncRunRepository = NewSyncRunRepository(gormDb)
	return &repositories
}
