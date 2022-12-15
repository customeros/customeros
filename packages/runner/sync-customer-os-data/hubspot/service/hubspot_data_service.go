package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	hubspotEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/repository"
	"gorm.io/gorm"
	"log"
	"strings"
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

func (s *hubspotDataService) GetContactsForSync(batchSize int) []entity.ContactData {
	hubspotContacts, err := repository.GetContacts(s.getDb(), batchSize)
	if err != nil {
		log.Print(err)
		return nil
	}
	customerOsContacts := []entity.ContactData{}
	for _, v := range hubspotContacts {
		hubspotContactProperties, err := repository.GetContactProperties(s.getDb(), v.AirbyteAbId, v.AirbyteContactsHashid)
		if err != nil {
			log.Print(err)
			continue
		}
		customerOsContacts = append(customerOsContacts, entity.ContactData{
			ExternalId:       v.Id,
			ExternalSystem:   "hubspot",
			FirstName:        hubspotContactProperties.FirstName,
			LastName:         hubspotContactProperties.LastName,
			CreatedAt:        v.CreateDate.UTC(),
			PrimaryEmail:     hubspotContactProperties.Email,
			AdditionalEmails: strings.Split(hubspotContactProperties.AdditionalEmails, ";"),
			PrimaryE164:      hubspotContactProperties.PhoneNumber,
			Readonly:         true,
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
	err := s.getDb().AutoMigrate(&hubspotEntity.SyncStatusContact{})
	if err != nil {
		log.Print(err)
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
