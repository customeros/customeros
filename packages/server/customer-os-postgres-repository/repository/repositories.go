package repository

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"gorm.io/gorm"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type Repositories struct {
	Db      *gorm.DB
	AsyncDb *gorm.DB

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
	FlowRepository                               FlowRepository
	FlowAgentExecutionRepository                 FlowAgentExecutionRepository
	FlowAgentRegistryRepository                  FlowAgentRegistryRepository
	FlowDeadEventsRepository                     FlowDeadEventsRepository
	FlowEdgeRepository                           FlowEdgeRepository
	FlowExecutionRepository                      FlowExecutionRepository
	FlowListenerRegistryRepository               FlowListenerRegistryRepository
	FlowNodeRepository                           FlowNodeRepository
	FlowTransitionsRegistryRepository            FlowTransitionsRegistryRepository
	FlowWebhooksRepository                       FlowWebhooksRepository
	GoogleServiceAccountKeyRepository            GoogleServiceAccountKeyRepository
	IndustryMappingRepository                    IndustryMappingRepository
	MailStackDomainRepository                    MailStackDomainRepository
	MailstackBuyRequestRepository                MailstackBuyRequestRepository
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
	TrackingRepository                           TrackingRepository
	UserEmailImportPageTokenRepository           UserEmailImportStateRepository
	UserWorkingScheduleRepository                UserWorkingScheduleRepository
	GlobalOrganizationRepository                 GlobalOrganizationRepository
	GlobalOrganizationWebsiteToProcessRepository GlobalOrganizationWebsiteToProcessRepository
}

func InitRepositories(postgresDB *config.PostgresDB) *Repositories {
	repositories := &Repositories{
		Db:      postgresDB.GormDB,
		AsyncDb: postgresDB.AsyncGormDB,

		OAuthTokenRepository:               NewOAuthTokenRepository(postgresDB.AsyncGormDB),
		RawEmailRepository:                 NewRawEmailRepository(postgresDB.AsyncGormDB),
		GoogleServiceAccountKeyRepository:  NewGoogleServiceAccountKeyRepository(postgresDB.AsyncGormDB),
		UserEmailImportPageTokenRepository: NewUserEmailImportStateRepository(postgresDB.AsyncGormDB),

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
		FlowRepository:                               NewFlowRepository(postgresDB.GormDB),
		FlowAgentExecutionRepository:                 NewFlowAgentExecutionRepository(postgresDB.GormDB),
		FlowAgentRegistryRepository:                  NewFlowAgentRegistryRepository(postgresDB.GormDB),
		FlowDeadEventsRepository:                     NewFlowDeadEventsRepository(postgresDB.GormDB),
		FlowEdgeRepository:                           NewFlowEdgeRepository(postgresDB.GormDB),
		FlowExecutionRepository:                      NewFlowExecutionRepository(postgresDB.GormDB),
		FlowListenerRegistryRepository:               NewFlowListenerRegistryRepository(postgresDB.GormDB),
		FlowNodeRepository:                           NewFlowNodeRepository(postgresDB.GormDB),
		FlowTransitionsRegistryRepository:            NewFlowTransitionsRegistryRepository(postgresDB.GormDB),
		FlowWebhooksRepository:                       NewFlowWebhooksRepository(postgresDB.GormDB),
		IndustryMappingRepository:                    NewIndustryMappingRepository(postgresDB.GormDB),
		MailStackDomainRepository:                    NewMailStackDomainRepository(postgresDB.GormDB),
		MailstackBuyRequestRepository:                NewMailstackBuyRequestRepository(postgresDB.GormDB),
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
		TrackingRepository:                           NewTrackingRepository(postgresDB.GormDB),
		UserWorkingScheduleRepository:                NewUserWorkingScheduleRepository(postgresDB.GormDB),
		GlobalOrganizationRepository:                 NewGlobalOrganizationRepository(postgresDB.GormDB),
		GlobalOrganizationWebsiteToProcessRepository: NewGlobalOrganizationWebsiteToProcessRepository(postgresDB.GormDB),
	}

	return repositories
}

func (r *Repositories) Migration(postgresDB *config.PostgresDB) {
	err := postgresDB.GormDB.AutoMigrate(
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
		&entity.Flow{},
		&entity.FlowAgentRegistry{},
		&entity.FlowAgentExecution{},
		&entity.FlowDeadEvents{},
		&entity.FlowEdge{},
		&entity.FlowExecution{},
		&entity.FlowListenerRegistry{},
		&entity.FlowNode{},
		&entity.FlowTransitionsRegistry{},
		&entity.FlowWebhooks{},
		&entity.IndustryMapping{},
		&entity.MailStackDomain{},
		&entity.MailstackBuyRequest{},
		&entity.MailstackBuyRequestDomain{},
		&entity.MailstackReputationEntity{},
		&entity.OrganizationWebsiteHostingPlatform{},
		&entity.PersonalEmailProvider{},
		&entity.PersonalIntegration{},
		&entity.PostmarkApiKey{},
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
		&entity.UserWorkingSchedule{},
		&entity.GlobalOrganization{},
		&entity.GlobalOrganizationWebsiteToProcess{},
	)
	if err != nil {
		panic(err)
	}

	err = postgresDB.AsyncGormDB.AutoMigrate(
		&entity.GoogleServiceAccountKey{},
		&entity.OAuthTokenEntity{},
		&entity.RawEmail{},
		&entity.UserEmailImportState{},
		&entity.UserEmailImportStateHistory{},
	)
	if err != nil {
		panic(err)
	}
}

func (r *Repositories) InitData(ctx context.Context, postgresRepos *Repositories) {
	err := r.FlowAgentRegistryRepository.Initialize(ctx)
	if err != nil {
		panic(err)
	}

	err = r.FlowListenerRegistryRepository.Initialize(ctx)
	if err != nil {
		panic(err)
	}

	err = r.FlowTransitionsRegistryRepository.Initialize(ctx)
	if err != nil {
		panic(err)
	}
}
