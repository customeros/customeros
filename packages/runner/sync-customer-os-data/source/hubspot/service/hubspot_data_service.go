package service

import (
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	sourceEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/hubspot"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type hubspotDataService struct {
	airbyteStoreDb *config.AirbyteStoreDB
	tenant         string
	instance       string
	contactsRaw    map[string]string
	companiesRaw   map[string]string
	ownersRaw      map[string]string
	notesRaw       map[string]string
	emailsRaw      map[string]string
	meetingsRaw    map[string]string
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.SourceDataService {
	return &hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		contactsRaw:    map[string]string{},
		companiesRaw:   map[string]string{},
		ownersRaw:      map[string]string{},
		notesRaw:       map[string]string{},
		emailsRaw:      map[string]string{},
		meetingsRaw:    map[string]string{},
	}
}

func (s *hubspotDataService) Refresh() {
	err := s.getDb().AutoMigrate(&sourceEntity.SyncStatus{})
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
	s.ownersRaw = make(map[string]string)
	s.contactsRaw = make(map[string]string)
	s.companiesRaw = make(map[string]string)
	s.notesRaw = make(map[string]string)
	s.emailsRaw = make(map[string]string)
	s.meetingsRaw = make(map[string]string)
}

func (s *hubspotDataService) GetContactsForSync(batchSize int, runId string) []entity.ContactData {
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, repository.ContactEntity)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	var contacts []entity.ContactData
	for _, v := range airbyteRecords {
		contact := entity.ContactData{}
		outputJSON, err := hubspot.MapContact(v.AirbyteData)
		if err != nil {
			logrus.Error(err)
			continue
		}
		err = json.Unmarshal([]byte(outputJSON), &contact)
		if err != nil {
			logrus.Error(err)
			continue
		}
		contact.ExternalSyncId = contact.ExternalId
		contact.ExternalSystem = s.SourceId()
		for _, textCustomField := range contact.TextCustomFields {
			textCustomField.ExternalSystem = s.SourceId()
		}
		contact.Id = ""

		s.contactsRaw[contact.ExternalSyncId] = v.AirbyteAbId
		contacts = append(contacts, contact)
	}
	return contacts
}

func (s *hubspotDataService) GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData {
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, repository.CompanyEntity)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	var organizations []entity.OrganizationData
	for _, v := range airbyteRecords {
		organization := entity.OrganizationData{}
		outputJSON, err := hubspot.MapOrganization(v.AirbyteData)
		if err != nil {
			logrus.Error(err)
			continue
		}
		err = json.Unmarshal([]byte(outputJSON), &organization)
		if err != nil {
			logrus.Error(err)
			continue
		}
		organization.ExternalSyncId = organization.ExternalId
		organization.ExternalSystem = s.SourceId()
		organization.Id = ""

		s.companiesRaw[organization.ExternalSyncId] = v.AirbyteAbId
		organizations = append(organizations, organization)
	}
	return organizations
}

func (s *hubspotDataService) GetUsersForSync(batchSize int, runId string) []entity.UserData {
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, repository.OwnerEntity)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	var users []entity.UserData
	for _, v := range airbyteRecords {
		user := entity.UserData{}
		outputJSON, err := hubspot.MapUser(v.AirbyteData)
		if err != nil {
			logrus.Error(err)
			continue
		}
		err = json.Unmarshal([]byte(outputJSON), &user)
		if err != nil {
			logrus.Error(err)
			continue
		}
		user.ExternalSyncId = user.ExternalId
		user.ExternalSystem = s.SourceId()
		user.Id = ""

		s.ownersRaw[user.ExternalSyncId] = v.AirbyteAbId
		users = append(users, user)
	}
	return users
}

func (s *hubspotDataService) GetNotesForSync(batchSize int, runId string) []entity.NoteData {
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, repository.NoteEntity)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	var notes []entity.NoteData
	for _, v := range airbyteRecords {
		note := entity.NoteData{}
		outputJSON, err := hubspot.MapNote(v.AirbyteData)
		if err != nil {
			logrus.Error(err)
			continue
		}
		err = json.Unmarshal([]byte(outputJSON), &note)
		if err != nil {
			logrus.Error(err)
			continue
		}
		note.ExternalSyncId = note.ExternalId
		note.ExternalSystem = s.SourceId()
		note.Id = ""

		s.notesRaw[note.ExternalSyncId] = v.AirbyteAbId
		notes = append(notes, note)
	}
	return notes
}

