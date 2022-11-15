package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type RepositoryContainer struct {
	Drivers                         Drivers
	ContactRepository               ContactRepository
	FieldSetRepository              FieldSetRepository
	EntityDefinitionRepository      EntityDefinitionRepository
	FieldSetDefinitionRepository    FieldSetDefinitionRepository
	CustomFieldDefinitionRepository CustomFieldDefinitionRepository
	CustomFieldRepository           CustomFieldRepository
	ConversationRepository          ConversationRepository
	ContactTypeRepository           ContactTypeRepository
	CompanyRepository               CompanyRepository
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
	container.EntityDefinitionRepository = NewEntityDefinitionRepository(driver, &container)
	container.FieldSetDefinitionRepository = NewFieldSetDefinitionRepository(driver, &container)
	container.CustomFieldDefinitionRepository = NewCustomFieldDefinitionRepository(driver, &container)
	container.ConversationRepository = NewConversationRepository(driver, &container)
	container.ContactRepository = NewContactRepository(driver, &container)
	container.FieldSetRepository = NewFieldSetRepository(driver, &container)
	container.CustomFieldRepository = NewCustomFieldRepository(driver, &container)
	container.ContactTypeRepository = NewContactTypeRepository(driver, &container)
	container.CompanyRepository = NewCompanyRepository(driver, &container)
	return &container
}
