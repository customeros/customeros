package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type Repositories struct {
	Db      *gorm.DB
	AsyncDb *gorm.DB

	AgentCapabilityRegistryRepository            AgentCapabilityRegistryRepository
	AgentsRepository                             AgentsRepository
	AgentExecutionRepository                     AgentExecutionRepository
	AgentRegistryRepository                      AgentRegistryRepository
	AiLocationMappingRepository                  AiLocationMappingRepository
	AiPromptLogRepository                        AiPromptLogRepository
	ApiBillableEventRepository                   ApiBillableEventRepository
	AppKeyRepository                             AppKeyRepository
	BrowserAutomationRunRepository               BrowserAutomationRunRepository
	BrowserAutomationRunResultRepository         BrowserAutomationRunResultRepository
	BrowserConfigRepository                      BrowserConfigRepository
	CacheEmailEnrowRepository                    CacheEmailEnrowRepository
	CacheEmailScrubbyRepository                  CacheEmailScrubbyRepository
	CacheEmailTrueinboxRepository                CacheEmailTrueinboxRepository
	CacheEmailValidationDomainRepository         CacheEmailValidationDomainRepository
	CacheEmailValidationRepository               CacheEmailValidationRepository
	CacheIpDataRepository                        CacheIpDataRepository
	CacheIpHunterRepository                      CacheIpHunterRepository
	CacheIPIdentifyRepository                    CacheIPIdentifyRepository
	CommonRepository                             CommonRepository
	CosApiEnrichPersonTempResultRepository       CosApiEnrichPersonTempResultRepository
	CurrencyRateRepository                       CurrencyRateRepository
	CustomerOsIdsRepository                      CustomerOsIdsRepository
	EmailLookupRepository                        EmailLookupRepository
	EmailMessageRepository                       EmailMessageRepository
	EmailTrackingRepository                      EmailTrackingRepository
	EmailValidationRecordRepository              EmailValidationRecordRepository
	EmailValidationRequestBulkRepository         EmailValidationRequestBulkRepository
	EnrichDetailsBetterContactRepository         EnrichDetailsBetterContactRepository
	EnrichDetailsBrandfetchRepository            EnrichDetailsBrandfetchRepository
	EnrichDetailsPrefilterTrackingRepository     EnrichDetailsPrefilterTrackingRepository
	EnrichDetailsScrapInRepository               EnrichDetailsScrapInRepository
	EnrichDetailsTrackingRepository              EnrichDetailsTrackingRepository
	EventBufferRepository                        EventBufferRepository
	ExternalAppKeysRepository                    ExternalAppKeysRepository
	FlowsRepository                              FlowsRepository
	FlowEdgeRepository                           FlowEdgeRepository
	FlowExecutionRepository                      FlowExecutionRepository
	FlowNodeRepository                           FlowNodeRepository
	FlowTransitionsRegistryRepository            FlowTransitionsRegistryRepository
	GlobalOrganizationRepository                 GlobalOrganizationRepository
	GlobalOrganizationWebsiteToProcessRepository GlobalOrganizationWebsiteToProcessRepository
	GoogleServiceAccountKeyRepository            GoogleServiceAccountKeyRepository
	MailStackDomainRepository                    MailStackDomainRepository
	MailstackBuyRequestRepository                MailstackBuyRequestRepository
	MagicLinkRepository                          MagicLinkRepository
	OAuthTokenRepository                         OAuthTokenRepository
	OranizationWebsiteHostingPlatformRepository  OrganizationWebsiteHostingPlatformRepository
	PersonalEmailProviderRepository              PersonalEmailProviderRepository
	PersonalIntegrationRepository                PersonalIntegrationRepository
	PostmarkApiKeyRepository                     PostmarkApiKeyRepository
	RawEmailRepository                           RawEmailRepository
	SlackChannelNotificationRepository           SlackChannelNotificationRepository
	SlackChannelRepository                       SlackChannelRepository
	SlackSettingsRepository                      SlackSettingsRepository
	StatsApiCallsRepository                      StatsApiCallsRepository
	TableViewDefinitionRepository                TableViewDefinitionRepository
	TenantRepository                             TenantRepository
	TenantSettingsEmailExclusionRepository       TenantSettingsEmailExclusionRepository
	TenantSettingsMailboxRepository              TenantSettingsMailboxRepository
	TenantSettingsOpportunityStageRepository     TenantSettingsOpportunityStageRepository
	TenantSettingsRepository                     TenantSettingsRepository
	TenantWebhookApiKeyRepository                TenantWebhookApiKeyRepository
	TenantWebhookRepository                      TenantWebhookRepository
	TrackingAllowedOriginRepository              TrackingAllowedOriginRepository
	UserEmailImportPageTokenRepository           UserEmailImportStateRepository
	UserWorkingScheduleRepository                UserWorkingScheduleRepository
	WebhooksRepository                           WebhooksRepository
	WebSessionRepository                         WebSessionRepository
	WebTrackerEventsRepository                   WebTrackerEventsRepository
}