func (s *hubspotDataService) GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData {
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, repository.EmailMessageEntity)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	var emailMessages []entity.EmailMessageData
	for _, v := range airbyteRecords {
		emailMessage := entity.EmailMessageData{}
		outputJSON, err := hubspot.MapEmailMessage(v.AirbyteData)
		if err != nil {
			logrus.Error(err)
			continue
		}
		err = json.Unmarshal([]byte(outputJSON), &emailMessage)
		if err != nil {
			logrus.Error(err)
			continue
		}
		emailMessage.ExternalSyncId = emailMessage.ExternalId
		emailMessage.ExternalSystem = s.SourceId()
		emailMessage.Id = ""

		s.emailsRaw[emailMessage.ExternalSyncId] = v.AirbyteAbId
		emailMessages = append(emailMessages, emailMessage)
	}
	return emailMessages
}

func (s *hubspotDataService) GetIssuesForSync(batchSize int, runId string) []entity.IssueData {
	// no need to implement
	return nil
}

func (s *hubspotDataService) GetMeetingsForSync(batchSize int, runId string) []entity.MeetingData {
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, repository.MeetingEntity)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	var meetings []entity.MeetingData
	for _, v := range airbyteRecords {
		meeting := entity.MeetingData{}
		outputJSON, err := hubspot.MapMeeting(v.AirbyteData)
		if err != nil {
			logrus.Error(err)
			continue
		}
		err = json.Unmarshal([]byte(outputJSON), &meeting)
		if err != nil {
			logrus.Error(err)
			continue
		}
		meeting.ExternalSyncId = meeting.ExternalId
		meeting.ExternalSystem = s.SourceId()
		meeting.Id = ""

		s.meetingsRaw[meeting.ExternalSyncId] = v.AirbyteAbId
		meetings = append(meetings, meeting)
	}
	return meetings
}

func (s *hubspotDataService) MarkContactProcessed(externalSyncId, runId string, synced bool) error {
	airbyteAbId, ok := s.contactsRaw[externalSyncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), repository.ContactEntity, airbyteAbId, synced, runId, externalSyncId)
		if err != nil {
			logrus.Errorf("error while marking contact with external reference %s as synced for hubspot", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkOrganizationProcessed(externalSyncId, runId string, synced bool) error {
	airbyteAbId, ok := s.companiesRaw[externalSyncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), repository.CompanyEntity, airbyteAbId, synced, runId, externalSyncId)
		if err != nil {
			logrus.Errorf("error while marking company with external reference %s as synced for hubspot", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkUserProcessed(externalSyncId, runId string, synced bool) error {
	airbyteAbId, ok := s.ownersRaw[externalSyncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), repository.OwnerEntity, airbyteAbId, synced, runId, externalSyncId)
		if err != nil {
			logrus.Errorf("error while marking owner with external reference %s as synced for hubspot", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkNoteProcessed(externalSyncId, runId string, synced bool) error {
	airbyteAbId, ok := s.notesRaw[externalSyncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), repository.NoteEntity, airbyteAbId, synced, runId, externalSyncId)
		if err != nil {
			logrus.Errorf("error while marking note with external reference %s as synced for hubspot", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkIssueProcessed(externalSyncId, runId string, synced bool) error {
	// no need to implement
	return nil
}

func (s *hubspotDataService) MarkEmailMessageProcessed(externalSyncId, runId string, synced bool) error {
	airbyteAbId, ok := s.emailsRaw[externalSyncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), repository.EmailMessageEntity, airbyteAbId, synced, runId, externalSyncId)
		if err != nil {
			logrus.Errorf("error while marking email with external reference %s as synced for hubspot", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkMeetingProcessed(externalSyncId, runId string, synced bool) error {
	airbyteAbId, ok := s.meetingsRaw[externalSyncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), repository.MeetingEntity, airbyteAbId, synced, runId, externalSyncId)
		if err != nil {
			logrus.Errorf("error while marking meeting with external reference %s as synced for hubspot", externalSyncId)
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) GetInteractionEventsForSync(batchSize int, runId string) []entity.InteractionEventData {
	return nil
}

func (s *hubspotDataService) MarkInteractionEventProcessed(externalSyncId, runId string, synced bool) error {
	return nil
}
