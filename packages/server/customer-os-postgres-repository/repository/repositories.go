package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	AppKeyRepository                   AppKeyRepository
	PersonalIntegrationRepository      PersonalIntegrationRepository
	AiPromptLogRepository              AiPromptLogRepository
	PersonalEmailProviderRepository    PersonalEmailProviderRepository
	TenantWebhookApiKeyRepository      TenantWebhookApiKeyRepository
	TenantWebhookRepository            TenantWebhookRepository
	SlackChannelRepository             SlackChannelRepository
	PostmarkApiKeyRepository           PostmarkApiKeyRepository
	GoogleServiceAccountKeyRepository  GoogleServiceAccountKeyRepository
	CurrencyRateRepository             CurrencyRateRepository
	EventBufferRepository              EventBufferRepository
	TableViewDefinitionRepository      TableViewDefinitionRepository
	TrackingAllowedOriginRepository    TrackingAllowedOriginRepository
	TechLimitRepository                TechLimitRepository
	EmailExclusionRepository           EmailExclusionRepository
	ExternalAppKeysRepository          ExternalAppKeysRepository
	EnrichDetailsBetterContactRepository EnrichDetailsBetterContactRepository
	UserEmailImportPageTokenRepository UserEmailImportStateRepository
	RawEmailRepository                 RawEmailRepository
	OAuthTokenRepository               OAuthTokenRepository
	SlackSettingsRepository            SlackSettingsRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		AppKeyRepository:                   NewAppKeyRepo(db),
		PersonalIntegrationRepository:      NewPersonalIntegrationsRepo(db),
		AiPromptLogRepository:              NewAiPromptLogRepository(db),
		PersonalEmailProviderRepository:    NewPersonalEmailProviderRepository(db),
		TenantWebhookApiKeyRepository:      NewTenantWebhookApiKeyRepo(db),
		TenantWebhookRepository:            NewTenantWebhookRepo(db),
		SlackChannelRepository:             NewSlackChannelRepository(db),
		PostmarkApiKeyRepository:           NewPostmarkApiKeyRepo(db),
		GoogleServiceAccountKeyRepository:  NewGoogleServiceAccountKeyRepository(db),
		CurrencyRateRepository:             NewCurrencyRateRepository(db),
		EventBufferRepository:              NewEventBufferRepository(db),
		TableViewDefinitionRepository:      NewTableViewDefinitionRepository(db),
		TrackingAllowedOriginRepository:    NewTrackingAllowedOriginRepository(db),
		TechLimitRepository:                NewTechLimitRepository(db),
		EmailExclusionRepository:           NewEmailExclusionRepository(db),
		ExternalAppKeysRepository:          NewExternalAppKeysRepository(db),
		EnrichDetailsBetterContactRepository: NewEnrichDetailsBetterContactRepository(db),
		UserEmailImportPageTokenRepository: NewUserEmailImportStateRepository(db),
		RawEmailRepository:                 NewRawEmailRepository(db),
		OAuthTokenRepository:               NewOAuthTokenRepository(db),
		SlackSettingsRepository:            NewSlackSettingsRepository(db),
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

	err = db.AutoMigrate(&entity.TechLimit{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.TrackingAllowedOrigin{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.EmailExclusion{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.ExternalAppKeys{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.EnrichDetailsBetterContact{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.UserEmailImportState{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.UserEmailImportStateHistory{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.RawEmail{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.OAuthTokenEntity{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.SlackSettingsEntity{})
	if err != nil {
		panic(err)
	}

}
