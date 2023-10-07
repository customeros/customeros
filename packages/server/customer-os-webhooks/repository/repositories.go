package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository/postgres/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	Drivers                  Drivers
	CommonRepositories       *commonRepository.Repositories
	SyncRunWebhookRepository repository.SyncRunWebhookRepository
	ExternalSystemRepository ExternalSystemRepository
	UserRepository           UserRepository
	LocationRepository       LocationRepository
	OrganizationRepository   OrganizationRepository
	TenantRepository         TenantRepository
	EmailRepository          EmailRepository
	PhoneNumberRepository    PhoneNumberRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
		CommonRepositories:       commonRepository.InitRepositories(gormDb, driver),
		SyncRunWebhookRepository: repository.NewSyncRunWebhookRepository(gormDb),
	}
	repositories.ExternalSystemRepository = NewExternalSystemRepository(driver)
	repositories.UserRepository = NewUserRepository(driver)
	repositories.LocationRepository = NewLocationRepository(driver)
	repositories.OrganizationRepository = NewOrganizationRepository(driver)
	repositories.TenantRepository = NewTenantRepository(driver)
	repositories.EmailRepository = NewEmailRepository(driver)
	repositories.PhoneNumberRepository = NewPhoneNumberRepository(driver)

	err := gormDb.AutoMigrate(&postgresentity.SyncRunWebhook{})
	if err != nil {
		panic(err)
	}

	return &repositories
}
