package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type Repositories struct {
	Db *gorm.DB

	AiLocationMappingRepository                 AiLocationMappingRepository
	AiPromptLogRepository                       AiPromptLogRepository
	ApiBillableEventRepository                  ApiBillableEventRepository
	AppKeyRepository                            AppKeyRepository
	BrowserAutomationRunRepository              BrowserAutomationRunRepository
	BrowserAutomationRunResultRepository        BrowserAutomationRunResultRepository
	BrowserConfigRepository                     BrowserConfigRepository
	CacheEmailEnrowRepository                   CacheEmailEnrowRepository
	CacheEmailScrubbyRepository                 CacheEmailScrubbyRepository
	CacheEmailTrueinboxRepository               CacheEmailTrueinboxRepository
	CacheEmailValidationDomainRepository        CacheEmailValidationDomainRepository
	CacheEmailValidationRepository              CacheEmailValidationRepository
	CacheIpDataRepository                       CacheIpDataRepository
	CacheIpHunterRepository                     CacheIpHunterRepository
	CommonRepository                            CommonRepository
	CosApiEnrichPersonTempResultRepository      CosApiEnrichPersonTempResultRepository
	CurrencyRateRepository                      CurrencyRateRepository
	CustomerOsIdsRepository                     CustomerOsIdsRepository
	EmailLookupRepository                       EmailLookupRepository
	EmailMessageRepository                      EmailMessageRepository
	EmailTrackingRepository                     EmailTrackingRepository
	EmailValidationRecordRepository             EmailValidationRecordRepository
	EmailValidationRequestBulkRepository        EmailValidationRequestBulkRepository
	EnrichDetailsBetterContactRepository        EnrichDetailsBetterContactRepository
	EnrichDetailsBrandfetchRepository           EnrichDetailsBrandfetchRepository
	EnrichDetailsPrefilterTrackingRepository    EnrichDetailsPrefilterTrackingRepository
	EnrichDetailsScrapInRepository              EnrichDetailsScrapInRepository
	EnrichDetailsTrackingRepository             EnrichDetailsTrackingRepository
	EventBufferRepository                       EventBufferRepository
	ExternalAppKeysRepository                   ExternalAppKeysRepository
	FlowActionRegistryRepository                FlowActionRegistryRepository
	FlowListenerRegistryRepository              FlowListenerRegistryRepository
	FlowTransitionsRegistryRepository           FlowTransitionsRegistryRepository
	FlowWebhooksRepository                      FlowWebhooksRepository
	GoogleServiceAccountKeyRepository           GoogleServiceAccountKeyRepository
	IndustryMappingRepository                   IndustryMappingRepository
	MailStackDomainRepository                   MailStackDomainRepository
	MailstackBuyRequestRepository               MailstackBuyRequestRepository
	OAuthTokenRepository                        OAuthTokenRepository
	OranizationWebsiteHostingPlatformRepository OrganizationWebsiteHostingPlatformRepository
	PersonalEmailProviderRepository             PersonalEmailProviderRepository
	PersonalIntegrationRepository               PersonalIntegrationRepository
	PostmarkApiKeyRepository                    PostmarkApiKeyRepository
	RawEmailRepository                          RawEmailRepository
	SlackChannelNotificationRepository          SlackChannelNotificationRepository
	SlackChannelRepository                      SlackChannelRepository
	SlackSettingsRepository                     SlackSettingsRepository
	StatsApiCallsRepository                     StatsApiCallsRepository
	TableViewDefinitionRepository               TableViewDefinitionRepository
	TenantRepository                            TenantRepository
	TenantSettingsEmailExclusionRepository      TenantSettingsEmailExclusionRepository
	TenantSettingsMailboxRepository             TenantSettingsMailboxRepository
	TenantSettingsOpportunityStageRepository    TenantSettingsOpportunityStageRepository
	TenantSettingsRepository                    TenantSettingsRepository
	TenantWebhookApiKeyRepository               TenantWebhookApiKeyRepository
	TenantWebhookRepository                     TenantWebhookRepository
	TrackingAllowedOriginRepository             TrackingAllowedOriginRepository
	TrackingRepository                          TrackingRepository
	UserEmailImportPageTokenRepository          UserEmailImportStateRepository
	UserWorkingScheduleRepository               UserWorkingScheduleRepository
	WorkflowRepository                          WorkflowRepository
	GlobalOrganizationRepository                GlobalOrganizationRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		Db: db,

		AiLocationMappingRepository:                 NewAiLocationMappingRepository(db),
		AiPromptLogRepository:                       NewAiPromptLogRepository(db),
		ApiBillableEventRepository:                  NewApiBillableEventRepository(db),
		AppKeyRepository:                            NewAppKeyRepo(db),
		BrowserAutomationRunRepository:              NewBrowserAutomationRunRepository(db),
		BrowserAutomationRunResultRepository:        NewBrowserAutomationRunResultRepository(db),
		BrowserConfigRepository:                     NewBrowserConfigRepository(db),
		CacheEmailEnrowRepository:                   NewCacheEmailEnrowRepository(db),
		CacheEmailScrubbyRepository:                 NewCacheEmailScrubbyRepository(db),
		CacheEmailTrueinboxRepository:               NewCacheEmailTrueinboxRepository(db),
		CacheEmailValidationDomainRepository:        NewCacheEmailValidationDomainRepository(db),
		CacheEmailValidationRepository:              NewCacheEmailValidationRepository(db),
		CacheIpDataRepository:                       NewCacheIpDataRepository(db),
		CacheIpHunterRepository:                     NewCacheIpHunterRepository(db),
		CommonRepository:                            NewCommonRepository(db),
		CosApiEnrichPersonTempResultRepository:      NewCosApiEnrichPersonTempResultRepository(db),
		CurrencyRateRepository:                      NewCurrencyRateRepository(db),
		CustomerOsIdsRepository:                     NewCustomerOsIdsRepository(db),
		EmailLookupRepository:                       NewEmailLookupRepository(db),
		EmailMessageRepository:                      NewEmailMessageRepository(db),
		EmailTrackingRepository:                     NewEmailTrackingRepository(db),
		EmailValidationRecordRepository:             NewEmailValidationRecordRepository(db),
		EmailValidationRequestBulkRepository:        NewEmailValidationRequestBulkRepository(db),
		EnrichDetailsBetterContactRepository:        NewEnrichDetailsBetterContactRepository(db),
		EnrichDetailsBrandfetchRepository:           NewEnrichDetailsBrandfetchRepository(db),
		EnrichDetailsPrefilterTrackingRepository:    NewEnrichDetailsPrefilterTrackingRepository(db),
		EnrichDetailsScrapInRepository:              NewEnrichDetailsScrapInRepository(db),
		EnrichDetailsTrackingRepository:             NewEnrichDetailsTrackingRepository(db),
		EventBufferRepository:                       NewEventBufferRepository(db),
		ExternalAppKeysRepository:                   NewExternalAppKeysRepository(db),
		FlowActionRegistryRepository:                NewFlowActionRegistryRepository(db),
		FlowListenerRegistryRepository:              NewFlowListenerRegistryRepository(db),
		FlowTransitionsRegistryRepository:           NewFlowTransitionsRegistryRepository(db),
		FlowWebhooksRepository:                      NewFlowWebhooksRepository(db),
		GoogleServiceAccountKeyRepository:           NewGoogleServiceAccountKeyRepository(db),
		IndustryMappingRepository:                   NewIndustryMappingRepository(db),
		MailStackDomainRepository:                   NewMailStackDomainRepository(db),
		MailstackBuyRequestRepository:               NewMailstackBuyRequestRepository(db),
		OAuthTokenRepository:                        NewOAuthTokenRepository(db),
		OranizationWebsiteHostingPlatformRepository: NewOrganizationWebsiteHostingPlatformRepository(db),
		PersonalEmailProviderRepository:             NewPersonalEmailProviderRepository(db),
		PersonalIntegrationRepository:               NewPersonalIntegrationsRepo(db),
		PostmarkApiKeyRepository:                    NewPostmarkApiKeyRepo(db),
		RawEmailRepository:                          NewRawEmailRepository(db),
		SlackChannelNotificationRepository:          NewSlackChannelNotificationRepository(db),
		SlackChannelRepository:                      NewSlackChannelRepository(db),
		SlackSettingsRepository:                     NewSlackSettingsRepository(db),
		StatsApiCallsRepository:                     NewStatsApiCallsRepository(db),
		TableViewDefinitionRepository:               NewTableViewDefinitionRepository(db),
		TenantRepository:                            NewTenantRepository(db),
		TenantSettingsEmailExclusionRepository:      NewEmailExclusionRepository(db),
		TenantSettingsMailboxRepository:             NewTenantSettingsMailboxRepository(db),
		TenantSettingsOpportunityStageRepository:    NewTenantSettingsOpportunityStageRepository(db),
		TenantSettingsRepository:                    NewTenantSettingsRepository(db),
		TenantWebhookApiKeyRepository:               NewTenantWebhookApiKeyRepository(db),
		TenantWebhookRepository:                     NewTenantWebhookRepo(db),
		TrackingAllowedOriginRepository:             NewTrackingAllowedOriginRepository(db),
		TrackingRepository:                          NewTrackingRepository(db),
		UserEmailImportPageTokenRepository:          NewUserEmailImportStateRepository(db),
		UserWorkingScheduleRepository:               NewUserWorkingScheduleRepository(db),
		WorkflowRepository:                          NewWorkflowRepository(db),
		GlobalOrganizationRepository:                NewGlobalOrganizationRepository(db),
	}

	return repositories
}

