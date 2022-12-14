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

func (s *hubspotDataService) GetContactsForSync() []entity.ContactEntity {
	s.refresh()
	contacts, err := repository.GetContacts(s.getDb())
	if err != nil {
		return nil
	}
	log.Printf("%v", len(contacts))
	return []entity.ContactEntity{}
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.DataService {
	return &hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
	}
}

func (s *hubspotDataService) refresh() {
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
