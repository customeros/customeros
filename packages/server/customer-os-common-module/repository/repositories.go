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
	UserSettingsRepository repository.UserSettingsRepository
	AppKeyRepository       repository.AppKeyRepository
	AiPromptLogRepository repository.AiPromptLogRepository
	UserRepository         neo4jrepo.UserRepository
	TenantRepository       neo4jrepo.TenantRepository
	CountryRepository      neo4jrepo.CountryRepository
	StateRepository        neo4jrepo.StateRepository
}

func InitRepositories(db *gorm.DB, driver *neo4j.DriverWithContext) *Repositories {
	repositories := &Repositories{
		UserSettingsRepository: repository.NewUserSettingsRepository(db),
		AppKeyRepository:       repository.NewAppKeyRepo(db),
		AiPromptLogRepository: repository.NewAiPromptLogRepository(db),
		UserRepository:         neo4jrepo.NewUserRepository(driver),
		TenantRepository:       neo4jrepo.NewTenantRepository(driver),
		CountryRepository:      neo4jrepo.NewCountryRepository(driver),
		StateRepository:        neo4jrepo.NewStateRepository(driver),
	}

	var err error

	err = db.AutoMigrate(&entity.AppKey{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	err = db.AutoMigrate(&entity.AiPromptLog{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return repositories
}
