package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type LinkedinConnectionRequestWriteRepository interface {
}

type linkedinConnectionRequestWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewLinkedinConnectionRequestWriteRepository(driver *neo4j.DriverWithContext, database string) LinkedinConnectionRequestWriteRepository {
	return &linkedinConnectionRequestWriteRepository{
		driver:   driver,
		database: database,
	}
}
