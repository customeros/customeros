package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	AiLocationMappingRepository                 AiLocationMappingRepository
	AiPromptLogRepository                       AiPromptLogRepository
	AppKeyRepository                            AppKeyRepository
	PersonalIntegrationRepository               PersonalIntegrationRepository
	PersonalEmailProviderRepository             PersonalEmailProviderRepository
	TenantWebhookApiKeyRepository               TenantWebhookApiKeyRepository
	TenantWebhookRepository                     TenantWebhookRepository
	SlackChannelRepository                      SlackChannelRepository
	PostmarkApiKeyRepository                    PostmarkApiKeyRepository
	GoogleServiceAccountKeyRepository           GoogleServiceAccountKeyRepository
	CurrencyRateRepository                      CurrencyRateRepository
	EventBufferRepository                       EventBufferRepository
	TableViewDefinitionRepository               TableViewDefinitionRepository
	TrackingAllowedOriginRepository             TrackingAllowedOriginRepository
	ExternalAppKeysRepository                   ExternalAppKeysRepository
	EnrichDetailsBetterContactRepository        EnrichDetailsBetterContactRepository
	EnrichDetailsScrapInRepository              EnrichDetailsScrapInRepository
	EnrichDetailsBrandfetchRepository           EnrichDetailsBrandfetchRepository
	EnrichDetailsPrefilterTrackingRepository    EnrichDetailsPrefilterTrackingRepository
	EnrichDetailsTrackingRepository             EnrichDetailsTrackingRepository
	UserEmailImportPageTokenRepository          UserEmailImportStateRepository
	RawEmailRepository                          RawEmailRepository
	OAuthTokenRepository                        OAuthTokenRepository
	SlackSettingsRepository                     SlackSettingsRepository
	SlackChannelNotificationRepository          SlackChannelNotificationRepository
	WorkflowRepository                          WorkflowRepository
	IndustryMappingRepository                   IndustryMappingRepository
	TrackingRepository                          TrackingRepository
	TenantSettingsRepository                    TenantSettingsRepository
	TenantSettingsOpportunityStageRepository    TenantSettingsOpportunityStageRepository
	TenantSettingsMailboxRepository             TenantSettingsMailboxRepository
	TenantSettingsEmailExclusionRepository      TenantSettingsEmailExclusionRepository
	EmailLookupRepository                       EmailLookupRepository
	EmailTrackingRepository                     EmailTrackingRepository
	TenantRepository                            TenantRepository
	CacheIpDataRepository                       CacheIpDataRepository
	CacheIpHunterRepository                     CacheIpHunterRepository
	CacheEmailValidationRepository              CacheEmailValidationRepository
	CacheEmailValidationDomainRepository        CacheEmailValidationDomainRepository
	StatsApiCallsRepository                     StatsApiCallsRepository
	CosApiEnrichPersonTempResultRepository      CosApiEnrichPersonTempResultRepository
	OranizationWebsiteHostingPlatformRepository OrganizationWebsiteHostingPlatformRepository
	CustomerOsIdsRepository                     CustomerOsIdsRepository
	CacheEmailScrubbyRepository                 CacheEmailScrubbyRepository
	CacheEmailTrueinboxRepository               CacheEmailTrueinboxRepository
	CacheEmailEnrowRepository                   CacheEmailEnrowRepository
	EmailValidationRecordRepository             EmailValidationRecordRepository
	EmailValidationRequestBulkRepository        EmailValidationRequestBulkRepository
	ApiBillableEventRepository                  ApiBillableEventRepository
	MailStackDomainRepository                   MailStackDomainRepository
	BrowserConfigRepository                     BrowserConfigRepository
	BrowserAutomationRunRepository              BrowserAutomationRunRepository
	BrowserAutomationRunResultRepository        BrowserAutomationRunResultRepository
	EmailMessageRepository                      EmailMessageRepository
	UserWorkingScheduleRepository               UserWorkingScheduleRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		AiLocationMappingRepository:                 NewAiLocationMappingRepository(db),
		AiPromptLogRepository:                       NewAiPromptLogRepository(db),
		AppKeyRepository:                            NewAppKeyRepo(db),
		PersonalIntegrationRepository:               NewPersonalIntegrationsRepo(db),
		PersonalEmailProviderRepository:             NewPersonalEmailProviderRepository(db),
		TenantWebhookApiKeyRepository:               NewTenantWebhookApiKeyRepository(db),
		TenantWebhookRepository:                     NewTenantWebhookRepo(db),
		SlackChannelRepository:                      NewSlackChannelRepository(db),
		PostmarkApiKeyRepository:                    NewPostmarkApiKeyRepo(db),
		GoogleServiceAccountKeyRepository:           NewGoogleServiceAccountKeyRepository(db),
		CurrencyRateRepository:                      NewCurrencyRateRepository(db),
		EventBufferRepository:                       NewEventBufferRepository(db),
		TableViewDefinitionRepository:               NewTableViewDefinitionRepository(db),
		TrackingAllowedOriginRepository:             NewTrackingAllowedOriginRepository(db),
		ExternalAppKeysRepository:                   NewExternalAppKeysRepository(db),
		EnrichDetailsBetterContactRepository:        NewEnrichDetailsBetterContactRepository(db),
		EnrichDetailsScrapInRepository:              NewEnrichDetailsScrapInRepository(db),
		EnrichDetailsBrandfetchRepository:           NewEnrichDetailsBrandfetchRepository(db),
		EnrichDetailsPrefilterTrackingRepository:    NewEnrichDetailsPrefilterTrackingRepository(db),
		EnrichDetailsTrackingRepository:             NewEnrichDetailsTrackingRepository(db),
		UserEmailImportPageTokenRepository:          NewUserEmailImportStateRepository(db),
		RawEmailRepository:                          NewRawEmailRepository(db),
		OAuthTokenRepository:                        NewOAuthTokenRepository(db),
		SlackSettingsRepository:                     NewSlackSettingsRepository(db),
		SlackChannelNotificationRepository:          NewSlackChannelNotificationRepository(db),
		WorkflowRepository:                          NewWorkflowRepository(db),
		IndustryMappingRepository:                   NewIndustryMappingRepository(db),
		TrackingRepository:                          NewTrackingRepository(db),
		TenantSettingsRepository:                    NewTenantSettingsRepository(db),
		TenantSettingsOpportunityStageRepository:    NewTenantSettingsOpportunityStageRepository(db),
		TenantSettingsMailboxRepository:             NewTenantSettingsMailboxRepository(db),
		TenantSettingsEmailExclusionRepository:      NewEmailExclusionRepository(db),
		EmailLookupRepository:                       NewEmailLookupRepository(db),
		EmailTrackingRepository:                     NewEmailTrackingRepository(db),
		TenantRepository:                            NewTenantRepository(db),
		CacheIpDataRepository:                       NewCacheIpDataRepository(db),
		CacheIpHunterRepository:                     NewCacheIpHunterRepository(db),
		CacheEmailValidationRepository:              NewCacheEmailValidationRepository(db),
		CacheEmailValidationDomainRepository:        NewCacheEmailValidationDomainRepository(db),
		StatsApiCallsRepository:                     NewStatsApiCallsRepository(db),
		CosApiEnrichPersonTempResultRepository:      NewCosApiEnrichPersonTempResultRepository(db),
		OranizationWebsiteHostingPlatformRepository: NewOrganizationWebsiteHostingPlatformRepository(db),
		CustomerOsIdsRepository:                     NewCustomerOsIdsRepository(db),
		CacheEmailScrubbyRepository:                 NewCacheEmailScrubbyRepository(db),
		CacheEmailTrueinboxRepository:               NewCacheEmailTrueinboxRepository(db),
		CacheEmailEnrowRepository:                   NewCacheEmailEnrowRepository(db),
		EmailValidationRecordRepository:             NewEmailValidationRecordRepository(db),
		EmailValidationRequestBulkRepository:        NewEmailValidationRequestBulkRepository(db),
		ApiBillableEventRepository:                  NewApiBillableEventRepository(db),
		MailStackDomainRepository:                   NewMailStackDomainRepository(db),
		BrowserConfigRepository:                     NewBrowserConfigRepository(db),
		BrowserAutomationRunRepository:              NewBrowserAutomationRunRepository(db),
		BrowserAutomationRunResultRepository:        NewBrowserAutomationRunResultRepository(db),
		EmailMessageRepository:                      NewEmailMessageRepository(db),
		UserWorkingScheduleRepository:               NewUserWorkingScheduleRepository(db),
	}

	return repositories
}

