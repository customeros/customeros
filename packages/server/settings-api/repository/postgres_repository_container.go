package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"gorm.io/gorm"
)

type PostgresRepositories struct {
	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jrepository.Repositories

	TenantSettingsRepository TenantSettingsRepository
}

func InitRepositories(cfg *config.Config, db *gorm.DB, driver *neo4j.DriverWithContext) *PostgresRepositories {
	p := &PostgresRepositories{
		PostgresRepositories:     postgresRepository.InitRepositories(db),
		Neo4jRepositories:        neo4jrepository.InitNeo4jRepositories(driver, cfg.Neo4j.Database),
		TenantSettingsRepository: NewTenantSettingsRepository(db),
	}

	err := db.AutoMigrate(entity.TenantSettings{})
	if err != nil {
		panic(err)
	}

	return p
}
