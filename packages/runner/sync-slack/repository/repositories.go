package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository/neo4j"
	postgresrepo "github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/repository/postgres"
	commonrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"gorm.io/gorm"
)

type Dbs struct {
	Neo4jDriver *neo4j.DriverWithContext
	GormDb      *gorm.DB
}

type Repositories struct {
	Dbs Dbs

	CommonRepositories          *commonrepo.Repositories
	TenantSettingsRepository    postgresrepo.TenantSettingsRepository
	SlackSyncSettingsRepository postgresrepo.SlackSyncSettingsRepository
	SlackSyncRunRepository      postgresrepo.SlackSyncRunRepository

	TenantRepository       neo4jrepo.TenantRepository
	OrganizationRepository neo4jrepo.OrganizationRepository
}

func InitRepositories(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver: driver,
			GormDb:      gormDb,
		},
		CommonRepositories:          commonrepo.InitRepositories(gormDb, driver),
		TenantSettingsRepository:    postgresrepo.NewTenantSettingsRepository(gormDb),
		SlackSyncSettingsRepository: postgresrepo.NewSlackSyncSettingsRepository(gormDb),
		SlackSyncRunRepository:      postgresrepo.NewSlackSyncRunRepository(gormDb),

		TenantRepository:       neo4jrepo.NewTenantRepository(driver),
		OrganizationRepository: neo4jrepo.NewOrganizationRepository(driver),
	}
	return &repositories
}
