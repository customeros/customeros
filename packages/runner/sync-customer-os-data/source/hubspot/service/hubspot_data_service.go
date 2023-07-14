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

const (
	ContactEntity      = "contacts"
	CompanyEntity      = "companies"
	OwnerEntity        = "owners"
	NoteEntity         = "engagements_notes"
	MeetingEntity      = "engagements_meetings"
	EmailMessageEntity = "engagements_emails"
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
	_, ok := s.processingIds[ContactEntity]
	if !ok {
		s.processingIds[ContactEntity] = map[string]string{}
	}
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, ContactEntity)
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
		contact.SyncId = v.AirbyteAbId
		contact.ExternalSystem = s.SourceId()
		for _, textCustomField := range contact.TextCustomFields {
			textCustomField.ExternalSystem = s.SourceId()
		}
		contact.Id = ""

		s.processingIds[ContactEntity][contact.SyncId] = contact.ExternalId
		contacts = append(contacts, contact)
	}
	return contacts
}

func (s *hubspotDataService) GetUsersForSync(batchSize int, runId string) []entity.UserData {
	_, ok := s.processingIds[OwnerEntity]
	if !ok {
		s.processingIds[OwnerEntity] = map[string]string{}
	}
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, OwnerEntity)
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
		user.SyncId = v.AirbyteAbId
		user.ExternalSystem = s.SourceId()
		user.Id = ""

		s.processingIds[OwnerEntity][user.SyncId] = user.ExternalId
		users = append(users, user)
	}
	return users
}

func (s *hubspotDataService) GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData {
	_, ok := s.processingIds[CompanyEntity]
	if !ok {
		s.processingIds[CompanyEntity] = map[string]string{}
	}
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, CompanyEntity)
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
		organization.SyncId = v.AirbyteAbId
		organization.ExternalSystem = s.SourceId()
		organization.Id = ""

		s.processingIds[CompanyEntity][organization.SyncId] = organization.ExternalId
		organizations = append(organizations, organization)
	}
	return organizations
}

func (s *hubspotDataService) GetNotesForSync(batchSize int, runId string) []entity.NoteData {
	_, ok := s.processingIds[NoteEntity]
	if !ok {
		s.processingIds[NoteEntity] = map[string]string{}
	}
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, NoteEntity)
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
		note.SyncId = v.AirbyteAbId
		note.ExternalSystem = s.SourceId()
		note.Id = ""

		s.processingIds[NoteEntity][note.SyncId] = note.ExternalId
		notes = append(notes, note)
	}
	return notes
}

func (s *hubspotDataService) GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData {
	_, ok := s.processingIds[EmailMessageEntity]
	if !ok {
		s.processingIds[EmailMessageEntity] = map[string]string{}
	}
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, EmailMessageEntity)
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
		emailMessage.SyncId = v.AirbyteAbId
		emailMessage.ExternalSystem = s.SourceId()
		emailMessage.Id = ""

		s.processingIds[EmailMessageEntity][emailMessage.SyncId] = emailMessage.ExternalId
		emailMessages = append(emailMessages, emailMessage)
	}
	return emailMessages
}

func (s *hubspotDataService) GetIssuesForSync(batchSize int, runId string) []entity.IssueData {
	// no need to implement
	return nil
}

func (s *hubspotDataService) GetMeetingsForSync(batchSize int, runId string) []entity.MeetingData {
	_, ok := s.processingIds[MeetingEntity]
	if !ok {
		s.processingIds[MeetingEntity] = map[string]string{}
	}
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, MeetingEntity)
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
		meeting.SyncId = v.AirbyteAbId
		meeting.ExternalSystem = s.SourceId()
		meeting.Id = ""

		s.processingIds[MeetingEntity][meeting.SyncId] = meeting.ExternalId
		meetings = append(meetings, meeting)
	}
	return meetings
}

func (s *hubspotDataService) MarkProcessed(processingEntity, syncId, runId string, synced bool) error {
	externalId, ok := s.processingIds[processingEntity][syncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), processingEntity, syncId, synced, runId, externalId)
		if err != nil {
			logrus.Errorf("error while marking %s with external reference %s as synced for %s", processingEntity, externalId, s.SourceId())
		}
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkContactProcessed(syncId, runId string, synced bool) error {
	err := s.MarkProcessed(ContactEntity, syncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkOrganizationProcessed(syncId, runId string, synced bool) error {
	err := s.MarkProcessed(CompanyEntity, syncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkUserProcessed(syncId, runId string, synced bool) error {
	err := s.MarkProcessed(OwnerEntity, syncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkNoteProcessed(syncId, runId string, synced bool) error {
	err := s.MarkProcessed(NoteEntity, syncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkIssueProcessed(syncId, runId string, synced bool) error {
	// no need to implement
	return nil
}

func (s *hubspotDataService) MarkEmailMessageProcessed(syncId, runId string, synced bool) error {
	err := s.MarkProcessed(EmailMessageEntity, syncId, runId, synced)
	if err != nil {
		return err
	}
	return nil
}

func (s *hubspotDataService) MarkMeetingProcessed(syncId, runId string, synced bool) error {
	err := s.MarkProcessed(MeetingEntity, syncId, runId, synced)
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
