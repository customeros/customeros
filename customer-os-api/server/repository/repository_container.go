package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type RepositoryContainer struct {
	Drivers                    Drivers
	EntityDefinitionRepository EntityDefinitionRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.Driver
}

func InitRepos(driver *neo4j.Driver) *RepositoryContainer {
	container := RepositoryContainer{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
	}
	container.EntityDefinitionRepository = NewEntityDefinitionRepository(driver, &container)
	return &container
}
