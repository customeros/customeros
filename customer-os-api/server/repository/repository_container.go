package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type RepositoryContainer struct {
	Drivers                         Drivers
	EntityDefinitionRepository      EntityDefinitionRepository
	FieldSetDefinitionRepository    FieldSetDefinitionRepository
	CustomFieldDefinitionRepository CustomFieldDefinitionRepository
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
	container.FieldSetDefinitionRepository = NewFieldSetDefinitionRepository(driver, &container)
	container.CustomFieldDefinitionRepository = NewCustomFieldDefinitionRepository(driver, &container)
	return &container
}
