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
	processingIds  map[string]map[string]string
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.SourceDataService {
	return &hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		processingIds:  map[string]map[string]string{},
	}
}

func (s *hubspotDataService) Start() {
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
	s.processingIds = make(map[string]map[string]string)
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
		contact.ExternalSyncId = v.AirbyteAbId
		contact.ExternalSystem = s.SourceId()
		for _, textCustomField := range contact.TextCustomFields {
			textCustomField.ExternalSystem = s.SourceId()
		}
		contact.Id = ""

		s.processingIds[repository.ContactEntity][contact.ExternalSyncId] = contact.ExternalId
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
		organization.ExternalSyncId = v.AirbyteAbId
		organization.ExternalSystem = s.SourceId()
		organization.Id = ""

		s.processingIds[repository.CompanyEntity][organization.ExternalSyncId] = organization.ExternalId
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
		user.ExternalSyncId = v.AirbyteAbId
		user.ExternalSystem = s.SourceId()
		user.Id = ""

		s.processingIds[repository.OwnerEntity][user.ExternalSyncId] = user.ExternalId
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
		note.ExternalSyncId = v.AirbyteAbId
		note.ExternalSystem = s.SourceId()
		note.Id = ""

		s.processingIds[repository.NoteEntity][note.ExternalSyncId] = note.ExternalId
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
		emailMessage.ExternalSyncId = v.AirbyteAbId
		emailMessage.ExternalSystem = s.SourceId()
		emailMessage.Id = ""

		s.processingIds[repository.EmailMessageEntity][emailMessage.ExternalSyncId] = emailMessage.ExternalId
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
		meeting.ExternalSyncId = v.AirbyteAbId
		meeting.ExternalSystem = s.SourceId()
		meeting.Id = ""

		s.processingIds[repository.MeetingEntity][meeting.ExternalSyncId] = meeting.ExternalId
		meetings = append(meetings, meeting)
	}
	return meetings
}

func (s *hubspotDataService) MarkProcessed(processingEntity, externalSyncId, runId string, synced bool) error {
	externalId, ok := s.processingIds[processingEntity][externalSyncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), processingEntity, externalSyncId, synced, runId, externalId)
		if err != nil {
			logrus.Errorf("error while marking %s with external reference %s as synced for %s", processingEntity, externalId, s.SourceId())
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkContactProcessed(externalSyncId, runId string, synced bool) error {
	err := s.MarkProcessed(repository.ContactEntity, externalSyncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkOrganizationProcessed(externalSyncId, runId string, synced bool) error {
	err := s.MarkProcessed(repository.CompanyEntity, externalSyncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkUserProcessed(externalSyncId, runId string, synced bool) error {
	err := s.MarkProcessed(repository.OwnerEntity, externalSyncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkNoteProcessed(externalSyncId, runId string, synced bool) error {
	err := s.MarkProcessed(repository.NoteEntity, externalSyncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkIssueProcessed(externalSyncId, runId string, synced bool) error {
	// no need to implement
	return nil
}

func (s *hubspotDataService) MarkEmailMessageProcessed(externalSyncId, runId string, synced bool) error {
	err := s.MarkProcessed(repository.EmailMessageEntity, externalSyncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkMeetingProcessed(externalSyncId, runId string, synced bool) error {
	err := s.MarkProcessed(repository.MeetingEntity, externalSyncId, runId, synced)
	if err != nil {
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
