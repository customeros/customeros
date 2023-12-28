package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repositories struct {
	ContractReadRepository ContractReadRepository
	UserReadRepository     UserReadRepository
	UserWriteRepository    UserWriteRepository
}

func InitNeo4jRepositories(driver *neo4j.DriverWithContext, neo4jDatabase string) *Repositories {
	repositories := Repositories{
		ContractReadRepository: NewContractReadRepository(driver, neo4jDatabase),
		UserReadRepository:     NewUserReadRepository(driver, neo4jDatabase),
		UserWriteRepository:    NewUserWriteRepository(driver, neo4jDatabase),
	}
	return &repositories
}
