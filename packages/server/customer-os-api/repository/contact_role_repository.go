package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ContactRoleRepository interface {
}

type contactRoleRepository struct {
	driver *neo4j.Driver
}

func NewContactRoleRepository(driver *neo4j.Driver) ContactRoleRepository {
	return &contactRoleRepository{
		driver: driver,
	}
}