func (r *Repositories) Migration(db *gorm.DB) {

	//err = db.AutoMigrate(&entity.AppKey{})
	//if err != nil {
	//	panic(err)
	//}

	err := db.AutoMigrate(
		&entity.Tenant{},
		&entity.AiLocationMapping{},
		&entity.AiPromptLog{},
		&entity.PersonalIntegration{},
		&entity.PersonalEmailProvider{},
		&entity.TenantWebhookApiKey{},
		&entity.TenantWebhook{},
		&entity.SlackChannel{},
		&entity.PostmarkApiKey{},
		&entity.GoogleServiceAccountKey{},
		&entity.CurrencyRate{},
		&entity.EventBuffer{},
		&entity.TableViewDefinition{},
		&entity.TrackingAllowedOrigin{},
		&entity.TenantSettingsEmailExclusion{},
		&entity.ExternalAppKeys{},
		&entity.EnrichDetailsBetterContact{},
		&entity.EnrichDetailsScrapIn{},
		&entity.EnrichDetailsBrandfetch{},
		&entity.UserEmailImportState{},
		&entity.UserEmailImportStateHistory{},
		&entity.RawEmail{},
		&entity.OAuthTokenEntity{},
		&entity.SlackSettingsEntity{},
		&entity.SlackChannelNotification{},
		&entity.Workflow{},
		&entity.IndustryMapping{},
		&entity.Tracking{},
		&entity.EnrichDetailsPreFilterTracking{},
		&entity.EnrichDetailsTracking{},
		&entity.TenantSettings{},
		&entity.TenantSettingsOpportunityStage{},
		&entity.TenantSettingsMailbox{},
		&entity.EmailLookup{},
		&entity.EmailTracking{},
		&entity.CacheIpData{},
		&entity.CacheIpHunter{},
		&entity.CacheEmailValidation{},
		&entity.CacheEmailValidationDomain{},
		&entity.CacheEmailScrubby{},
		&entity.StatsApiCalls{},
		&entity.CosApiEnrichPersonTempResult{},
		&entity.OrganizationWebsiteHostingPlatform{},
		&entity.CustomerOsIds{},
		&entity.CacheEmailTrueinbox{},
		&entity.CacheEmailEnrow{},
		&entity.EmailValidationRecord{},
		&entity.EmailValidationRequestBulk{},
		&entity.ApiBillableEvent{},
		&entity.MailStackDomain{},
		&entity.EmailMessage{},
		&entity.UserWorkingSchedule{})
	if err != nil {
		panic(err)
	}
}
