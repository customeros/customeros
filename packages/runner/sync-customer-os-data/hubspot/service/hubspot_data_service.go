package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	hubspotEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/repository"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

type hubspotDataService struct {
	airbyteStoreDb *config.AirbyteStoreDB
	tenant         string
	contacts       map[string]hubspotEntity.Contact
	companies      map[string]hubspotEntity.Company
	owners         map[string]hubspotEntity.Owner
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.DataService {
	return &hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		contacts:       map[string]hubspotEntity.Contact{},
		companies:      map[string]hubspotEntity.Company{},
		owners:         map[string]hubspotEntity.Owner{},
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
		// set main contact fields
		contactForCustomerOs := entity.ContactData{
			ExternalId:       v.Id,
			ExternalSystem:   s.SourceId(),
			FirstName:        hubspotContactProperties.FirstName,
			LastName:         hubspotContactProperties.LastName,
			JobTitle:         hubspotContactProperties.JobTitle,
			CreatedAt:        v.CreateDate.UTC(),
			PrimaryEmail:     hubspotContactProperties.Email,
			AdditionalEmails: strings.Split(hubspotContactProperties.AdditionalEmails, ";"),
			PrimaryE164:      hubspotContactProperties.PhoneNumber,
			Readonly:         true,
		}
		// set reference to primary company
		if hubspotContactProperties.PrimaryCompanyExternalId.Valid {
			contactForCustomerOs.PrimaryCompanyExternalId = strconv.FormatFloat(hubspotContactProperties.PrimaryCompanyExternalId.Float64, 'f', 0, 64)
		}
		// set reference to all linked companies
		var companiesExternalIds []int64
		v.CompaniesExternalIds.AssignTo(&companiesExternalIds)
		if companiesExternalIds != nil {
			var strCompaniesExternalIds []string
			for _, v := range companiesExternalIds {
				companyExternalId := strconv.FormatInt(v, 10)
				strCompaniesExternalIds = append(strCompaniesExternalIds, companyExternalId)
			}
			contactForCustomerOs.CompaniesExternalIds = strCompaniesExternalIds
		}

		customerOsContacts = append(customerOsContacts, contactForCustomerOs)
		s.contacts[v.Id] = v
	}
	return customerOsContacts
}

func (s *hubspotDataService) GetCompaniesForSync(batchSize int) []entity.CompanyData {
	hubspotCompanies, err := repository.GetCompanies(s.getDb(), batchSize)
	if err != nil {
		log.Print(err)
		return nil
	}
	customerOsCompanies := []entity.CompanyData{}
	for _, v := range hubspotCompanies {
		hubspotCompanyProperties, err := repository.GetCompanyProperties(s.getDb(), v.AirbyteAbId, v.AirbyteCompaniesHashid)
		if err != nil {
			log.Print(err)
			continue
		}
		customerOsCompanies = append(customerOsCompanies, entity.CompanyData{
			ExternalId:     v.Id,
			ExternalSystem: s.SourceId(),
			Name:           hubspotCompanyProperties.Name,
			Description:    hubspotCompanyProperties.Description,
			Domain:         hubspotCompanyProperties.Domain,
			Website:        hubspotCompanyProperties.Website,
			Industry:       hubspotCompanyProperties.Industry,
			IsPublic:       hubspotCompanyProperties.IsPublic,
			CreatedAt:      v.CreateDate.UTC(),
			Readonly:       true,
		})
		s.companies[v.Id] = v
	}
	return customerOsCompanies
}

func (s *hubspotDataService) GetUsersForSync(batchSize int) []entity.UserData {
	hubspotOwners, err := repository.GetOwners(s.getDb(), batchSize)
	if err != nil {
		log.Print(err)
		return nil
	}
	customerOsUsers := []entity.UserData{}
	for _, v := range hubspotOwners {
		customerOsUsers = append(customerOsUsers, entity.UserData{
			ExternalId:     v.Id,
			ExternalSystem: s.SourceId(),
			FirstName:      v.FirstName,
			LastName:       v.LastName,
			Email:          v.Email,
			CreatedAt:      v.CreateDate.UTC(),
			Readonly:       true,
		})
		s.owners[v.Id] = v
	}
	return customerOsUsers
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

func (s *hubspotDataService) MarkCompanyProcessed(externalId string, synced bool) error {
	company, ok := s.companies[externalId]
	if ok {
		err := repository.MarkCompanyProcessed(s.getDb(), company, synced)
		if err != nil {
			log.Printf("error while marking company with external reference %s as synced for hubspot", externalId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkUserProcessed(externalId string, synced bool) error {
	owner, ok := s.owners[externalId]
	if ok {
		err := repository.MarkOwnerProcessed(s.getDb(), owner, synced)
		if err != nil {
			log.Printf("error while marking owner with external reference %s as synced for hubspot", externalId)
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
	err = s.getDb().AutoMigrate(&hubspotEntity.SyncStatusCompany{})
	if err != nil {
		log.Print(err)
	}
	err = s.getDb().AutoMigrate(&hubspotEntity.SyncStatusOwner{})
	if err != nil {
		log.Print(err)
	}
}

func (s *hubspotDataService) getDb() *gorm.DB {
	return s.airbyteStoreDb.GetDBHandler(&config.Context{
		Schema: config.CommonSchemaPrefix + s.tenant,
	})
}

func (s *hubspotDataService) SourceId() string {
	return "hubspot"
}

func (s *hubspotDataService) Close() {
	s.contacts = make(map[string]hubspotEntity.Contact)
	s.companies = make(map[string]hubspotEntity.Company)
	s.owners = make(map[string]hubspotEntity.Owner)
}
