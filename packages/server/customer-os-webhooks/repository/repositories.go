package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository/postgres/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	Drivers       Drivers
	neo4jDatabase string

	Neo4jRepositories    *neo4jrepository.Repositories
	PostgresRepositories *postgresRepository.Repositories

	SyncRunWebhookRepository     repository.SyncRunWebhookRepository
	UserRepository               UserRepository
	LocationRepository           LocationRepository
	OrganizationRepository       OrganizationRepository
	ContactRepository            ContactRepository
	TenantRepository             TenantRepository
	EmailRepository              EmailRepository
	PhoneNumberRepository        PhoneNumberRepository
	LogEntryRepository           LogEntryRepository
	InteractionSessionRepository InteractionSessionRepository
	InteractionEventRepository   InteractionEventRepository
	CommentRepository            CommentRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB, neo4jDatabase string) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
		neo4jDatabase:            neo4jDatabase,
		PostgresRepositories:     postgresRepository.InitRepositories(gormDb),
		Neo4jRepositories:        neo4jrepository.InitNeo4jRepositories(driver, neo4jDatabase),
		SyncRunWebhookRepository: repository.NewSyncRunWebhookRepository(gormDb),
	}
	repositories.UserRepository = NewUserRepository(driver)
	repositories.LocationRepository = NewLocationRepository(driver)
	repositories.OrganizationRepository = NewOrganizationRepository(driver)
	repositories.ContactRepository = NewContactRepository(driver, neo4jDatabase)
	repositories.TenantRepository = NewTenantRepository(driver)
	repositories.EmailRepository = NewEmailRepository(driver)
	repositories.PhoneNumberRepository = NewPhoneNumberRepository(driver)
	repositories.LogEntryRepository = NewLogEntryRepository(driver)
	repositories.InteractionSessionRepository = NewInteractionSessionRepository(driver, neo4jDatabase)
	repositories.InteractionEventRepository = NewInteractionEventRepository(driver, neo4jDatabase)
	repositories.CommentRepository = NewCommentRepository(driver, neo4jDatabase)

	err := gormDb.AutoMigrate(&postgresentity.SyncRunWebhook{})
	if err != nil {
		panic(err)
	}

	return &repositories
}
