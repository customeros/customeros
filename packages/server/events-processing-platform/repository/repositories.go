package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jrepository "github.com/openline-ai/customer-os-neo4j-repository/repository"
	cmn_repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	repository "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"gorm.io/gorm"
)

type Drivers struct {
	Neo4jDriver *neo4j.DriverWithContext
}

type Repositories struct {
	Drivers Drivers

	Neo4jRepositories       *neo4jrepository.Repositories
	CommonRepositories      *cmn_repository.Repositories
	CustomerOsIdsRepository repository.CustomerOsIdsRepository

	ActionRepository             ActionRepository
	CommentRepository            CommentRepository
	ContactRepository            ContactRepository
	ContractRepository           ContractRepository
	CountryRepository            CountryRepository
	CustomFieldRepository        CustomFieldRepository
	EmailRepository              EmailRepository
	ExternalSystemRepository     ExternalSystemRepository
	IssueRepository              IssueRepository
	InteractionEventRepository   InteractionEventRepository
	InteractionSessionRepository InteractionSessionRepository
	JobRoleRepository            JobRoleRepository
	LocationRepository           LocationRepository
	OpportunityRepository        OpportunityRepository
	OrganizationRepository       OrganizationRepository
	PhoneNumberRepository        PhoneNumberRepository
	PlayerRepository             PlayerRepository
	ServiceLineItemRepository    ServiceLineItemRepository
	SocialRepository             SocialRepository
	TagRepository                TagRepository
	TimelineEventRepository      TimelineEventRepository
}

func InitRepos(driver *neo4j.DriverWithContext, neo4jDatabase string, gormDb *gorm.DB, log logger.Logger) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
		Neo4jRepositories:            neo4jrepository.InitNeo4jRepositories(driver, neo4jDatabase),
		CommonRepositories:           cmn_repository.InitRepositories(gormDb, driver),
		CustomerOsIdsRepository:      repository.NewCustomerOsIdsRepository(gormDb),
		PhoneNumberRepository:        NewPhoneNumberRepository(driver),
		EmailRepository:              NewEmailRepository(driver),
		ContactRepository:            NewContactRepository(driver),
		OrganizationRepository:       NewOrganizationRepository(driver, neo4jDatabase),
		LocationRepository:           NewLocationRepository(driver),
		CountryRepository:            NewCountryRepository(driver),
		JobRoleRepository:            NewJobRoleRepository(driver),
		SocialRepository:             NewSocialRepository(driver),
		InteractionEventRepository:   NewInteractionEventRepository(driver, neo4jDatabase),
		InteractionSessionRepository: NewInteractionSessionRepository(driver, neo4jDatabase),
		ActionRepository:             NewActionRepository(driver),
		IssueRepository:              NewIssueRepository(driver, neo4jDatabase),
		TagRepository:                NewTagRepository(driver),
		PlayerRepository:             NewPlayerRepository(driver),
		ExternalSystemRepository:     NewExternalSystemRepository(driver),
		TimelineEventRepository:      NewTimelineEventRepository(driver, log),
		CustomFieldRepository:        NewCustomFieldRepository(driver),
		CommentRepository:            NewCommentRepository(driver, neo4jDatabase),
		OpportunityRepository:        NewOpportunityRepository(driver, neo4jDatabase),
		ContractRepository:           NewContractRepository(driver, neo4jDatabase),
		ServiceLineItemRepository:    NewServiceLineItemRepository(driver, neo4jDatabase),
	}

	err := gormDb.AutoMigrate(&entity.CustomerOsIds{})
	if err != nil {
		panic(err)
	}

	return &repositories
}
