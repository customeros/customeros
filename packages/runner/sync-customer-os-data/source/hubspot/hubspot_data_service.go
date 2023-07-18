package hubspot

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	sourceEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	ContactsTableSuffix      = "contacts"
	CompaniesTableSuffix     = "companies"
	OwnersTableSuffix        = "owners"
	NotesTableSuffix         = "engagements_notes"
	MeetingsTableSuffix      = "engagements_meetings"
	EmailMessagesTableSuffix = "engagements_emails"
)

var sourceTableSuffixByDataType = map[string][]string{
	string(common.USERS):          {OwnersTableSuffix},
	string(common.CONTACTS):       {ContactsTableSuffix},
	string(common.ORGANIZATIONS):  {CompaniesTableSuffix},
	string(common.NOTES):          {NotesTableSuffix},
	string(common.MEETINGS):       {MeetingsTableSuffix},
	string(common.EMAIL_MESSAGES): {EmailMessagesTableSuffix},
}

type hubspotDataService struct {
	airbyteStoreDb *config.AirbyteStoreDB
	tenant         string
	instance       string
	processingIds  map[string]source.ProcessingEntity
	dataFuncs      map[common.SyncedEntityType]func(int, string) []any
}

func NewHubspotDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) source.SourceDataService {
	dataService := hubspotDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		processingIds:  map[string]source.ProcessingEntity{},
	}
	dataService.dataFuncs = map[common.SyncedEntityType]func(int, string) []any{}
	dataService.dataFuncs[common.USERS] = dataService.GetUsersForSync
	dataService.dataFuncs[common.ORGANIZATIONS] = dataService.GetOrganizationsForSync
	dataService.dataFuncs[common.CONTACTS] = dataService.GetContactsForSync
	dataService.dataFuncs[common.NOTES] = dataService.GetNotesForSync
	dataService.dataFuncs[common.MEETINGS] = dataService.GetMeetingsForSync
	dataService.dataFuncs[common.EMAIL_MESSAGES] = dataService.GetEmailMessagesForSync
	return &dataService
}

func (s *hubspotDataService) GetDataForSync(dataType common.SyncedEntityType, batchSize int, runId string) []interface{} {
	if ok := s.dataFuncs[dataType]; ok != nil {
		return s.dataFuncs[dataType](batchSize, runId)
	} else {
		logrus.Warnf("No %s data function for %s", s.SourceId(), dataType)
		return nil
	}
}

func (s *hubspotDataService) Init() {
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
	s.processingIds = make(map[string]source.ProcessingEntity)
}

func (s *hubspotDataService) GetUsersForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.USERS)
	var users []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			logrus.Panic(err) // alexb handle errors
			return nil
		}
		for _, v := range airbyteRecords {
			if len(users) >= batchSize {
				break
			}
			outputJSON, err := MapUser(v.AirbyteData)
			user, err := source.MapJsonToUser(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				logrus.Panic(err) // alexb handle errors
				continue
			}
			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  user.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			users = append(users, user)
		}
	}
	return users
}

func (s *hubspotDataService) GetOrganizationsForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.ORGANIZATIONS)

	var organizations []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			logrus.Panic(err) // alexb handle errors
			return nil
		}
		for _, v := range airbyteRecords {
			if len(organizations) >= batchSize {
				break
			}
			outputJSON, err := MapOrganization(v.AirbyteData)
			organization, err := source.MapJsonToOrganization(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				logrus.Panic(err) // alexb handle errors
				continue
			}

			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  organization.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			organizations = append(organizations, organization)
		}
	}
	return organizations
}

func (s *hubspotDataService) GetContactsForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.CONTACTS)

	var contacts []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			logrus.Panic(err) // alexb handle errors
			return nil
		}
		for _, v := range airbyteRecords {
			if len(contacts) >= batchSize {
				break
			}
			outputJSON, err := MapContact(v.AirbyteData)
			contact, err := source.MapJsonToContact(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				logrus.Panic(err) // alexb handle errors
				continue
			}

			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  contact.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			contacts = append(contacts, contact)
		}
	}
	return contacts
}

func (s *hubspotDataService) GetNotesForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.NOTES)
	var notes []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			logrus.Panic(err) // alexb handle errors
			return nil
		}
		for _, v := range airbyteRecords {
			if len(notes) >= batchSize {
				break
			}
			outputJSON, err := MapNote(v.AirbyteData)
			note, err := source.MapJsonToNote(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				logrus.Panic(err) // alexb handle errors
				continue
			}

			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  note.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			notes = append(notes, note)
		}
	}
	return notes
}

func (s *hubspotDataService) GetEmailMessagesForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.EMAIL_MESSAGES)
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffixByDataType[currentEntity][0])
	if err != nil {
		logrus.Panic(err) // alexb handle errors
		return nil
	}
	var emailMessages []any
	for _, v := range airbyteRecords {
		if len(emailMessages) >= batchSize {
			break
		}
		outputJSON, err := MapEmailMessage(v.AirbyteData)
		emailMessage, err := source.MapJsonToEmailMessage(outputJSON, v.AirbyteAbId, s.SourceId())
		if err != nil {
			logrus.Panic(err) // alexb handle errors
			continue
		}

		s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
			ExternalId:  emailMessage.ExternalId,
			Entity:      currentEntity,
			TableSuffix: sourceTableSuffixByDataType[currentEntity][0],
		}
		emailMessages = append(emailMessages, emailMessage)
	}
	return emailMessages
}

func (s *hubspotDataService) GetMeetingsForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.MEETINGS)
	airbyteRecords, err := repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffixByDataType[currentEntity][0])
	if err != nil {
		logrus.Panic(err) // alexb handle errors
		return nil
	}
	var meetings []any
	for _, v := range airbyteRecords {
		if len(meetings) >= batchSize {
			break
		}
		outputJSON, err := MapMeeting(v.AirbyteData)
		meeting, err := source.MapJsonToMeeting(outputJSON, v.AirbyteAbId, s.SourceId())
		if err != nil {
			logrus.Panic(err) // alexb handle errors
			continue
		}

		s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
			ExternalId:  meeting.ExternalId,
			Entity:      currentEntity,
			TableSuffix: sourceTableSuffixByDataType[currentEntity][0],
		}
		meetings = append(meetings, meeting)
	}
	return meetings
}

func (s *hubspotDataService) MarkProcessed(syncId, runId string, synced, skipped bool, reason string) error {
	v, ok := s.processingIds[syncId]
	if ok {
		err := repository.MarkProcessed(s.getDb(), v.Entity, v.TableSuffix, syncId, synced, skipped, runId, v.ExternalId, reason)
		if err != nil {
			logrus.Errorf("error while marking %s with external reference %s as synced for %s", v.Entity, v.ExternalId, s.SourceId())
		}
		return err
	}
	return nil
}
