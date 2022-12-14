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
		syncContacts(dataService)
	}

	log.Printf("found %d tenants to sync", len(tenantsToSync))
}

func syncContacts(dataService common.DataService) {
	for {
		contacts := dataService.GetContactsForSync(batchSize)
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
