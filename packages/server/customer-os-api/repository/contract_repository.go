package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ContractRepository interface {
}

type contractRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContractRepository(driver *neo4j.DriverWithContext) ContractRepository {
	return &contractRepository{
		driver: driver,
	}
}