func InitRepositories(postgresDB *config.PostgresDB) *Repositories {
	repositories := &Repositories{
		Db:      postgresDB.GormDB,
		AsyncDb: postgresDB.AsyncGormDB,

		OAuthTokenRepository:               NewOAuthTokenRepository(postgresDB.AsyncGormDB),
		RawEmailRepository:                 NewRawEmailRepository(postgresDB.AsyncGormDB),
		GoogleServiceAccountKeyRepository:  NewGoogleServiceAccountKeyRepository(postgresDB.AsyncGormDB),
		UserEmailImportPageTokenRepository: NewUserEmailImportStateRepository(postgresDB.AsyncGormDB),

		AgentCapabilityRegistryRepository:            NewAgentCapabilityRegistryRepository(postgresDB.GormDB),
		AgentsRepository:                             NewAgentsRepository(postgresDB.GormDB),
		AgentExecutionRepository:                     NewAgentExecutionRepository(postgresDB.GormDB),
		AgentRegistryRepository:                      NewAgentRegistryRepository(postgresDB.GormDB),
		AiLocationMappingRepository:                  NewAiLocationMappingRepository(postgresDB.GormDB),
		AiPromptLogRepository:                        NewAiPromptLogRepository(postgresDB.GormDB),
		ApiBillableEventRepository:                   NewApiBillableEventRepository(postgresDB.GormDB),
		AppKeyRepository:                             NewAppKeyRepo(postgresDB.GormDB),
		BrowserAutomationRunRepository:               NewBrowserAutomationRunRepository(postgresDB.GormDB),
		BrowserAutomationRunResultRepository:         NewBrowserAutomationRunResultRepository(postgresDB.GormDB),
		BrowserConfigRepository:                      NewBrowserConfigRepository(postgresDB.GormDB),
		CacheEmailEnrowRepository:                    NewCacheEmailEnrowRepository(postgresDB.GormDB),
		CacheEmailScrubbyRepository:                  NewCacheEmailScrubbyRepository(postgresDB.GormDB),
		CacheEmailTrueinboxRepository:                NewCacheEmailTrueinboxRepository(postgresDB.GormDB),
		CacheEmailValidationDomainRepository:         NewCacheEmailValidationDomainRepository(postgresDB.GormDB),
		CacheEmailValidationRepository:               NewCacheEmailValidationRepository(postgresDB.GormDB),
		CacheIpDataRepository:                        NewCacheIpDataRepository(postgresDB.GormDB),
		CacheIpHunterRepository:                      NewCacheIpHunterRepository(postgresDB.GormDB),
		CacheIPIdentifyRepository:                    NewCacheIPIdentifyRepository(postgresDB.GormDB),
		CommonRepository:                             NewCommonRepository(postgresDB),
		CosApiEnrichPersonTempResultRepository:       NewCosApiEnrichPersonTempResultRepository(postgresDB.GormDB),
		CurrencyRateRepository:                       NewCurrencyRateRepository(postgresDB.GormDB),
		CustomerOsIdsRepository:                      NewCustomerOsIdsRepository(postgresDB.GormDB),
		EmailLookupRepository:                        NewEmailLookupRepository(postgresDB.GormDB),
		EmailMessageRepository:                       NewEmailMessageRepository(postgresDB.GormDB),
		EmailTrackingRepository:                      NewEmailTrackingRepository(postgresDB.GormDB),
		EmailValidationRecordRepository:              NewEmailValidationRecordRepository(postgresDB.GormDB),
		EmailValidationRequestBulkRepository:         NewEmailValidationRequestBulkRepository(postgresDB.GormDB),
		EnrichDetailsBetterContactRepository:         NewEnrichDetailsBetterContactRepository(postgresDB.GormDB),
		EnrichDetailsBrandfetchRepository:            NewEnrichDetailsBrandfetchRepository(postgresDB.GormDB),
		EnrichDetailsPrefilterTrackingRepository:     NewEnrichDetailsPrefilterTrackingRepository(postgresDB.GormDB),
		EnrichDetailsScrapInRepository:               NewEnrichDetailsScrapInRepository(postgresDB.GormDB),
		EnrichDetailsTrackingRepository:              NewEnrichDetailsTrackingRepository(postgresDB.GormDB),
		EventBufferRepository:                        NewEventBufferRepository(postgresDB.GormDB),
		ExternalAppKeysRepository:                    NewExternalAppKeysRepository(postgresDB.GormDB),
		FlowsRepository:                              NewFlowsRepository(postgresDB.GormDB),
		FlowEdgeRepository:                           NewFlowEdgeRepository(postgresDB.GormDB),
		FlowExecutionRepository:                      NewFlowExecutionRepository(postgresDB.GormDB),
		FlowNodeRepository:                           NewFlowNodeRepository(postgresDB.GormDB),
		FlowTransitionsRegistryRepository:            NewFlowTransitionsRegistryRepository(postgresDB.GormDB),
		GlobalOrganizationRepository:                 NewGlobalOrganizationRepository(postgresDB.GormDB),
		GlobalOrganizationWebsiteToProcessRepository: NewGlobalOrganizationWebsiteToProcessRepository(postgresDB.GormDB),
		MailStackDomainRepository:                    NewMailStackDomainRepository(postgresDB.GormDB),
		MailstackBuyRequestRepository:                NewMailstackBuyRequestRepository(postgresDB.GormDB),
		MagicLinkRepository:                          NewMagicLinkRepository(postgresDB.GormDB),
		OranizationWebsiteHostingPlatformRepository:  NewOrganizationWebsiteHostingPlatformRepository(postgresDB.GormDB),
		PersonalEmailProviderRepository:              NewPersonalEmailProviderRepository(postgresDB.GormDB),
		PersonalIntegrationRepository:                NewPersonalIntegrationsRepo(postgresDB.GormDB),
		PostmarkApiKeyRepository:                     NewPostmarkApiKeyRepo(postgresDB.GormDB),
		SlackChannelNotificationRepository:           NewSlackChannelNotificationRepository(postgresDB.GormDB),
		SlackChannelRepository:                       NewSlackChannelRepository(postgresDB.GormDB),
		SlackSettingsRepository:                      NewSlackSettingsRepository(postgresDB.GormDB),
		StatsApiCallsRepository:                      NewStatsApiCallsRepository(postgresDB.GormDB),
		TableViewDefinitionRepository:                NewTableViewDefinitionRepository(postgresDB.GormDB),
		TenantRepository:                             NewTenantRepository(postgresDB.GormDB),
		TenantSettingsEmailExclusionRepository:       NewEmailExclusionRepository(postgresDB.GormDB),
		TenantSettingsMailboxRepository:              NewTenantSettingsMailboxRepository(postgresDB.GormDB),
		TenantSettingsOpportunityStageRepository:     NewTenantSettingsOpportunityStageRepository(postgresDB.GormDB),
		TenantSettingsRepository:                     NewTenantSettingsRepository(postgresDB.GormDB),
		TenantWebhookApiKeyRepository:                NewTenantWebhookApiKeyRepository(postgresDB.GormDB),
		TenantWebhookRepository:                      NewTenantWebhookRepo(postgresDB.GormDB),
		TrackingAllowedOriginRepository:              NewTrackingAllowedOriginRepository(postgresDB.GormDB),
		UserWorkingScheduleRepository:                NewUserWorkingScheduleRepository(postgresDB.GormDB),
		WebhooksRepository:                           NewWebhooksRepository(postgresDB.GormDB),
		WebSessionRepository:                         NewWebSessionRepository(postgresDB.GormDB),
		WebTrackerEventsRepository:                   NewWebTrackerEventsRepository(postgresDB.GormDB),
	}

	return repositories
}

