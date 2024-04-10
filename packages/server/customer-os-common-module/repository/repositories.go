package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j"
)

type Repositories struct {
	UserRepository   neo4jrepo.UserRepository
	TenantRepository neo4jrepo.TenantRepository
	StateRepository  neo4jrepo.StateRepository
}

func InitRepositories(driver *neo4j.DriverWithContext) *Repositories {
	repositories := &Repositories{
		UserRepository:   neo4jrepo.NewUserRepository(driver),
		TenantRepository: neo4jrepo.NewTenantRepository(driver),
		StateRepository:  neo4jrepo.NewStateRepository(driver),
	}

	return repositories
}
