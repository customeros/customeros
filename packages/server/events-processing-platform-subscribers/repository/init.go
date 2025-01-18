package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	neo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type Repositories struct {
	Neo4jRepositories    *neo.Repositories
	PostgresRepositories *postgres.Repositories
}

func InitRepos(driver *neo4j.DriverWithContext, database string, postgresDB *config.PostgresDB) *Repositories {

	repos := Repositories{}

	repos.Neo4jRepositories = neo.InitNeo4jRepositories(driver, database)
	repos.PostgresRepositories = postgres.InitRepositories(postgresDB)

	return &repos
}
