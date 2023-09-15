package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/hubspot"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/intercom"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/pipedrive"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/salesforce"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/slack"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/zendesk_support"
	localutils "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type SyncCustomerOsDataService interface {
	Sync(ctx context.Context, runId string)
}

type syncService struct {
	repositories   *repository.Repositories
	services       *Services
	cfg            *config.Config
	syncServiceMap map[string]map[common.SyncedEntityType]SyncService
	log            logger.Logger
}

func NewSyncCustomerOsDataService(repositories *repository.Repositories, services *Services, cfg *config.Config, log logger.Logger) SyncCustomerOsDataService {
	service := syncService{
		repositories:   repositories,
		services:       services,
		cfg:            cfg,
		log:            log,
		syncServiceMap: make(map[string]map[common.SyncedEntityType]SyncService),
	}
	// sample to populate map
	service.syncServiceMap[string(entity.AirbyteSourceHubspot)] = map[common.SyncedEntityType]SyncService{
		common.USERS: services.UserDefaultSyncService,
	}
	return &service
}

func (s *syncService) Sync(parentCtx context.Context, runId string) {
	tenantsToSync, err := s.repositories.TenantSyncSettingsRepository.GetTenantsForSync()
	if err != nil {
		s.log.Error("failed to get tenants for sync")
		return
	}

	for _, v := range tenantsToSync {
		syncDate := utils.Now()
		syncRunDtls := entity.SyncRun{
			StartAt:              syncDate,
			RunId:                runId,
			TenantSyncSettingsId: v.ID,
		}
		ctx := localutils.WithCustomContext(parentCtx, &localutils.CustomContext{Tenant: v.Tenant, Source: v.Source, RunId: runId})

		dataService, err := s.sourceDataService(v)
		if err != nil {
			s.log.Errorf("failed to get data service for tenant %v: %v", v.Tenant, err)
			continue
		}

		defer func() {
			dataService.Close()
		}()

		s.syncExternalSystem(ctx, dataService, v.Tenant)

		completed, failed, skipped := s.userSyncService(v).Sync(ctx, dataService, syncDate, v.Tenant, runId, s.cfg.SyncCustomerOsData.BatchSize)
		syncRunDtls.CompletedUsers = completed
		syncRunDtls.FailedUsers = failed
		syncRunDtls.SkippedUsers = skipped

		completed, failed, skipped = s.organizationSyncService(v).Sync(ctx, dataService, syncDate, v.Tenant, runId, s.cfg.SyncCustomerOsData.BatchSize)
		syncRunDtls.CompletedOrganizations = completed
		syncRunDtls.FailedOrganizations = failed
		syncRunDtls.SkippedOrganizations = skipped

		completed, failed, skipped = s.contactSyncService(v).Sync(ctx, dataService, syncDate, v.Tenant, runId, s.cfg.SyncCustomerOsData.BatchSize)
		syncRunDtls.CompletedContacts = completed
		syncRunDtls.FailedContacts = failed
		syncRunDtls.SkippedContacts = skipped

		completed, failed, skipped = s.issueSyncService(v).Sync(ctx, dataService, syncDate, v.Tenant, runId, s.cfg.SyncCustomerOsData.BatchSize)
		syncRunDtls.CompletedIssues = completed
		syncRunDtls.FailedIssues = failed
		syncRunDtls.SkippedIssues = skipped

		completed, failed, skipped = s.noteSyncService(v).Sync(ctx, dataService, syncDate, v.Tenant, runId, s.cfg.SyncCustomerOsData.BatchSize)
		syncRunDtls.CompletedNotes = completed
		syncRunDtls.FailedNotes = failed
		syncRunDtls.SkippedNotes = skipped

		completed, failed, skipped = s.emailMessageSyncService(v).Sync(ctx, dataService, syncDate, v.Tenant, runId, s.cfg.SyncCustomerOsData.BatchSize)
		syncRunDtls.CompletedEmailMessages = completed
		syncRunDtls.FailedEmailMessages = failed
		syncRunDtls.SkippedEmailMessages = skipped

		completed, failed, skipped = s.meetingSyncService(v).Sync(ctx, dataService, syncDate, v.Tenant, runId, s.cfg.SyncCustomerOsData.BatchSize)
		syncRunDtls.CompletedMeetings = completed
		syncRunDtls.FailedMeetings = failed
		syncRunDtls.SkippedMeetings = skipped

		completed, failed, skipped = s.interactionEventSyncService(v).Sync(ctx, dataService, syncDate, v.Tenant, runId, s.cfg.SyncCustomerOsData.BatchSize)
		syncRunDtls.CompletedInteractionEvents = completed
		syncRunDtls.FailedInteractionEvents = failed
		syncRunDtls.SkippedInteractionEvents = skipped

		syncRunDtls.SumTotalFailed()
		syncRunDtls.SumTotalSkipped()
		syncRunDtls.SumTotalCompleted()
		syncRunDtls.EndAt = utils.Now()

		s.repositories.SyncRunRepository.Save(syncRunDtls)
	}
}

