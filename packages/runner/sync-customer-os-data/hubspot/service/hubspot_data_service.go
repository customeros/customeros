package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	hubspotEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/hubspot/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type hubspotDataService struct {
	airbyteStoreDb *config.AirbyteStoreDB
	tenant         string
	contacts       map[string]hubspotEntity.Contact
	companies      map[string]hubspotEntity.Company
	owners         map[string]hubspotEntity.Owner
	notes          map[string]hubspotEntity.Note
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.DataService {
	return &hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		contacts:       map[string]hubspotEntity.Contact{},
		companies:      map[string]hubspotEntity.Company{},
		owners:         map[string]hubspotEntity.Owner{},
		notes:          map[string]hubspotEntity.Note{},
	}
}

func (s *hubspotDataService) GetContactsForSync(batchSize int, runId string) []entity.ContactData {
	hubspotContacts, err := repository.GetContacts(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	customerOsContacts := []entity.ContactData{}
	for _, v := range hubspotContacts {
		hubspotContactProperties, err := repository.GetContactProperties(s.getDb(), v.AirbyteAbId, v.AirbyteContactsHashid)
		if err != nil {
			logrus.Error(err)
			continue
		}
		// set main contact fields
		contactForCustomerOs := entity.ContactData{
			ExternalId:          v.Id,
			ExternalSystem:      s.SourceId(),
			FirstName:           hubspotContactProperties.FirstName,
			LastName:            hubspotContactProperties.LastName,
			JobTitle:            hubspotContactProperties.JobTitle,
			CreatedAt:           v.CreateDate.UTC(),
			PrimaryEmail:        hubspotContactProperties.Email,
			AdditionalEmails:    strings.Split(hubspotContactProperties.AdditionalEmails, ";"),
			PrimaryE164:         hubspotContactProperties.PhoneNumber,
			UserOwnerExternalId: hubspotContactProperties.OwnerId,
			Readonly:            true,
			Country:             hubspotContactProperties.Country,
			State:               hubspotContactProperties.State,
			City:                hubspotContactProperties.City,
			Address:             hubspotContactProperties.Address,
			Zip:                 hubspotContactProperties.Zip,
			Fax:                 hubspotContactProperties.Fax,
		}
		// set reference to primary company
		if hubspotContactProperties.PrimaryCompanyExternalId.Valid {
			contactForCustomerOs.PrimaryOrganizationExternalId = strconv.FormatFloat(hubspotContactProperties.PrimaryCompanyExternalId.Float64, 'f', 0, 64)
		}
		// set reference to all linked companies
		var companiesExternalIds []any
		v.CompaniesExternalIds.AssignTo(&companiesExternalIds)
		var strCompaniesExternalIds []string
		for _, c := range companiesExternalIds {
			if _, ok := c.(string); ok {
				strCompaniesExternalIds = append(strCompaniesExternalIds, c.(string))
			} else if _, ok := c.(int64); ok {
				companyExternalId := strconv.FormatInt(c.(int64), 10)
				strCompaniesExternalIds = append(strCompaniesExternalIds, companyExternalId)
			} else if _, ok := c.(float64); ok {
				companyExternalId := strconv.FormatFloat(c.(float64), 'f', 0, 64)
				strCompaniesExternalIds = append(strCompaniesExternalIds, companyExternalId)
			}
		}
		contactForCustomerOs.OrganizationsExternalIds = strCompaniesExternalIds
		// set custom fields
		var textCustomFields []entity.TextCustomField
		if len(hubspotContactProperties.LifecycleStage) > 0 {
			textCustomFields = append(textCustomFields, entity.TextCustomField{
				Name:           "Hubspot Lifecycle Stage",
				Value:          hubspotContactProperties.LifecycleStage,
				ExternalSystem: s.SourceId(),
			})
		}
		contactForCustomerOs.TextCustomFields = textCustomFields

		// set contact type
		if isCustomerContactType(hubspotContactProperties.LifecycleStage) {
			contactForCustomerOs.ContactTypeName = "CUSTOMER"
		} else {
			contactForCustomerOs.ContactTypeName = "NOT_SET"
		}

		customerOsContacts = append(customerOsContacts, contactForCustomerOs)
		s.contacts[v.Id] = v
	}
	return customerOsContacts
}

func isCustomerContactType(hubspotLifecycleStage string) bool {
	customerLifecycleStages := map[string]bool{
		"lead": true, "subscriber": true, "marketingqualifiedlead": true, "salesqualifiedlead": true,
		"opportunity": true, "evangelist": true, "customer": true}
	return customerLifecycleStages[hubspotLifecycleStage]
}

func (s *hubspotDataService) GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData {
	hubspotCompanies, err := repository.GetCompanies(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	customerOsOrganizations := []entity.OrganizationData{}
	for _, v := range hubspotCompanies {
		hubspotCompanyProperties, err := repository.GetCompanyProperties(s.getDb(), v.AirbyteAbId, v.AirbyteCompaniesHashid)
		if err != nil {
			logrus.Error(err)
			continue
		}
		customerOsOrganizations = append(customerOsOrganizations, entity.OrganizationData{
			ExternalId:           v.Id,
			ExternalSystem:       s.SourceId(),
			Name:                 hubspotCompanyProperties.Name,
			Description:          hubspotCompanyProperties.Description,
			Domain:               hubspotCompanyProperties.Domain,
			Website:              hubspotCompanyProperties.Website,
			Industry:             hubspotCompanyProperties.Industry,
			IsPublic:             hubspotCompanyProperties.IsPublic,
			CreatedAt:            v.CreateDate.UTC(),
			Readonly:             true,
			Country:              hubspotCompanyProperties.Country,
			State:                hubspotCompanyProperties.State,
			City:                 hubspotCompanyProperties.City,
			Address:              hubspotCompanyProperties.Address,
			Address2:             hubspotCompanyProperties.Address2,
			Zip:                  hubspotCompanyProperties.Zip,
			Phone:                hubspotCompanyProperties.Phone,
			OrganizationTypeName: "COMPANY",
		})
		s.companies[v.Id] = v
	}
	return customerOsOrganizations
}

func (s *hubspotDataService) GetUsersForSync(batchSize int, runId string) []entity.UserData {
	hubspotOwners, err := repository.GetOwners(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
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

func (s *hubspotDataService) GetNotesForSync(batchSize int, runId string) []entity.NoteData {
	hubspotNotes, err := repository.GetNotes(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	customerOsNotes := []entity.NoteData{}
	for _, v := range hubspotNotes {
		hubspotNoteProperties, err := repository.GetNoteProperties(s.getDb(), v.AirbyteAbId, v.AirbyteNotesHashid)
		if err != nil {
			logrus.Error(err)
			continue
		}
		// set main fields
		noteForCustomerOs := entity.NoteData{
			ExternalId:     v.Id,
			ExternalSystem: s.SourceId(),
			CreatedAt:      v.CreateDate.UTC(),
			Html:           hubspotNoteProperties.NoteBody,
			UserExternalId: hubspotNoteProperties.OwnerId,
		}
		// set reference to all linked contacts
		var contactsExternalIds []any
		v.ContactsExternalIds.AssignTo(&contactsExternalIds)
		var strContactsExternalIds []string
		for _, c := range contactsExternalIds {
			if _, ok := c.(string); ok {
				strContactsExternalIds = append(strContactsExternalIds, c.(string))
			} else if _, ok := c.(int64); ok {
				contactExternalId := strconv.FormatInt(c.(int64), 10)
				strContactsExternalIds = append(strContactsExternalIds, contactExternalId)
			} else if _, ok := c.(float64); ok {
				contactExternalId := strconv.FormatFloat(c.(float64), 'f', 0, 64)
				strContactsExternalIds = append(strContactsExternalIds, contactExternalId)
			}
		}
		noteForCustomerOs.ContactsExternalIds = strContactsExternalIds
		customerOsNotes = append(customerOsNotes, noteForCustomerOs)
		s.notes[v.Id] = v
	}
	return customerOsNotes
}

func (s *hubspotDataService) MarkContactProcessed(externalId, runId string, synced bool) error {
	contact, ok := s.contacts[externalId]
	if ok {
		err := repository.MarkContactProcessed(s.getDb(), contact, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking contact with external reference %s as synced for hubspot", externalId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkOrganizationProcessed(externalId, runId string, synced bool) error {
	company, ok := s.companies[externalId]
	if ok {
		err := repository.MarkCompanyProcessed(s.getDb(), company, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking company with external reference %s as synced for hubspot", externalId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkUserProcessed(externalId, runId string, synced bool) error {
	owner, ok := s.owners[externalId]
	if ok {
		err := repository.MarkOwnerProcessed(s.getDb(), owner, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking owner with external reference %s as synced for hubspot", externalId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkNoteProcessed(externalId, runId string, synced bool) error {
	note, ok := s.notes[externalId]
	if ok {
		err := repository.MarkNoteProcessed(s.getDb(), note, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking note with external reference %s as synced for hubspot", externalId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) Refresh() {
	err := s.getDb().AutoMigrate(&hubspotEntity.SyncStatusContact{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&hubspotEntity.SyncStatusCompany{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&hubspotEntity.SyncStatusOwner{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&hubspotEntity.SyncStatusNote{})
	if err != nil {
		logrus.Error(err)
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
	s.notes = make(map[string]hubspotEntity.Note)
}
