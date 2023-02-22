package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	localEntity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/zendesk_support/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/zendesk_support/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type zendeskSupportDataService struct {
	airbyteStoreDb *config.AirbyteStoreDB
	tenant         string
	instance       string
	users          map[string]localEntity.User
}

func NewZendeskSupportDataService(airbyteStoreDb *config.AirbyteStoreDB, tenant string) common.SourceDataService {
	return &zendeskSupportDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		users:          map[string]localEntity.User{},
	}
}

func (s *zendeskSupportDataService) Refresh() {
	err := s.getDb().AutoMigrate(&localEntity.SyncStatusUser{})
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
	s.users = make(map[string]localEntity.User)
}

func (z *zendeskSupportDataService) SourceId() string {
	return string(entity.AirbyteSourceZendeskSupport)
}

func (z *zendeskSupportDataService) GetContactsForSync(batchSize int, runId string) []entity.ContactData {
	//TODO implement me
	return nil
}

func (z *zendeskSupportDataService) GetOrganizationsForSync(batchSize int, runId string) []entity.OrganizationData {
	//TODO implement me
	return nil
}

func (s *zendeskSupportDataService) GetUsersForSync(batchSize int, runId string) []*entity.UserData {
	zendeskUsers, err := repository.GetUsers(s.getDb(), batchSize, runId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	customerOsUsers := make([]*entity.UserData, 0, len(zendeskUsers))
	for _, v := range zendeskUsers {
		userData := entity.UserData{
			ExternalId:     strconv.FormatInt(v.Id, 10),
			ExternalSystem: s.SourceId(),
			Name:           v.Name,
			Email:          v.Email,
			PhoneNumber:    v.Phone,
			CreatedAt:      v.CreateDate.UTC(),
			UpdatedAt:      v.UpdatedDate.UTC(),
			ExternalSyncId: strconv.FormatInt(v.Id, 10),
		}
		customerOsUsers = append(customerOsUsers, &userData)

		s.users[userData.ExternalSyncId] = v
	}
	return customerOsUsers
}

func (z zendeskSupportDataService) GetNotesForSync(batchSize int, runId string) []entity.NoteData {
	//TODO implement me
	return nil
}

func (z zendeskSupportDataService) GetEmailMessagesForSync(batchSize int, runId string) []entity.EmailMessageData {
	//TODO implement me
	return nil
}

func (z zendeskSupportDataService) MarkContactProcessed(externalId, runId string, synced bool) error {
	//TODO implement me
	return nil
}

func (z zendeskSupportDataService) MarkOrganizationProcessed(externalId, runId string, synced bool) error {
	//TODO implement me
	return nil
}

func (s *zendeskSupportDataService) MarkUserProcessed(externalId, runId string, synced bool) error {
	user, ok := s.users[externalId]
	if ok {
		err := repository.MarkUserProcessed(s.getDb(), user, synced, runId)
		if err != nil {
			logrus.Errorf("error while marking owner with external reference %s as synced for hubspot", externalId)
		}
		return err
	}
	return nil
}

func (z zendeskSupportDataService) MarkNoteProcessed(externalId, runId string, synced bool) error {
	//TODO implement me
	return nil
}

func (z zendeskSupportDataService) MarkEmailMessageProcessed(externalId, runId string, synced bool) error {
	//TODO implement me
	return nil
}
