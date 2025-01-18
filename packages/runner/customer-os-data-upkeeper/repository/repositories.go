package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	neo4jRepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

// deprecated
type Repositories struct {
	PostgresRepositories *postgresRepository.Repositories
	Neo4jRepositories    *neo4jRepository.Repositories
}

// deprecated
func InitRepositories(cfg *config.Config, driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB) *Repositories {
	repositories := Repositories{
		PostgresRepositories: postgresRepository.InitRepositories(postgresDB),
		Neo4jRepositories:    neo4jRepository.InitNeo4jRepositories(driver, cfg.Common.Infrastructure.Neo4jConfig.Database),
	}
	return &repositories
}
