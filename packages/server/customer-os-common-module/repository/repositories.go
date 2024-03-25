package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	AppKeyRepository                  repository.AppKeyRepository
	PersonalIntegrationRepository     repository.PersonalIntegrationRepository
	AiPromptLogRepository             repository.AiPromptLogRepository
	WhitelistDomainRepository         repository.WhitelistDomainRepository
	PersonalEmailProviderRepository   repository.PersonalEmailProviderRepository
	UserRepository                    neo4jrepo.UserRepository
	TenantRepository                  neo4jrepo.TenantRepository
	StateRepository                   neo4jrepo.StateRepository
	TenantWebhookApiKeyRepository     repository.TenantWebhookApiKeyRepository
	TenantWebhookRepository           repository.TenantWebhookRepository
	SlackChannelRepository            repository.SlackChannelRepository
	PostmarkApiKeyRepository          repository.PostmarkApiKeyRepository
	GoogleServiceAccountKeyRepository repository.GoogleServiceAccountKeyRepository
	CurrencyRateRepository            repository.CurrencyRateRepository
	EventBufferRepository             repository.EventBufferRepository
}

func InitRepositories(db *gorm.DB, driver *neo4j.DriverWithContext) *Repositories {
	repositories := &Repositories{
		AppKeyRepository:                  repository.NewAppKeyRepo(db),
		PersonalIntegrationRepository:     repository.NewPersonalIntegrationsRepo(db),
		AiPromptLogRepository:             repository.NewAiPromptLogRepository(db),
		WhitelistDomainRepository:         repository.NewWhitelistDomainRepository(db),
		PersonalEmailProviderRepository:   repository.NewPersonalEmailProviderRepository(db),
		UserRepository:                    neo4jrepo.NewUserRepository(driver),
		TenantRepository:                  neo4jrepo.NewTenantRepository(driver),
		StateRepository:                   neo4jrepo.NewStateRepository(driver),
		TenantWebhookApiKeyRepository:     repository.NewTenantWebhookApiKeyRepo(db),
		TenantWebhookRepository:           repository.NewTenantWebhookRepo(db),
		SlackChannelRepository:            repository.NewSlackChannelRepository(db),
		PostmarkApiKeyRepository:          repository.NewPostmarkApiKeyRepo(db),
		GoogleServiceAccountKeyRepository: repository.NewGoogleServiceAccountKeyRepository(db),
		CurrencyRateRepository:            repository.NewCurrencyRateRepository(db),
		EventBufferRepository:             repository.NewEventBufferRepository(db),
	}

	return repositories
}

func Migration(db *gorm.DB) {

	var err error

	err = db.AutoMigrate(&commonEntity.AppKey{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.AiPromptLog{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.WhitelistDomain{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.PersonalIntegration{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.PersonalEmailProvider{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.TenantWebhookApiKey{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.TenantWebhook{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.SlackChannel{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.PostmarkApiKey{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.GoogleServiceAccountKey{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.CurrencyRate{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&commonEntity.EventBuffer{})
	if err != nil {
		panic(err)
	}

}
