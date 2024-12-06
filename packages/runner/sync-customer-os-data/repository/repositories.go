package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
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
	ExternalSystemRepository   ExternalSystemRepository
	OrganizationRepository     OrganizationRepository
	UserRepository             UserRepository
	InteractionEventRepository InteractionEventRepository
	MeetingRepository          MeetingRepository
}

func InitRepos(driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB, airbyteStoreDb *config.RawDataStoreDB, log logger.Logger) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver:    driver,
			GormDB:         postgresDB.GormDB,
			RawDataStoreDB: airbyteStoreDb,
		},
		PostgresRepositories:         postgresRepository.InitRepositories(postgresDB),
		TenantSyncSettingsRepository: NewTenantSyncSettingsRepository(postgresDB.GormDB),
		TenantSettingsRepository:     NewTenantSettingsRepository(postgresDB.GormDB),
		SyncRunRepository:            NewSyncRunRepository(postgresDB.GormDB),
		ContactRepository:            NewContactRepository(driver),
		EmailRepository:              NewEmailRepository(driver),
		ExternalSystemRepository:     NewExternalSystemRepository(driver),
		OrganizationRepository:       NewOrganizationRepository(driver, log),
		UserRepository:               NewUserRepository(driver),
		InteractionEventRepository:   NewInteractionEventRepository(driver),
		MeetingRepository:            NewMeetingRepository(driver),
	}
	return &repositories
}