func (r *Repositories) Migration(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.AiLocationMapping{},
		&entity.AiPromptLog{},
		&entity.ApiBillableEvent{},
		&entity.CacheEmailEnrow{},
		&entity.CacheEmailScrubby{},
		&entity.CacheEmailTrueinbox{},
		&entity.CacheEmailValidation{},
		&entity.CacheEmailValidationDomain{},
		&entity.CacheIpData{},
		&entity.CacheIpHunter{},
		&entity.CosApiEnrichPersonTempResult{},
		&entity.CurrencyRate{},
		&entity.CustomerOsIds{},
		&entity.DMARCMonitoring{},
		&entity.EmailLookup{},
		&entity.EmailMessage{},
		&entity.EmailTracking{},
		&entity.EmailValidationRecord{},
		&entity.EmailValidationRequestBulk{},
		&entity.EnrichDetailsBetterContact{},
		&entity.EnrichDetailsBrandfetch{},
		&entity.EnrichDetailsPreFilterTracking{},
		&entity.EnrichDetailsScrapIn{},
		&entity.EnrichDetailsTracking{},
		&entity.EventBuffer{},
		&entity.ExternalAppKeys{},
		&entity.FlowActionRegistry{},
		&entity.FlowListenerRegistry{},
		&entity.FlowTransitionsRegistry{},
		&entity.FlowWebhooks{},
		&entity.GoogleServiceAccountKey{},
		&entity.IndustryMapping{},
		&entity.MailStackDomain{},
		&entity.MailstackBuyRequest{},
		&entity.MailstackBuyRequestDomain{},
		&entity.MailstackReputationEntity{},
		&entity.OAuthTokenEntity{},
		&entity.OrganizationWebsiteHostingPlatform{},
		&entity.PersonalEmailProvider{},
		&entity.PersonalIntegration{},
		&entity.PostmarkApiKey{},
		&entity.RawEmail{},
		&entity.SlackChannel{},
		&entity.SlackChannelNotification{},
		&entity.SlackSettingsEntity{},
		&entity.StatsApiCalls{},
		&entity.TableViewDefinition{},
		&entity.Tenant{},
		&entity.TenantSettings{},
		&entity.TenantSettingsEmailExclusion{},
		&entity.TenantSettingsMailbox{},
		&entity.TenantSettingsOpportunityStage{},
		&entity.TenantWebhook{},
		&entity.TenantWebhookApiKey{},
		&entity.Tracking{},
		&entity.TrackingAllowedOrigin{},
		&entity.UserEmailImportState{},
		&entity.UserEmailImportStateHistory{},
		&entity.UserWorkingSchedule{},
		&entity.Workflow{},
		&entity.GlobalOrganization{},
	)
	if err != nil {
		panic(err)
	}
}

func (r *Repositories) InitData(ctx context.Context, postgresRepos *Repositories) {
	err := r.FlowActionRegistryRepository.InitializeActions(ctx)
	if err != nil {
		panic(err)
	}
	err = r.FlowListenerRegistryRepository.InitializeFlowListenerEvents(ctx)
	if err != nil {
		panic(err)
	}

	err = r.FlowTransitionsRegistryRepository.InitializeFlowTransitions(ctx)
	if err != nil {
		panic(err)
	}
}
