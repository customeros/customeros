package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type RepositoryContainer struct {
	Drivers                         Drivers
	CompanyRepository               CompanyRepository
	ContactGroupRepository          ContactGroupRepository
	ContactRepository               ContactRepository
	ContactTypeRepository           ContactTypeRepository
	ConversationRepository          ConversationRepository
	CustomFieldDefinitionRepository CustomFieldDefinitionRepository
	CustomFieldRepository           CustomFieldRepository
	EntityDefinitionRepository      EntityDefinitionRepository
	FieldSetDefinitionRepository    FieldSetDefinitionRepository
	FieldSetRepository              FieldSetRepository
	UserRepository                  UserRepository
}

type Drivers struct {
	Neo4jDriver *neo4j.Driver
}

func InitRepos(driver *neo4j.Driver) *RepositoryContainer {
	container := RepositoryContainer{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
	}
	container.CompanyRepository = NewCompanyRepository(driver, &container)
	container.ContactGroupRepository = NewContactGroupRepository(driver, &container)
	container.ContactRepository = NewContactRepository(driver, &container)
	container.ContactTypeRepository = NewContactTypeRepository(driver, &container)
	container.ConversationRepository = NewConversationRepository(driver, &container)
	container.CustomFieldDefinitionRepository = NewCustomFieldDefinitionRepository(driver, &container)
	container.CustomFieldRepository = NewCustomFieldRepository(driver, &container)
	container.EntityDefinitionRepository = NewEntityDefinitionRepository(driver, &container)
	container.FieldSetDefinitionRepository = NewFieldSetDefinitionRepository(driver, &container)
	container.FieldSetRepository = NewFieldSetRepository(driver, &container)
	container.UserRepository = NewUserRepository(driver, &container)
	return &container
}
