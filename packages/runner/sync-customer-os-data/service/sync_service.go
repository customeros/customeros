package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/service"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"log"
)

const batchSize = 100

type SyncService interface {
	Sync()
}

type syncService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewSyncService(repositories *repository.Repositories, services *Services) SyncService {
	return &syncService{
		repositories: repositories,
		services:     services,
	}
}

func (s *syncService) Sync() {
	tenantsToSync, err := s.repositories.TenantSyncSettingsRepository.GetTenantsForSync()
	if err != nil {
		log.Print("failed to get tenants for sync")
		return
	}
	for _, v := range tenantsToSync {
		dataService, err := s.dataService(v)
		if err != nil {
			continue
		}
		syncContacts(dataService, v.Tenant)
	}
}

func syncContacts(dataService common.DataService, tenant string) {
	for {
		contacts := dataService.GetContactsForSync(batchSize)
		if len(contacts) == 0 {
			log.Printf("no contacts found for sync from %s for tenant %s", dataService.SourceName(), tenant)
			break
		}
		log.Printf("syncing %d contacts from %s for tenant %s", len(contacts), dataService.SourceName(), tenant)
		for _, v := range contacts {
			dataService.MarkContactSynced(v.ExternalId)
		}
		if len(contacts) < batchSize {
			break
		}
	}
}

func (s *syncService) dataService(tenantToSync entity.TenantSyncSettings) (common.DataService, error) {
	switch tenantToSync.Source {
	case entity.HUBSPOT:
		dataService := service.NewHubspotDataService(s.repositories.Dbs.AirbyteStoreDB, tenantToSync.Tenant)
		dataService.Refresh()
		return dataService, nil
	}
	return nil, fmt.Errorf("unkown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
}
