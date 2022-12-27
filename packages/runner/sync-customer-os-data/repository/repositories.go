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
	SyncRunRepository            SyncRunRepository
	ContactRepository            ContactRepository
	ExternalSystemRepository     ExternalSystemRepository
	CompanyRepository            CompanyRepository
	RoleRepository               RoleRepository
	UserRepository               UserRepository
	NoteRepository               NoteRepository
}

func InitRepos(driver *neo4j.Driver, controlDb *gorm.DB, airbyteStoreDb *config.AirbyteStoreDB) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver:    driver,
			ControlDb:      controlDb,
			AirbyteStoreDB: airbyteStoreDb,
		},
		TenantSyncSettingsRepository: NewTenantSyncSettingsRepository(controlDb),
		SyncRunRepository:            NewSyncRunRepository(controlDb),
		ContactRepository:            NewContactRepository(driver),
		ExternalSystemRepository:     NewExternalSystemRepository(driver),
		CompanyRepository:            NewCompanyRepository(driver),
		RoleRepository:               NewRoleRepository(driver),
		UserRepository:               NewUserRepository(driver),
		NoteRepository:               NewNoteRepository(driver),
	}
	return &repositories
}
