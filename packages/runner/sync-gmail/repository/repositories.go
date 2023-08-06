package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Neo4jDriver *neo4j.DriverWithContext

	CommonRepositories *commonRepository.Repositories

	//pg repositories
	RawEmailRepository              RawEmailRepository
	PersonalEmailProviderRepository PersonalEmailProviderRepository

	//neo4j repositories
	TenantRepository           TenantRepository
	UserRepository             UserRepository
	EmailRepository            EmailRepository
	ExternalSystemRepository   ExternalSystemRepository
	InteractionEventRepository InteractionEventRepository
	ContactRepository          ContactRepository
	OrganizationRepository     OrganizationRepository
	WorkspaceRepository        WorkspaceRepository
	AnalysisRepository         AnalysisRepository
	ActionItemRepository       ActionItemRepository
	DomainRepository           DomainRepository
}

func InitRepos(driver *neo4j.DriverWithContext, gormDb *gorm.DB) *Repositories {
	repositories := Repositories{
		Neo4jDriver: driver,

		CommonRepositories: commonRepository.InitRepositories(gormDb, driver),

		RawEmailRepository:              NewRawEmailRepository(gormDb),
		PersonalEmailProviderRepository: NewPersonalEmailProviderRepository(gormDb),

		TenantRepository:           NewTenantRepository(driver),
		UserRepository:             NewUserRepository(driver),
		EmailRepository:            NewEmailRepository(driver),
		ExternalSystemRepository:   NewExternalSystemRepository(driver),
		InteractionEventRepository: NewInteractionEventRepository(driver),
		ContactRepository:          NewContactRepository(driver),
		OrganizationRepository:     NewOrganizationRepository(driver),
		WorkspaceRepository:        NewWorkspaceRepository(driver),
		AnalysisRepository:         NewAnalysisRepository(driver),
		ActionItemRepository:       NewActionItemRepository(driver),
		DomainRepository:           NewDomainRepository(driver),
	}

	var err error

	err = gormDb.AutoMigrate(&entity.PersonalEmailProvider{})
	if err != nil {
		panic(err)
	}

	return &repositories
}
