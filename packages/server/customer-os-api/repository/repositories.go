package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Repositories struct {
	Drivers                       Drivers
	ActionRepository              ActionRepository
	CompanyRepository             CompanyRepository
	ContactGroupRepository        ContactGroupRepository
	ContactRepository             ContactRepository
	ContactTypeRepository         ContactTypeRepository
	ConversationRepository        ConversationRepository
	MessageRepository             MessageRepository
	CustomFieldTemplateRepository CustomFieldTemplateRepository
	CustomFieldRepository         CustomFieldRepository
	EntityTemplateRepository      EntityTemplateRepository
	FieldSetTemplateRepository    FieldSetTemplateRepository
	FieldSetRepository            FieldSetRepository
	UserRepository                UserRepository
	ExternalSystemRepository      ExternalSystemRepository
	NoteRepository                NoteRepository
	ContactRoleRepository         ContactRoleRepository
	AddressRepository             AddressRepository
	EmailRepository               EmailRepository
	PhoneNumberRepository         PhoneNumberRepository
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
	repositories.ActionRepository = NewActionRepository(driver)
	repositories.CompanyRepository = NewCompanyRepository(driver)
	repositories.ContactGroupRepository = NewContactGroupRepository(driver)
	repositories.ContactRepository = NewContactRepository(driver)
	repositories.ContactTypeRepository = NewContactTypeRepository(driver)
	repositories.ConversationRepository = NewConversationRepository(driver)
	repositories.MessageRepository = NewMessageRepository(driver)
	repositories.CustomFieldTemplateRepository = NewCustomFieldTemplateRepository(driver)
	repositories.CustomFieldRepository = NewCustomFieldRepository(driver)
	repositories.EntityTemplateRepository = NewEntityTemplateRepository(driver, &repositories)
	repositories.FieldSetTemplateRepository = NewFieldSetTemplateRepository(driver, &repositories)
	repositories.FieldSetRepository = NewFieldSetRepository(driver)
	repositories.UserRepository = NewUserRepository(driver)
	repositories.ExternalSystemRepository = NewExternalSystemRepository(driver)
	repositories.NoteRepository = NewNoteRepository(driver)
	repositories.ContactRoleRepository = NewContactRoleRepository(driver)
	repositories.AddressRepository = NewAddressRepository(driver)
	repositories.EmailRepository = NewEmailRepository(driver)
	repositories.PhoneNumberRepository = NewPhoneNumberRepository(driver)
	return &repositories
}
