package zendesk_support

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	sourceEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/entity"
	common_repository "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	UsersTableSuffix          = "users"
	OrganizationsTableSuffix  = "organizations"
	TicketsTableSuffix        = "tickets"
	TicketCommentsTableSuffix = "ticket_comments"
)

var sourceTableSuffixByDataType = map[string][]string{
	string(common.USERS):              {UsersTableSuffix},
	string(common.ORGANIZATIONS):      {OrganizationsTableSuffix, UsersTableSuffix},
	string(common.ISSUES):             {TicketsTableSuffix},
	string(common.NOTES):              {TicketCommentsTableSuffix},
	string(common.INTERACTION_EVENTS): {TicketCommentsTableSuffix},
}

type zendeskSupportDataService struct {
	airbyteStoreDb *config.AirbyteStoreDB
	tenant         string
	instance       string
	processingIds  map[string]source.ProcessingEntity
	dataFuncs      map[common.SyncedEntityType]func(int, string) []any
}

func (s *zendeskSupportDataService) GetDataForSync(dataType common.SyncedEntityType, batchSize int, runId string) []interface{} {
	if ok := s.dataFuncs[dataType]; ok != nil {
		return s.dataFuncs[dataType](batchSize, runId)
	} else {
		logrus.Warnf("No %s data function for %s", s.SourceId(), dataType)
		return nil
	}
}

func NewZendeskSupportDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) source.SourceDataService {
	dataService := zendeskSupportDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		processingIds:  map[string]source.ProcessingEntity{},
	}
	dataService.dataFuncs = map[common.SyncedEntityType]func(int, string) []any{}
	dataService.dataFuncs[common.USERS] = dataService.GetUsersForSync
	dataService.dataFuncs[common.ORGANIZATIONS] = dataService.GetOrganizationsForSync
	dataService.dataFuncs[common.ISSUES] = dataService.GetIssuesForSync
	dataService.dataFuncs[common.NOTES] = dataService.GetNotesForSync
	dataService.dataFuncs[common.INTERACTION_EVENTS] = dataService.GetInteractionEventsForSync
	return &dataService
}

func (s *zendeskSupportDataService) Init() {
	err := s.getDb().AutoMigrate(&sourceEntity.SyncStatus{})
	if err != nil {
		logrus.Error(err)
	}
}

func (s *zendeskSupportDataService) getDb() *gorm.DB {
	schemaName := s.SourceId()

	if len(s.instance) > 0 {
		schemaName = schemaName + "_" + s.instance
	}
	schemaName = schemaName + "_" + s.tenant
	return s.airbyteStoreDb.GetDBHandler(&config.Context{
		Schema: schemaName,
	})
}

func (s *zendeskSupportDataService) Close() {
	s.processingIds = make(map[string]source.ProcessingEntity)
}

func (s *zendeskSupportDataService) SourceId() string {
	return string(entity.AirbyteSourceZendeskSupport)
}

func (s *zendeskSupportDataService) GetContactsForSync(batchSize int, runId string) []entity.ContactData {
	return nil
}

func (s *zendeskSupportDataService) GetUsersForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.USERS)

	var users []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := common_repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
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

func (s *zendeskSupportDataService) GetOrganizationsForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.ORGANIZATIONS)

	var organizations []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := common_repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
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

func (s *zendeskSupportDataService) GetIssuesForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.ISSUES)

	var issues []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := common_repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			logrus.Panic(err) // alexb handle errors
			return nil
		}
		for _, v := range airbyteRecords {
			if len(issues) >= batchSize {
				break
			}
			outputJSON, err := MapIssue(v.AirbyteData)
			issue, err := source.MapJsonToIssue(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				logrus.Panic(err) // alexb handle errors
				continue
			}

			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  issue.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			issues = append(issues, issue)
		}
	}
	return issues
}

func (s *zendeskSupportDataService) GetNotesForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.NOTES)

	var notes []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := common_repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
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

func (s *zendeskSupportDataService) GetInteractionEventsForSync(batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.INTERACTION_EVENTS)

	var interactionEvents []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := common_repository.GetAirbyteUnprocessedRecords(s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			logrus.Panic(err) // alexb handle errors
			return nil
		}
		for _, v := range airbyteRecords {
			if len(interactionEvents) >= batchSize {
				break
			}
			outputJSON, err := MapInteractionEvent(v.AirbyteData)
			interactionEvent, err := source.MapJsonToInteractionEvent(outputJSON, v.AirbyteAbId, s.SourceId())
			if err != nil {
				logrus.Panic(err) // alexb handle errors
				continue
			}

			s.processingIds[v.AirbyteAbId] = source.ProcessingEntity{
				ExternalId:  interactionEvent.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			interactionEvents = append(interactionEvents, interactionEvent)
		}
	}
	return interactionEvents
}

func (s *zendeskSupportDataService) MarkProcessed(syncId, runId string, synced, skipped bool, reason string) error {
	v, ok := s.processingIds[syncId]
	if ok {
		err := common_repository.MarkProcessed(s.getDb(), v.Entity, v.TableSuffix, syncId, synced, skipped, runId, v.ExternalId, reason)
		if err != nil {
			logrus.Errorf("error while marking %s with external reference %s as synced for %s", v.Entity, v.ExternalId, s.SourceId())
		}
		return err
	}
	return nil
}
