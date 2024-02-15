package repository

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	AppKeyRepository                repository.AppKeyRepository
	PersonalIntegrationRepository   repository.PersonalIntegrationRepository
	AiPromptLogRepository           repository.AiPromptLogRepository
	WhitelistDomainRepository       repository.WhitelistDomainRepository
	PersonalEmailProviderRepository repository.PersonalEmailProviderRepository
	UserRepository                  neo4jrepo.UserRepository
	TenantRepository                neo4jrepo.TenantRepository
	StateRepository                 neo4jrepo.StateRepository
	TenantApiKeyRepository          repository.TenantApiKeyRepository
	TenantWebhookRepository         repository.TenantWebhookRepository
}

func InitRepositories(db *gorm.DB, driver *neo4j.DriverWithContext) *Repositories {
	repositories := &Repositories{
		AppKeyRepository:                repository.NewAppKeyRepo(db),
		PersonalIntegrationRepository:   repository.NewPersonalIntegrationsRepo(db),
		AiPromptLogRepository:           repository.NewAiPromptLogRepository(db),
		WhitelistDomainRepository:       repository.NewWhitelistDomainRepository(db),
		PersonalEmailProviderRepository: repository.NewPersonalEmailProviderRepository(db),
		UserRepository:                  neo4jrepo.NewUserRepository(driver),
		TenantRepository:                neo4jrepo.NewTenantRepository(driver),
		StateRepository:                 neo4jrepo.NewStateRepository(driver),
		TenantApiKeyRepository:          repository.NewTenantApiKeyRepo(db),
		TenantWebhookRepository:         repository.NewTenantWebhookRepo(db),
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

	err = db.AutoMigrate(&entity.WhitelistDomain{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	err = db.AutoMigrate(&entity.PersonalIntegration{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	err = db.AutoMigrate(&entity.PersonalEmailProvider{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	err = db.AutoMigrate(&entity.TenantApiKey{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	err = db.AutoMigrate(&entity.TenantWebhook{})
	if err != nil {
		log.Print(err)
		panic(err)
	}

	return repositories
}
