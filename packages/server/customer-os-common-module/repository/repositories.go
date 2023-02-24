package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
	"log"
)

type Repositories struct {
	AppKeyRepository  repository.AppKeyRepository
	UserRepository    neo4jrepo.UserRepository
	CountryRepository neo4jrepo.CountryRepository
}

func InitRepositories(db *gorm.DB, driver *neo4j.DriverWithContext) *Repositories {
	repositories := &Repositories{
		AppKeyRepository:  repository.NewAppKeyRepo(db),
		UserRepository:    neo4jrepo.NewUserRepository(driver),
		CountryRepository: neo4jrepo.NewCountryRepository(driver),
	}

	var err error

	err = db.AutoMigrate(&entity.AppKey{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return repositories
}