func (r *Repositories) Migration(postgresDB *config.PostgresDB) {
	err := postgresDB.GormDB.AutoMigrate(
		&postgres_entity.AgentRegistry{},
		&postgres_entity.AgentExecution{},
		&postgres_entity.Agents{},
		&postgres_entity.AiLocationMapping{},
		&postgres_entity.AiPromptLog{},
		&postgres_entity.ApiBillableEvent{},
		&postgres_entity.CacheEmailEnrow{},
		&postgres_entity.CacheEmailScrubby{},
		&postgres_entity.CacheEmailTrueinbox{},
		&postgres_entity.CacheEmailValidation{},
		&postgres_entity.CacheEmailValidationDomain{},
		&postgres_entity.CacheIpData{},
		&postgres_entity.CacheIpHunter{},
		&postgres_entity.CacheIPIdentify{},
		&postgres_entity.CosApiEnrichPersonTempResult{},
		&postgres_entity.CurrencyRate{},
		&postgres_entity.CustomerOsIds{},
		&postgres_entity.DMARCMonitoring{},
		&postgres_entity.EmailLookup{},
		&postgres_entity.EmailMessage{},
		&postgres_entity.EmailTracking{},
		&postgres_entity.EmailValidationRecord{},
		&postgres_entity.EmailValidationRequestBulk{},
		&postgres_entity.EnrichDetailsBetterContact{},
		&postgres_entity.EnrichDetailsBrandfetch{},
		&postgres_entity.EnrichDetailsPreFilterTracking{},
		&postgres_entity.EnrichDetailsScrapIn{},
		&postgres_entity.EnrichDetailsTracking{},
		&postgres_entity.EventBuffer{},
		&postgres_entity.ExternalAppKeys{},
		&postgres_entity.Flows{},
		&postgres_entity.FlowEdge{},
		&postgres_entity.FlowExecution{},
		&postgres_entity.FlowNode{},
		&postgres_entity.FlowTransitionsRegistry{},
		&postgres_entity.GlobalOrganization{},
		&postgres_entity.GlobalOrganizationWebsiteToProcess{},
		&postgres_entity.MailStackDomain{},
		&postgres_entity.MailstackBuyRequest{},
		&postgres_entity.MailstackBuyRequestDomain{},
		&postgres_entity.MailstackReputationEntity{},
		&postgres_entity.MagicLink{},
		&postgres_entity.OrganizationWebsiteHostingPlatform{},
		&postgres_entity.PersonalEmailProvider{},
		&postgres_entity.PersonalIntegration{},
		&postgres_entity.PostmarkApiKey{},
		&postgres_entity.SlackChannel{},
		&postgres_entity.SlackChannelNotification{},
		&postgres_entity.SlackSettingsEntity{},
		&postgres_entity.StatsApiCalls{},
		&postgres_entity.TableViewDefinition{},
		&postgres_entity.Tenant{},
		&postgres_entity.TenantSettings{},
		&postgres_entity.TenantSettingsEmailExclusion{},
		&postgres_entity.TenantSettingsMailbox{},
		&postgres_entity.TenantSettingsOpportunityStage{},
		&postgres_entity.TenantWebhook{},
		&postgres_entity.TenantWebhookApiKey{},
		&postgres_entity.TrackingAllowedOrigin{},
		&postgres_entity.UserWorkingSchedule{},
		&postgres_entity.Webhooks{},
		&postgres_entity.WebSession{},
		&postgres_entity.WebTrackerEvents{},
	)
	if err != nil {
		panic(err)
	}

	err = postgresDB.AsyncGormDB.AutoMigrate(
		&postgres_entity.GoogleServiceAccountKey{},
		&postgres_entity.OAuthTokenEntity{},
		&postgres_entity.RawEmail{},
		&postgres_entity.UserEmailImportState{},
		&postgres_entity.UserEmailImportStateHistory{},
	)
	if err != nil {
		panic(err)
	}
}

func (r *Repositories) InitData(ctx context.Context, postgresRepos *Repositories) {
	err := r.AgentRegistryRepository.Initialize(ctx)
	if err != nil {
		panic(err)
	}

	err = r.FlowTransitionsRegistryRepository.Initialize(ctx)
	if err != nil {
		panic(err)
	}
}
