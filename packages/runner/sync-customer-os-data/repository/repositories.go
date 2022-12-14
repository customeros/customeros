package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"gorm.io/gorm"
)

type Dbs struct {
	ControlDb      *gorm.DB
	Neo4jDriver    *neo4j.Driver
	AirbyteStoreDB *config.AirbyteStoreDB
}

type Repositories struct {
	Dbs                          Dbs
	TenantSyncSettingsRepository TenantSyncSettingsRepository
}

func InitRepos(driver *neo4j.Driver, controlDb *gorm.DB, airbyteStoreDb *config.AirbyteStoreDB) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver:    driver,
			ControlDb:      controlDb,
			AirbyteStoreDB: airbyteStoreDb,
		},
		TenantSyncSettingsRepository: NewTenantSyncSettingsRepository(controlDb),
	}
	return &repositories
}
