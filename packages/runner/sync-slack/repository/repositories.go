package repository

import (
	neo4jrepo "github.com/customeros/customeros/packages/runner/sync-slack/repository/neo4j"
	postgresrepo "github.com/customeros/customeros/packages/runner/sync-slack/repository/postgres"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	PostgresRepositories *postgresRepository.Repositories

	SlackSyncSettingsRepository postgresrepo.SlackSyncSettingsRepository
	SlackSyncRunRepository      postgresrepo.SlackSyncRunRepository

	TenantRepository       neo4jrepo.TenantRepository
	OrganizationRepository neo4jrepo.OrganizationRepository
}

func InitRepositories(driver *neo4j.DriverWithContext, postgresDB *config.PostgresDB) *Repositories {
	repositories := Repositories{
		PostgresRepositories: postgresRepository.InitRepositories(postgresDB),

		SlackSyncSettingsRepository: postgresrepo.NewSlackSyncSettingsRepository(postgresDB.GormDB),
		SlackSyncRunRepository:      postgresrepo.NewSlackSyncRunRepository(postgresDB.GormDB),

		TenantRepository:       neo4jrepo.NewTenantRepository(driver),
		OrganizationRepository: neo4jrepo.NewOrganizationRepository(driver),
	}
	return &repositories
}
