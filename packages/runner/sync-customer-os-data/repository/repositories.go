package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"gorm.io/gorm"
)

type Dbs struct {
	GormDB         *gorm.DB
	Neo4jDriver    *neo4j.DriverWithContext
	RawDataStoreDB *config.RawDataStoreDB
}

type Repositories struct {
	Dbs Dbs

	PostgresRepositories *postgresRepository.Repositories

	TenantSyncSettingsRepository TenantSyncSettingsRepository
	TenantSettingsRepository     TenantSettingsRepository
	SyncRunRepository            SyncRunRepository

	ContactRepository          ContactRepository
	EmailRepository            EmailRepository
	PhoneNumberRepository      PhoneNumberRepository
	LocationRepository         LocationRepository
	ExternalSystemRepository   ExternalSystemRepository
	OrganizationRepository     OrganizationRepository
	UserRepository             UserRepository
	InteractionEventRepository InteractionEventRepository
	MeetingRepository          MeetingRepository
}

func InitRepos(driver *neo4j.DriverWithContext, gormDB *gorm.DB, airbyteStoreDb *config.RawDataStoreDB, log logger.Logger) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver:    driver,
			GormDB:         gormDB,
			RawDataStoreDB: airbyteStoreDb,
		},
		PostgresRepositories:         postgresRepository.InitRepositories(gormDB),
		TenantSyncSettingsRepository: NewTenantSyncSettingsRepository(gormDB),
		TenantSettingsRepository:     NewTenantSettingsRepository(gormDB),
		SyncRunRepository:            NewSyncRunRepository(gormDB),
		ContactRepository:            NewContactRepository(driver),
		EmailRepository:              NewEmailRepository(driver),
		PhoneNumberRepository:        NewPhoneNumberRepository(driver),
		LocationRepository:           NewLocationRepository(driver),
		ExternalSystemRepository:     NewExternalSystemRepository(driver),
		OrganizationRepository:       NewOrganizationRepository(driver, log),
		UserRepository:               NewUserRepository(driver),
		InteractionEventRepository:   NewInteractionEventRepository(driver),
		MeetingRepository:            NewMeetingRepository(driver),
	}
	return &repositories
}
