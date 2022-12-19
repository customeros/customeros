package service

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/service"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"log"
	"time"
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

// TODO automatically create external system in neo4j and link with tenant
func (s *syncService) Sync() {
	tenantsToSync, err := s.repositories.TenantSyncSettingsRepository.GetTenantsForSync()
	if err != nil {
		log.Print("failed to get tenants for sync")
		return
	}

	for _, v := range tenantsToSync {
		dataService, err := s.dataService(v)
		if err != nil {
			log.Printf("failed to get data service for tenant %v: %v", v.Tenant, err)
			continue
		}

		defer dataService.Close()

		syncDate := time.Now().UTC()

		_ = s.repositories.ExternalSystemRepository.Merge(v.Tenant, dataService.SourceId())

		s.syncCompanies(dataService, syncDate, v.Tenant)
		s.syncContacts(dataService, syncDate, v.Tenant)
	}
}

func (s *syncService) syncContacts(dataService common.DataService, syncDate time.Time, tenant string) {
	for {
		contacts := dataService.GetContactsForSync(batchSize)
		if len(contacts) == 0 {
			log.Printf("no contacts found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		log.Printf("syncing %d contacts from %s for tenant %s", len(contacts), dataService.SourceId(), tenant)

		for _, v := range contacts {
			var failedSync = false

			contactId, err := s.repositories.ContactRepository.MergeContact(tenant, syncDate, v)
			if err != nil {
				failedSync = true
				log.Printf("failed merge contact with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			if len(v.PrimaryEmail) > 0 {
				if err = s.repositories.ContactRepository.MergePrimaryEmail(tenant, contactId, v.PrimaryEmail); err != nil {
					failedSync = true
					log.Printf("failed merge primary email for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			for _, additionalEmail := range v.AdditionalEmails {
				if err = s.repositories.ContactRepository.MergeAdditionalEmail(tenant, contactId, additionalEmail); err != nil {
					failedSync = true
					log.Printf("failed merge additional email for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if len(v.PrimaryE164) > 0 {
				if err = s.repositories.ContactRepository.MergePrimaryPhoneNumber(tenant, contactId, v.PrimaryE164); err != nil {
					failedSync = true
					log.Printf("failed merge primary phone number for contact with external reference %v , tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			log.Printf("successfully merged contact with id %v for tenant %v from %v", contactId, tenant, dataService.SourceId())
			if err := dataService.MarkContactProcessed(v.ExternalId, failedSync == false); err != nil {
				continue
			}
		}
		if len(contacts) < batchSize {
			break
		}
	}
}

func (s *syncService) syncCompanies(dataService common.DataService, syncDate time.Time, tenant string) {
	for {
		companies := dataService.GetCompaniesForSync(batchSize)
		if len(companies) == 0 {
			log.Printf("no companies found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		log.Printf("syncing %d companies from %s for tenant %s", len(companies), dataService.SourceId(), tenant)

		for _, v := range companies {
			var failedSync = false

			companyId, err := s.repositories.CompanyRepository.MergeCompany(tenant, syncDate, v)
			if err != nil {
				failedSync = true
				log.Printf("failed merge company with external reference %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			log.Printf("successfully merged company with id %v for tenant %v from %v", companyId, tenant, dataService.SourceId())
			if err := dataService.MarkCompanyProcessed(v.ExternalId, failedSync == false); err != nil {
				continue
			}
		}
		if len(companies) < batchSize {
			break
		}
	}
}

func (s *syncService) dataService(tenantToSync entity.TenantSyncSettings) (common.DataService, error) {
	// Use a map to store the different implementations of common.DataService as functions.
	dataServiceMap := map[entity.AirbyteSource]func() common.DataService{
		entity.HUBSPOT: func() common.DataService {
			return service.NewHubspotDataService(s.repositories.Dbs.AirbyteStoreDB, tenantToSync.Tenant)
		},
		// Add additional implementations here.
	}

	// Look up the corresponding implementation in the map using the tenantToSync.Source value.
	createDataService, ok := dataServiceMap[tenantToSync.Source]
	if !ok {
		// Return an error if the tenantToSync.Source value is not recognized.
		return nil, fmt.Errorf("unknown airbyte source %v, skipping sync for tenant %v", tenantToSync.Source, tenantToSync.Tenant)
	}

	// Call the createDataService function to create a new instance of common.DataService.
	dataService := createDataService()

	// Call the Refresh method on the dataService instance.
	dataService.Refresh()

	return dataService, nil
}
