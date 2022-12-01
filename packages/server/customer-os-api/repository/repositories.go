package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Repositories struct {
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

func InitRepos(driver *neo4j.Driver) *Repositories {
	repositories := Repositories{
		Drivers: Drivers{
			Neo4jDriver: driver,
		},
	}
	repositories.CompanyRepository = NewCompanyRepository(driver, &repositories)
	repositories.ContactGroupRepository = NewContactGroupRepository(driver, &repositories)
	repositories.ContactRepository = NewContactRepository(driver, &repositories)
	repositories.ContactTypeRepository = NewContactTypeRepository(driver, &repositories)
	repositories.ConversationRepository = NewConversationRepository(driver, &repositories)
	repositories.CustomFieldDefinitionRepository = NewCustomFieldDefinitionRepository(driver, &repositories)
	repositories.CustomFieldRepository = NewCustomFieldRepository(driver, &repositories)
	repositories.EntityDefinitionRepository = NewEntityDefinitionRepository(driver, &repositories)
	repositories.FieldSetDefinitionRepository = NewFieldSetDefinitionRepository(driver, &repositories)
	repositories.FieldSetRepository = NewFieldSetRepository(driver, &repositories)
	repositories.UserRepository = NewUserRepository(driver, &repositories)
	return &repositories
}
