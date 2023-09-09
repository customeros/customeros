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
	NoteDefaultSyncService             SyncService
	MeetingDefaultSyncService          SyncService
	EmailMessageDefaultSyncService     SyncService
	InteractionEventDefaultSyncService SyncService
}

func InitServices(cfg *config.Config, log logger.Logger, driver *neo4j.DriverWithContext, controlDb *gorm.DB, airbyteStoreDb *config.RawDataStoreDB, grpcClients *grpc_client.Clients) *Services {
	repositories := repository.InitRepos(driver, controlDb, airbyteStoreDb)

	services := new(Services)

	services.InitService = NewInitService(repositories, services, log)
	services.OrganizationService = NewOrganizationService(repositories)
	services.ContactService = NewContactService(repositories)
	services.UserService = NewUserService(repositories)

	services.UserDefaultSyncService = NewDefaultUserSyncService(repositories, log)
	services.OrganizationDefaultSyncService = NewDefaultOrganizationSyncService(repositories, services, log)
	services.ContactDefaultSyncService = NewDefaultContactSyncService(repositories, services, log)
	services.IssueDefaultSyncService = NewDefaultIssueSyncService(repositories, services, log)
	services.NoteDefaultSyncService = NewDefaultNoteSyncService(repositories, log)
	services.MeetingDefaultSyncService = NewDefaultMeetingSyncService(repositories, services, log)
	services.EmailMessageDefaultSyncService = NewDefaultEmailMessageSyncService(repositories, services, log)
	services.InteractionEventDefaultSyncService = NewDefaultInteractionEventSyncService(repositories, services, log)

	services.SyncToEventStoreService = NewSyncToEventStoreService(repositories, services, grpcClients, log)
	services.SyncCustomerOsDataService = NewSyncCustomerOsDataService(repositories, services, cfg, log)
	return services
}
