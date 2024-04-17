package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	AppKeyRepository                  AppKeyRepository
	PersonalIntegrationRepository     PersonalIntegrationRepository
	AiPromptLogRepository             AiPromptLogRepository
	WhitelistDomainRepository         WhitelistDomainRepository
	PersonalEmailProviderRepository   PersonalEmailProviderRepository
	TenantWebhookApiKeyRepository     TenantWebhookApiKeyRepository
	TenantWebhookRepository           TenantWebhookRepository
	SlackChannelRepository            SlackChannelRepository
	PostmarkApiKeyRepository          PostmarkApiKeyRepository
	GoogleServiceAccountKeyRepository GoogleServiceAccountKeyRepository
	CurrencyRateRepository            CurrencyRateRepository
	EventBufferRepository             EventBufferRepository
	TableViewDefinitionRepository     TableViewDefinitionRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		AppKeyRepository:                  NewAppKeyRepo(db),
		PersonalIntegrationRepository:     NewPersonalIntegrationsRepo(db),
		AiPromptLogRepository:             NewAiPromptLogRepository(db),
		WhitelistDomainRepository:         NewWhitelistDomainRepository(db),
		PersonalEmailProviderRepository:   NewPersonalEmailProviderRepository(db),
		TenantWebhookApiKeyRepository:     NewTenantWebhookApiKeyRepo(db),
		TenantWebhookRepository:           NewTenantWebhookRepo(db),
		SlackChannelRepository:            NewSlackChannelRepository(db),
		PostmarkApiKeyRepository:          NewPostmarkApiKeyRepo(db),
		GoogleServiceAccountKeyRepository: NewGoogleServiceAccountKeyRepository(db),
		CurrencyRateRepository:            NewCurrencyRateRepository(db),
		EventBufferRepository:             NewEventBufferRepository(db),
	}

	return repositories
}

func (r *Repositories) Migration(db *gorm.DB) {

	var err error

	err = db.AutoMigrate(&entity.AppKey{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.AiPromptLog{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.WhitelistDomain{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.PersonalIntegration{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.PersonalEmailProvider{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.TenantWebhookApiKey{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.TenantWebhook{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.SlackChannel{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.PostmarkApiKey{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.GoogleServiceAccountKey{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.CurrencyRate{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.EventBuffer{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.TableViewDefinition{})
	if err != nil {
		panic(err)
	}

}
