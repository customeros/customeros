package test

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

type TestDatabase struct {
	Driver *neo4j.DriverWithContext

	Neo4jContainer     testcontainers.Container
	Neo4jRepositories  *repository.Repositories
	PostgresRepository *gorm.DB
}
