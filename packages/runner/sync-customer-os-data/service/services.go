package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"gorm.io/gorm"
)

type Services struct {
	SyncCustomerOsDataService          SyncCustomerOsDataService
	SyncToEventStoreService            SyncToEventStoreService
	InitService                        InitService
	OrganizationService                OrganizationService
	UserDefaultSyncService             SyncService
	OrganizationDefaultSyncService     SyncService
	ContactDefaultSyncService          SyncService
	IssueDefaultSyncService            SyncService
	NoteDefaultSyncService             SyncService
	MeetingDefaultSyncService          SyncService
	EmailMessageDefaultSyncService     SyncService
	InteractionEventDefaultSyncService SyncService
}

func InitServices(cfg *config.Config, driver *neo4j.DriverWithContext, controlDb *gorm.DB, airbyteStoreDb *config.AirbyteStoreDB, grpcClients *grpc_client.Clients) *Services {
	repositories := repository.InitRepos(driver, controlDb, airbyteStoreDb)

	services := new(Services)

	services.InitService = NewInitService(repositories, services)
	services.OrganizationService = NewOrganizationService(repositories)

	services.UserDefaultSyncService = NewDefaultUserSyncService(repositories)
	services.OrganizationDefaultSyncService = NewDefaultOrganizationSyncService(repositories, services)
	services.ContactDefaultSyncService = NewDefaultContactSyncService(repositories, services)
	services.IssueDefaultSyncService = NewDefaultIssueSyncService(repositories, services)
	services.NoteDefaultSyncService = NewDefaultNoteSyncService(repositories)
	services.MeetingDefaultSyncService = NewDefaultMeetingSyncService(repositories, services)
	services.EmailMessageDefaultSyncService = NewDefaultEmailMessageSyncService(repositories, services)
	services.InteractionEventDefaultSyncService = NewDefaultInteractionEventSyncService(repositories)

	services.SyncToEventStoreService = NewSyncToEventStoreService(repositories, services, grpcClients)
	services.SyncCustomerOsDataService = NewSyncCustomerOsDataService(repositories, services, cfg)
	return services
}
