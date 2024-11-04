package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type LinkedinConnectionRequestReadRepository interface {
}

type linkedinConnectionRequestReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewLinkedinConnectionRequestReadRepository(driver *neo4j.DriverWithContext, database string) LinkedinConnectionRequestReadRepository {
	return &linkedinConnectionRequestReadRepository{
		driver:   driver,
		database: database,
	}
}
