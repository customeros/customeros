package slack

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	sourceentity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

const (
	UsersTableSuffix           = "users"
	ContactsTableSuffix        = "users"
	ChannelMessagesTableSuffix = "channel_messages"
	ThreadMessagesTableSuffix  = "thread_messages"
)

var sourceTableSuffixByDataType = map[string][]string{
	string(common.USERS):              {UsersTableSuffix},
	string(common.CONTACTS):           {ContactsTableSuffix},
	string(common.INTERACTION_EVENTS): {ChannelMessagesTableSuffix, ThreadMessagesTableSuffix},
}

type slackDataService struct {
	rawDataStoreDb *config.RawDataStoreDB
	tenant         string
	instance       string
	processingIds  map[string]source.ProcessingEntity
	dataFuncs      map[common.SyncedEntityType]func(context.Context, int, string) []any
	log            logger.Logger
}

func NewSlackDataService(rawDataStoreDb *config.RawDataStoreDB, tenant string, log logger.Logger) source.SourceDataService {
	dataService := slackDataService{
		rawDataStoreDb: rawDataStoreDb,
		tenant:         tenant,
		processingIds:  map[string]source.ProcessingEntity{},
		log:            log,
	}
	dataService.dataFuncs = map[common.SyncedEntityType]func(context.Context, int, string) []any{}
	dataService.dataFuncs[common.USERS] = dataService.GetUsersForSync
	dataService.dataFuncs[common.CONTACTS] = dataService.GetContactsForSync
	dataService.dataFuncs[common.INTERACTION_EVENTS] = dataService.GetInteractionEventsForSync
	return &dataService
}

func (s *slackDataService) GetDataForSync(ctx context.Context, dataType common.SyncedEntityType, batchSize int, runId string) []interface{} {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SlackDataService.GetDataForSync")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)
	span.LogFields(log.String("dataType", string(dataType)), log.Int("batchSize", batchSize))

	if ok := s.dataFuncs[dataType]; ok != nil {
		return s.dataFuncs[dataType](ctx, batchSize, runId)
	} else {
		s.log.Infof("No %s data function for %s", s.SourceId(), dataType)
		return nil
	}
}

func (s *slackDataService) Init() {
	err := s.getDb().AutoMigrate(&sourceentity.SyncStatusForOpenline{})
	if err != nil {
		s.log.Error(err)
	}
}

func (s *slackDataService) getDb() *gorm.DB {
	schemaName := s.SourceId()

	if len(s.instance) > 0 {
		schemaName = schemaName + "_" + s.instance
	}
	schemaName = schemaName + "_" + s.tenant
	return s.rawDataStoreDb.GetDBHandler(&config.Context{
		Schema: schemaName,
	})
}

func (s *slackDataService) SourceId() string {
	return string(entity.OpenlineSourceSlack)
}

func (s *slackDataService) Close() {
	s.processingIds = make(map[string]source.ProcessingEntity)
}

func (s *slackDataService) GetUsersForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.USERS)
	var users []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		rawRecords, err := repository.GetOpenlineUnprocessedRawRecords(ctx, s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range rawRecords {
			if len(users) >= batchSize {
				break
			}
			outputJSON, err := MapUser(v.Data)
			user, err := source.MapJsonToUser(outputJSON, v.RawId, s.SourceId())
			if err != nil {
				user = entity.UserData{
					BaseData: entity.BaseData{
						SyncId: v.RawId,
					},
				}
			}
			s.processingIds[v.RawId] = source.ProcessingEntity{
				ExternalId:  user.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			users = append(users, user)
		}
	}
	return users
}

func (s *slackDataService) GetContactsForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.CONTACTS)

	var contacts []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		rawRecords, err := repository.GetOpenlineUnprocessedRawRecords(ctx, s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range rawRecords {
			if len(contacts) >= batchSize {
				break
			}
			outputJSON, err := MapContact(v.Data)
			contact, err := source.MapJsonToContact(outputJSON, v.RawId, s.SourceId())
			if err != nil {
				contact = entity.ContactData{
					BaseData: entity.BaseData{
						SyncId: v.RawId,
					},
				}
			}

			s.processingIds[v.RawId] = source.ProcessingEntity{
				ExternalId:  contact.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			contacts = append(contacts, contact)
		}
	}
	return contacts
}

func (s *slackDataService) GetInteractionEventsForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.INTERACTION_EVENTS)

	var interactionEvents []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		rawRecords, err := repository.GetOpenlineUnprocessedRawRecords(ctx, s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix)
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range rawRecords {
			if len(interactionEvents) >= batchSize {
				break
			}
			outputJSON, err := MapInteractionEvent(v.Data)
			interactionEvent, err := source.MapJsonToInteractionEvent(outputJSON, v.RawId, s.SourceId())
			if err != nil {
				interactionEvent = entity.InteractionEventData{
					BaseData: entity.BaseData{
						SyncId: v.RawId,
					},
				}
			}

			s.processingIds[v.RawId] = source.ProcessingEntity{
				ExternalId:  interactionEvent.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			interactionEvents = append(interactionEvents, interactionEvent)
		}
	}
	return interactionEvents
}

func (s *slackDataService) MarkProcessed(ctx context.Context, syncId, runId string, synced, skipped bool, reason string) error {
	v, ok := s.processingIds[syncId]
	if ok {
		err := repository.MarkOpenlineRawRecordProcessed(ctx, s.getDb(), v.Entity, v.TableSuffix, syncId, synced, skipped, runId, v.ExternalId, reason)
		if err != nil {
			s.log.Errorf("error while marking %s with external reference %s as synced for %s", v.Entity, v.ExternalId, s.SourceId())
		}
		return err
	}
	return nil
}
