package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"gorm.io/gorm"
)

type Services struct {
	SyncCustomerOsDataService   SyncCustomerOsDataService
	SyncToEventStoreService     SyncToEventStoreService
	InitService                 InitService
	UserSyncService             UserSyncService
	OrganizationSyncService     OrganizationSyncService
	ContactSyncService          ContactSyncService
	IssueSyncService            IssueSyncService
	NoteSyncService             NoteSyncService
	MeetingSyncService          MeetingSyncService
	InteractionEventSyncService InteractionEventSyncService
}

func InitServices(driver *neo4j.DriverWithContext, controlDb *gorm.DB, airbyteStoreDb *config.AirbyteStoreDB, grpcClients *grpc_client.Clients) *Services {
	repositories := repository.InitRepos(driver, controlDb, airbyteStoreDb)

	services := new(Services)

	services.SyncCustomerOsDataService = NewSyncCustomerOsDataService(repositories, services)
	services.SyncToEventStoreService = NewSyncToEventStoreService(repositories, services, grpcClients)
	services.InitService = NewInitService(repositories, services)
	services.UserSyncService = NewUserSyncService(repositories)
	services.OrganizationSyncService = NewOrganizationSyncService(repositories)
	services.ContactSyncService = NewContactSyncService(repositories)
	services.IssueSyncService = NewIssueSyncService(repositories)
	services.NoteSyncService = NewNoteSyncService(repositories)
	services.MeetingSyncService = NewMeetingSyncService(repositories)
	services.InteractionEventSyncService = NewInteractionEventSyncService(repositories)

	return services
}