func (s *syncService) syncExternalSystem(ctx context.Context, dataService source.SourceDataService, tenant string) {
	_ = s.repositories.ExternalSystemRepository.Merge(ctx, tenant, dataService.SourceId())
}

func (s *syncService) sourceDataService(tenantToSync entity.TenantSyncSettings) (source.SourceDataService, error) {
	// Use a map to store the different implementations of source.SourceDataService as functions.
	dataServiceMap := map[string]func() source.SourceDataService{
		string(entity.AirbyteSourceHubspot): func() source.SourceDataService {
			return hubspot.NewHubspotDataService(s.repositories.Dbs.RawDataStoreDB, tenantToSync.Tenant, s.log)
		},
		string(entity.AirbyteSourceZendeskSupport): func() source.SourceDataService {
			return zendesk_support.NewZendeskSupportDataService(s.repositories.Dbs.RawDataStoreDB, tenantToSync.Tenant, s.log)
		},
		string(entity.AirbyteSourcePipedrive): func() source.SourceDataService {
			return pipedrive.NewPipedriveDataService(s.repositories.Dbs.RawDataStoreDB, tenantToSync.Tenant, s.log)
		},
		string(entity.AirbyteSourceIntercom): func() source.SourceDataService {
			return intercom.NewIntercomDataService(s.repositories.Dbs.RawDataStoreDB, tenantToSync.Tenant, s.log)
		},
		string(entity.AirbyteSourceSalesforce): func() source.SourceDataService {
			return salesforce.NewSalesforceDataService(s.repositories.Dbs.RawDataStoreDB, tenantToSync.Tenant, s.log)
		},
		string(entity.OpenlineSourceSlack): func() source.SourceDataService {
			return slack.NewSlackDataService(s.repositories.Dbs.RawDataStoreDB, tenantToSync.Tenant, s.log)
		},
		// Add additional implementations here.
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createDataService, ok := dataServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	// Call the createDataService function to create a new instance of source.SourceDataService.
	dataService := createDataService()

	// Call the Init method on the sourceDataService instance.
	dataService.Init()

	return dataService, nil
}

func (s *syncService) userSyncService(tenantSyncSettings entity.TenantSyncSettings) SyncService {
	if v, ok := s.syncServiceMap[tenantSyncSettings.Source]; ok {
		if u, ok := v[common.USERS]; ok {
			return u
		}
	}
	return s.services.UserDefaultSyncService
}

func (s *syncService) organizationSyncService(tenantSyncSettings entity.TenantSyncSettings) SyncService {
	if v, ok := s.syncServiceMap[tenantSyncSettings.Source]; ok {
		if u, ok := v[common.ORGANIZATIONS]; ok {
			return u
		}
	}
	return s.services.OrganizationDefaultSyncService
}

func (s *syncService) contactSyncService(tenantSyncSettings entity.TenantSyncSettings) SyncService {
	if v, ok := s.syncServiceMap[tenantSyncSettings.Source]; ok {
		if u, ok := v[common.CONTACTS]; ok {
			return u
		}
	}
	return s.services.ContactDefaultSyncService
}

func (s *syncService) issueSyncService(tenantSyncSettings entity.TenantSyncSettings) SyncService {
	if v, ok := s.syncServiceMap[tenantSyncSettings.Source]; ok {
		if u, ok := v[common.ISSUES]; ok {
			return u
		}
	}
	return s.services.IssueDefaultSyncService
}

func (s *syncService) noteSyncService(tenantSyncSettings entity.TenantSyncSettings) SyncService {
	if v, ok := s.syncServiceMap[tenantSyncSettings.Source]; ok {
		if u, ok := v[common.NOTES]; ok {
			return u
		}
	}
	return s.services.NoteDefaultSyncService
}

func (s *syncService) meetingSyncService(tenantSyncSettings entity.TenantSyncSettings) SyncService {
	if v, ok := s.syncServiceMap[tenantSyncSettings.Source]; ok {
		if u, ok := v[common.MEETINGS]; ok {
			return u
		}
	}
	return s.services.MeetingDefaultSyncService
}

func (s *syncService) emailMessageSyncService(tenantSyncSettings entity.TenantSyncSettings) SyncService {
	if v, ok := s.syncServiceMap[tenantSyncSettings.Source]; ok {
		if u, ok := v[common.EMAIL_MESSAGES]; ok {
			return u
		}
	}
	return s.services.EmailMessageDefaultSyncService
}

func (s *syncService) interactionEventSyncService(tenantSyncSettings entity.TenantSyncSettings) SyncService {
	if v, ok := s.syncServiceMap[tenantSyncSettings.Source]; ok {
		if u, ok := v[common.INTERACTION_EVENTS]; ok {
			return u
		}
	}
	return s.services.InteractionEventDefaultSyncService
}
