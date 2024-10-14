package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TagRepository interface {
}

type tagRepository struct {
	driver *neo4j.DriverWithContext
}

func NewTagRepository(driver *neo4j.DriverWithContext) TagRepository {
	return &tagRepository{
		driver: driver,
	}
}
