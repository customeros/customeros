package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	AiLocationMappingRepository              AiLocationMappingRepository
	AiPromptLogRepository                    AiPromptLogRepository
	AppKeyRepository                         AppKeyRepository
	FlowRepository                           FlowRepository
	FlowSequenceRepository                   FlowSequenceRepository
	FlowSequenceStepRepository               FlowSequenceStepRepository
	FlowSequenceContactRepository            FlowSequenceContactRepository
	FlowSequenceSenderRepository             FlowSequenceSenderRepository
	PersonalIntegrationRepository            PersonalIntegrationRepository
	PersonalEmailProviderRepository          PersonalEmailProviderRepository
	TenantWebhookApiKeyRepository            TenantWebhookApiKeyRepository
	TenantWebhookRepository                  TenantWebhookRepository
	SlackChannelRepository                   SlackChannelRepository
	PostmarkApiKeyRepository                 PostmarkApiKeyRepository
	GoogleServiceAccountKeyRepository        GoogleServiceAccountKeyRepository
	CurrencyRateRepository                   CurrencyRateRepository
	EventBufferRepository                    EventBufferRepository
	TableViewDefinitionRepository            TableViewDefinitionRepository
	TrackingAllowedOriginRepository          TrackingAllowedOriginRepository
	TechLimitRepository                      TechLimitRepository
	ExternalAppKeysRepository                ExternalAppKeysRepository
	EnrichDetailsBetterContactRepository     EnrichDetailsBetterContactRepository
	EnrichDetailsScrapInRepository           EnrichDetailsScrapInRepository
	EnrichDetailsPrefilterTrackingRepository EnrichDetailsPrefilterTrackingRepository
	EnrichDetailsTrackingRepository          EnrichDetailsTrackingRepository
	UserEmailImportPageTokenRepository       UserEmailImportStateRepository
	RawEmailRepository                       RawEmailRepository
	OAuthTokenRepository                     OAuthTokenRepository
	SlackSettingsRepository                  SlackSettingsRepository
	SlackChannelNotificationRepository       SlackChannelNotificationRepository
	ApiCacheRepository                       ApiCacheRepository
	WorkflowRepository                       WorkflowRepository
	IndustryMappingRepository                IndustryMappingRepository
	TrackingRepository                       TrackingRepository
	TenantSettingsRepository                 TenantSettingsRepository
	TenantSettingsOpportunityStageRepository TenantSettingsOpportunityStageRepository
	TenantSettingsMailboxRepository          TenantSettingsMailboxRepository
	TenantSettingsEmailExclusionRepository   TenantSettingsEmailExclusionRepository
	EmailLookupRepository                    EmailLookupRepository
	EmailTrackingRepository                  EmailTrackingRepository
	TenantRepository                         TenantRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		AiLocationMappingRepository:              NewAiLocationMappingRepository(db),
		AiPromptLogRepository:                    NewAiPromptLogRepository(db),
		AppKeyRepository:                         NewAppKeyRepo(db),
		FlowRepository:                           NewFlowRepository(db),
		FlowSequenceRepository:                   NewFlowSequenceRepository(db),
		FlowSequenceStepRepository:               NewFlowSequenceStepRepository(db),
		FlowSequenceContactRepository:            NewFlowSequenceContactRepository(db),
		FlowSequenceSenderRepository:             NewFlowSequenceSenderRepository(db),
		PersonalIntegrationRepository:            NewPersonalIntegrationsRepo(db),
		PersonalEmailProviderRepository:          NewPersonalEmailProviderRepository(db),
		TenantWebhookApiKeyRepository:            NewTenantWebhookApiKeyRepository(db),
		TenantWebhookRepository:                  NewTenantWebhookRepo(db),
		SlackChannelRepository:                   NewSlackChannelRepository(db),
		PostmarkApiKeyRepository:                 NewPostmarkApiKeyRepo(db),
		GoogleServiceAccountKeyRepository:        NewGoogleServiceAccountKeyRepository(db),
		CurrencyRateRepository:                   NewCurrencyRateRepository(db),
		EventBufferRepository:                    NewEventBufferRepository(db),
		TableViewDefinitionRepository:            NewTableViewDefinitionRepository(db),
		TrackingAllowedOriginRepository:          NewTrackingAllowedOriginRepository(db),
		TechLimitRepository:                      NewTechLimitRepository(db),
		ExternalAppKeysRepository:                NewExternalAppKeysRepository(db),
		EnrichDetailsBetterContactRepository:     NewEnrichDetailsBetterContactRepository(db),
		EnrichDetailsScrapInRepository:           NewEnrichDetailsScrapInRepository(db),
		EnrichDetailsPrefilterTrackingRepository: NewEnrichDetailsPrefilterTrackingRepository(db),
		EnrichDetailsTrackingRepository:          NewEnrichDetailsTrackingRepository(db),
		UserEmailImportPageTokenRepository:       NewUserEmailImportStateRepository(db),
		RawEmailRepository:                       NewRawEmailRepository(db),
		OAuthTokenRepository:                     NewOAuthTokenRepository(db),
		SlackSettingsRepository:                  NewSlackSettingsRepository(db),
		SlackChannelNotificationRepository:       NewSlackChannelNotificationRepository(db),
		ApiCacheRepository:                       NewApiCacheRepository(db),
		WorkflowRepository:                       NewWorkflowRepository(db),
		IndustryMappingRepository:                NewIndustryMappingRepository(db),
		TrackingRepository:                       NewTrackingRepository(db),
		TenantSettingsRepository:                 NewTenantSettingsRepository(db),
		TenantSettingsOpportunityStageRepository: NewTenantSettingsOpportunityStageRepository(db),
		TenantSettingsMailboxRepository:          NewTenantSettingsMailboxRepository(db),
		TenantSettingsEmailExclusionRepository:   NewEmailExclusionRepository(db),
		EmailLookupRepository:                    NewEmailLookupRepository(db),
		EmailTrackingRepository:                  NewEmailTrackingRepository(db),
		TenantRepository:                         NewTenantRepository(db),
	}

	return repositories
}

func (r *Repositories) Migration(db *gorm.DB) {

	var err error

	err = db.AutoMigrate(&entity.Tenant{})
	if err != nil {
		panic(err)
	}

	//err = db.AutoMigrate(&entity.AppKey{})
	//if err != nil {
	//	panic(err)
	//}

	err = db.AutoMigrate(&entity.AiLocationMapping{})
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

	err = db.AutoMigrate(&entity.TenantSettingsEmailExclusion{})
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

	err = db.AutoMigrate(&entity.EnrichDetailsScrapIn{})
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

	err = db.AutoMigrate(&entity.SlackChannelNotification{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.ApiCache{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.Workflow{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.IndustryMapping{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.Tracking{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.EnrichDetailsPreFilterTracking{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.EnrichDetailsTracking{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.TenantSettingsOpportunityStage{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.TenantSettingsMailbox{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.EmailLookup{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.EmailTracking{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.FlowSequenceStepTemplateVariable{}, &entity.Flow{}, &entity.FlowSequence{}, &entity.FlowSequenceStep{}, &entity.FlowSequenceContact{}, &entity.FlowSequenceSender{})
	if err != nil {
		panic(err)
	}
}
