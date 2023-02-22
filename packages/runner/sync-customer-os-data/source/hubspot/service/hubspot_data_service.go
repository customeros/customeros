package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	localEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/hubspot/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/hubspot/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type hubspotDataService struct {
	airbyteStoreDb *config.AirbyteStoreDB
	tenant         string
	instance       string
	contacts       map[string]localEntity.Contact
	companies      map[string]localEntity.Company
	owners         map[string]localEntity.Owner
	notes          map[string]localEntity.Note
	emails         map[string]localEntity.Email
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.SourceDataService {
	return &hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		contacts:       map[string]localEntity.Contact{},
		companies:      map[string]localEntity.Company{},
		owners:         map[string]localEntity.Owner{},
		notes:          map[string]localEntity.Note{},
		emails:         map[string]localEntity.Email{},
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
			UserExternalOwnerId: hubspotContactProperties.OwnerId,
			Country:             hubspotContactProperties.Country,
			Region:              hubspotContactProperties.State,
			Locality:            hubspotContactProperties.City,
			Address:             hubspotContactProperties.Address,
			Zip:                 hubspotContactProperties.Zip,
			DefaultLocationName: "Default location",
		}
		// set reference to linked organizations
		contactForCustomerOs.OrganizationsExternalIds = utils.ConvertJsonbToStringSlice(v.CompaniesExternalIds)
		// set reference to primary organization
		if hubspotContactProperties.PrimaryCompanyExternalId.Valid {
			contactForCustomerOs.PrimaryOrganizationExternalId = strconv.FormatFloat(hubspotContactProperties.PrimaryCompanyExternalId.Float64, 'f', 0, 64)
		}
		// add primary organization to organizations list
		contactForCustomerOs.OrganizationsExternalIds = append(contactForCustomerOs.OrganizationsExternalIds, contactForCustomerOs.PrimaryOrganizationExternalId)
		// remove any duplicated organizations
		contactForCustomerOs.OrganizationsExternalIds = utils.RemoveDuplicates(contactForCustomerOs.OrganizationsExternalIds)

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

		// set contact's tags
		if isCustomerTag(hubspotContactProperties.LifecycleStage) {
			contactForCustomerOs.TagName = "CUSTOMER"
		} else if isProspectTag(hubspotContactProperties.LifecycleStage) {
			contactForCustomerOs.TagName = "PROSPECT"
		}

		customerOsContacts = append(customerOsContacts, contactForCustomerOs)
		s.contacts[v.Id] = v
	}
	return customerOsContacts
}

func isCustomerTag(hubspotLifecycleStage string) bool {
	customerLifecycleStages := map[string]bool{
		"customer": true}
	return customerLifecycleStages[hubspotLifecycleStage]
}

func isProspectTag(hubspotLifecycleStage string) bool {
	prospectLifecycleStages := map[string]bool{
		"lead": true, "subscriber": true, "marketingqualifiedlead": true, "salesqualifiedlead": true,
		"opportunity": true}
	return prospectLifecycleStages[hubspotLifecycleStage]
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
			Country:              hubspotCompanyProperties.Country,
			Region:               hubspotCompanyProperties.State,
			Locality:             hubspotCompanyProperties.City,
			Address:              hubspotCompanyProperties.Address,
			Address2:             hubspotCompanyProperties.Address2,
			Zip:                  hubspotCompanyProperties.Zip,
			Phone:                hubspotCompanyProperties.Phone,
			OrganizationTypeName: "COMPANY",
			DefaultLocationName:  "Default location",
		})
		s.companies[v.Id] = v
	}
	return customerOsOrganizations
}

func (s *hubspotDataService) GetUsersForSync(batchSize int, runId string) []*entity.UserData {
	hubspotOwners, err := repository.GetOwners(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	customerOsUsers := make([]*entity.UserData, 0, len(hubspotOwners))
	for _, v := range hubspotOwners {
		userData := entity.UserData{
			ExternalId:      strconv.FormatInt(v.UserId, 10),
			ExternalOwnerId: v.Id,
			ExternalSystem:  s.SourceId(),
			FirstName:       v.FirstName,
			LastName:        v.LastName,
			Email:           v.Email,
			CreatedAt:       v.CreateDate.UTC(),
			UpdatedAt:       v.CreateDate.UTC(),
			ExternalSyncId:  v.Id,
		}
		customerOsUsers = append(customerOsUsers, &userData)

		s.owners[userData.ExternalSyncId] = v
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
			ExternalId:          v.Id,
			ExternalSystem:      s.SourceId(),
			CreatedAt:           v.CreateDate.UTC(),
			Html:                hubspotNoteProperties.NoteBody,
			UserExternalOwnerId: hubspotNoteProperties.OwnerId,
		}
		if hubspotNoteProperties.CreatedByUserId.Valid {
			noteForCustomerOs.UserExternalId = strconv.FormatFloat(hubspotNoteProperties.CreatedByUserId.Float64, 'f', 0, 64)
		}
		// set reference to all linked contacts
		noteForCustomerOs.ContactsExternalIds = utils.ConvertJsonbToStringSlice(v.ContactsExternalIds)
		// set reference to all linked companies
		noteForCustomerOs.OrganizationsExternalIds = utils.ConvertJsonbToStringSlice(v.CompaniesExternalIds)

		customerOsNotes = append(customerOsNotes, noteForCustomerOs)
		s.notes[v.Id] = v
	}
	return customerOsNotes
}

