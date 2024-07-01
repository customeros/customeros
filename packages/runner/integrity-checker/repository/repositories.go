package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/config"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
)

type Dbs struct {
	Neo4jDriver *neo4j.DriverWithContext
}

type Repositories struct {
	Dbs               Dbs
	Neo4jRepositories *neo4jRepository.Repositories
}

func InitRepositories(cfg *config.Config, driver *neo4j.DriverWithContext) *Repositories {
	repositories := Repositories{
		Dbs: Dbs{
			Neo4jDriver: driver,
		},
		Neo4jRepositories: neo4jRepository.InitNeo4jRepositories(driver, cfg.Neo4j.Database),
	}
	return &repositories
}
