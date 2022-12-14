package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	hubspotEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/repository"
	"gorm.io/gorm"
	"log"
)

type hubspotDataService struct {
	airbyteStoreDb *config.AirbyteStoreDB
	tenant         string
}

func (s *hubspotDataService) GetContactsForSync(batchSize int) []entity.ContactEntity {
	hubspotContacts, err := repository.GetContacts(s.getDb(), batchSize)
	if err != nil {
		log.Print(err)
		return nil
	}
	customerOsContacts := []entity.ContactEntity{}
	for _, v := range hubspotContacts {
		hubspotContactProperties, err := repository.GetContactProperties(s.getDb(), v.AirbyteAbId, v.AirbyteContactsHashid)
		if err != nil {
			log.Print(err)
			continue
		}
		customerOsContacts = append(customerOsContacts, entity.ContactEntity{
			ExternalReference: v.Id,
			ExternalSystem:    "hubspot",
			FirstName:         hubspotContactProperties.FirstName,
			LastName:          hubspotContactProperties.LastName,
		})
	}
	return customerOsContacts
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.DataService {
	return &hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
	}
}

func (s *hubspotDataService) Refresh() {
	// TODO automigrate only if table exists
	err := s.getDb().AutoMigrate(&hubspotEntity.Contact{})
	if err != nil {
		log.Print(err)
	}
}

func (s *hubspotDataService) getDb() *gorm.DB {
	return s.airbyteStoreDb.GetDBHandler(&config.Context{
		Schema: config.CommonSchemaPrefix + s.tenant,
	})
}
