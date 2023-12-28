package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"gorm.io/gorm"
)

type Services struct {
	SyncCustomerOsDataService          SyncCustomerOsDataService
	SyncToEventStoreService            SyncToEventStoreService
	InitService                        InitService
	OrganizationService                OrganizationService
	ContactService                     ContactService
	UserService                        UserService
	UserDefaultSyncService             SyncService
	OrganizationDefaultSyncService     SyncService
	ContactDefaultSyncService          SyncService
	IssueDefaultSyncService            SyncService
	LogEntryDefaultSyncService         SyncService
	MeetingDefaultSyncService          SyncService
	EmailMessageDefaultSyncService     SyncService
	InteractionEventDefaultSyncService SyncService
}

func InitServices(cfg *config.Config, log logger.Logger, driver *neo4j.DriverWithContext, controlDb *gorm.DB, airbyteStoreDb *config.RawDataStoreDB, grpcClients *grpc_client.Clients) *Services {
	repositories := repository.InitRepos(driver, controlDb, airbyteStoreDb, log)

	services := new(Services)

	services.InitService = NewInitService(repositories, services, log)
	services.OrganizationService = NewOrganizationService(repositories)
	services.ContactService = NewContactService(repositories)
	services.UserService = NewUserService(repositories)

	services.UserDefaultSyncService = NewDefaultUserSyncService(repositories, cfg, log)
	services.OrganizationDefaultSyncService = NewDefaultOrganizationSyncService(repositories, cfg, log)
	services.ContactDefaultSyncService = NewDefaultContactSyncService(repositories, cfg, log)
	services.IssueDefaultSyncService = NewDefaultIssueSyncService(repositories, cfg, log)
	services.LogEntryDefaultSyncService = NewDefaultLogEntrySyncService(repositories, cfg, log)
	services.MeetingDefaultSyncService = NewDefaultMeetingSyncService(repositories, services, log)
	services.EmailMessageDefaultSyncService = NewDefaultEmailMessageSyncService(repositories, services, log)
	services.InteractionEventDefaultSyncService = NewDefaultInteractionEventSyncService(repositories, cfg, log)

	services.SyncToEventStoreService = NewSyncToEventStoreService(repositories, services, grpcClients, log)
	services.SyncCustomerOsDataService = NewSyncCustomerOsDataService(repositories, services, cfg, log)
	return services
}