func (s *hubspotDataService) GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData {
	hubspotEmails, err := repository.GetEmails(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	customerOsEmails := []entity.EmailMessageData{}
	for _, v := range hubspotEmails {
		hubspotEmailProperties, err := repository.GetEmailProperties(s.getDb(), v.AirbyteAbId, v.AirbyteEmailsHashid)
		if err != nil {
			logrus.Error(err)
			continue
		}
		// set main fields
		emailForCustomerOS := entity.EmailMessageData{
			Html:           hubspotEmailProperties.EmailHtml,
			Subject:        hubspotEmailProperties.EmailSubject,
			CreatedAt:      v.CreateDate.UTC(),
			ExternalId:     v.Id,
			ExternalSystem: s.SourceId(),
			EmailThreadId:  hubspotEmailProperties.EmailThreadId,
			EmailMessageId: hubspotEmailProperties.EmailMessageId,
			FromEmail:      hubspotEmailProperties.EmailFromEmail,
			ToEmail:        emailsStringToArray(hubspotEmailProperties.EmailToEmail),
			CcEmail:        emailsStringToArray(hubspotEmailProperties.EmailCcEmail),
			BccEmail:       emailsStringToArray(hubspotEmailProperties.EmailBccEmail),
			FromFirstName:  hubspotEmailProperties.EmailFromFirstName,
			FromLastName:   hubspotEmailProperties.EmailFromLastName,
		}
		// set user id
		if hubspotEmailProperties.CreatedByUserId.Valid {
			emailForCustomerOS.UserExternalId = strconv.FormatFloat(hubspotEmailProperties.CreatedByUserId.Float64, 'f', 0, 64)
		}
		// set email message direction
		if hubspotEmailProperties.EmailDirection == "INCOMING_EMAIL" {
			emailForCustomerOS.Direction = entity.INBOUND
		} else {
			emailForCustomerOS.Direction = entity.OUTBOUND
		}
		// set reference to all linked contacts
		emailForCustomerOS.ContactsExternalIds = utils.ConvertJsonbToStringSlice(v.ContactsExternalIds)
		customerOsEmails = append(customerOsEmails, emailForCustomerOS)
		s.emails[v.Id] = v
	}
	return customerOsEmails
}

func emailsStringToArray(str string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ";")
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

func (s *hubspotDataService) MarkEmailMessageProcessed(externalId, runId string, synced bool) error {
	email, ok := s.emails[externalId]
	if ok {
		err := repository.MarkEmailProcessed(s.getDb(), email, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking email with external reference %s as synced for hubspot", externalId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) Refresh() {
	err := s.getDb().AutoMigrate(&localEntity.SyncStatusContact{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&localEntity.SyncStatusCompany{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&localEntity.SyncStatusOwner{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&localEntity.SyncStatusNote{})
	if err != nil {
		logrus.Error(err)
	}
	err = s.getDb().AutoMigrate(&localEntity.SyncStatusEmail{})
	if err != nil {
		logrus.Error(err)
	}
}

func (s *hubspotDataService) getDb() *gorm.DB {
	schemaName := s.SourceId()

	if len(s.instance) > 0 {
		schemaName = schemaName + "_" + s.instance
	}
	schemaName = schemaName + "_" + s.tenant
	return s.airbyteStoreDb.GetDBHandler(&config.Context{
		Schema: schemaName,
	})
}

func (s *hubspotDataService) SourceId() string {
	return string(entity.AirbyteSourceHubspot)
}

func (s *hubspotDataService) Close() {
	s.owners = make(map[string]localEntity.Owner)
	s.contacts = make(map[string]localEntity.Contact)
	s.companies = make(map[string]localEntity.Company)
	s.notes = make(map[string]localEntity.Note)
	s.emails = make(map[string]localEntity.Email)
}
