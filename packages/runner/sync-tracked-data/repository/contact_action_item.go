package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ContactActionItemRepository interface {
}

type userRepository struct {
	driver *neo4j.Driver
	repos  *Repositories
}

func NewContactActionItemRepository(driver *neo4j.Driver, repos *Repositories) ContactActionItemRepository {
	return &userRepository{
		driver: driver,
		repos:  repos,
	}
}
