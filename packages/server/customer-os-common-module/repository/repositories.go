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
	AppKeyRepository                    repository.AppKeyRepository
	PersonalIntegrationRepository       repository.PersonalIntegrationRepository
	AiPromptLogRepository               repository.AiPromptLogRepository
	ImportAllowedOrganizationRepository repository.ImportAllowedOrganizationRepository
	UserRepository                      neo4jrepo.UserRepository
	TenantRepository                    neo4jrepo.TenantRepository
	CountryRepository                   neo4jrepo.CountryRepository
	StateRepository                     neo4jrepo.StateRepository
}

func InitRepositories(db *gorm.DB, driver *neo4j.DriverWithContext) *Repositories {
	repositories := &Repositories{
		AppKeyRepository:                    repository.NewAppKeyRepo(db),
		PersonalIntegrationRepository:       repository.NewPersonalIntegrationsRepo(db),
		AiPromptLogRepository:               repository.NewAiPromptLogRepository(db),
		ImportAllowedOrganizationRepository: repository.NewImportAllowedOrganizationRepository(db),
		UserRepository:                      neo4jrepo.NewUserRepository(driver),
		TenantRepository:                    neo4jrepo.NewTenantRepository(driver),
		CountryRepository:                   neo4jrepo.NewCountryRepository(driver),
		StateRepository:                     neo4jrepo.NewStateRepository(driver),
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

	err = db.AutoMigrate(&entity.ImportAllowedOrganization{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	err = db.AutoMigrate(&entity.PersonalIntegration{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return repositories
}
