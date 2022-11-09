package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type EntityDefinitionRepository interface {
}

type entityDefinitionRepository struct {
	driver *neo4j.Driver
}

func NewEntityDefinitionRepository(driver *neo4j.Driver) EntityDefinitionRepository {
	return &entityDefinitionRepository{
		driver: driver,
	}
}
