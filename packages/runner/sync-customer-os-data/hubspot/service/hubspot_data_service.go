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
	contacts       map[string]hubspotEntity.Contact
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.DataService {
	return &hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		contacts:       map[string]hubspotEntity.Contact{},
	}
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
			ExternalId:     v.Id,
			ExternalSystem: "hubspot",
			FirstName:      hubspotContactProperties.FirstName,
			LastName:       hubspotContactProperties.LastName,
			CreatedAt:      v.CreateDate.UTC(),
			Readonly:       true,
		})
		s.contacts[v.Id] = v
	}
	return customerOsContacts
}

func (s *hubspotDataService) MarkContactProcessed(externalId string, synced bool) error {
	contact, ok := s.contacts[externalId]
	if ok {
		err := repository.MarkContactProcessed(s.getDb(), contact, synced)
		if err != nil {
			log.Printf("error while marking contact with external reference %s as synced for hubspot", externalId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) Refresh() {
	if s.getDb().Migrator().HasTable(hubspotEntity.Contact{}.TableName()) {
		err := s.getDb().AutoMigrate(&hubspotEntity.Contact{})
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *hubspotDataService) getDb() *gorm.DB {
	return s.airbyteStoreDb.GetDBHandler(&config.Context{
		Schema: config.CommonSchemaPrefix + s.tenant,
	})
}

func (s *hubspotDataService) SourceName() string {
	return "hubspot"
}

func (s *hubspotDataService) Close() {
	s.contacts = make(map[string]hubspotEntity.Contact)
}
