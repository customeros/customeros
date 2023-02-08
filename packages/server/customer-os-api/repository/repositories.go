package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Repositories struct {
	Drivers                       Drivers
	ActionRepository              ActionRepository
	OrganizationRepository        OrganizationRepository
	OrganizationTypeRepository    OrganizationTypeRepository
	ContactGroupRepository        ContactGroupRepository
	JobRepository                 JobRepository
	ConversationRepository        ConversationRepository
	CustomFieldTemplateRepository CustomFieldTemplateRepository
	CustomFieldRepository         CustomFieldRepository
	EntityTemplateRepository      EntityTemplateRepository
	FieldSetTemplateRepository    FieldSetTemplateRepository
	FieldSetRepository            FieldSetRepository
	UserRepository                UserRepository
	ExternalSystemRepository      ExternalSystemRepository
	NoteRepository                NoteRepository
	JobRoleRepository             JobRoleRepository
	PlaceRepository               PlaceRepository
	LocationRepository            LocationRepository
	EmailRepository               EmailRepository
	PhoneNumberRepository         PhoneNumberRepository
	TagRepository                 TagRepository
	SearchRepository              SearchRepository
	QueryRepository               QueryRepository
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
	repositories.OrganizationRepository = NewOrganizationRepository(driver)
	repositories.ContactGroupRepository = NewContactGroupRepository(driver)
	repositories.JobRepository = NewJobRepository(driver)
	repositories.ConversationRepository = NewConversationRepository(driver)
	repositories.CustomFieldTemplateRepository = NewCustomFieldTemplateRepository(driver)
	repositories.CustomFieldRepository = NewCustomFieldRepository(driver)
	repositories.EntityTemplateRepository = NewEntityTemplateRepository(driver, &repositories)
	repositories.FieldSetTemplateRepository = NewFieldSetTemplateRepository(driver, &repositories)
	repositories.FieldSetRepository = NewFieldSetRepository(driver)
	repositories.UserRepository = NewUserRepository(driver)
	repositories.ExternalSystemRepository = NewExternalSystemRepository(driver)
	repositories.NoteRepository = NewNoteRepository(driver)
	repositories.JobRoleRepository = NewJobRoleRepository(driver)
	repositories.PlaceRepository = NewPlaceRepository(driver)
	repositories.LocationRepository = NewLocationRepository(driver)
	repositories.EmailRepository = NewEmailRepository(driver)
	repositories.PhoneNumberRepository = NewPhoneNumberRepository(driver)
	repositories.OrganizationTypeRepository = NewOrganizationTypeRepository(driver)
	repositories.TagRepository = NewTagRepository(driver)
	repositories.SearchRepository = NewSearchRepository(driver)
	repositories.QueryRepository = NewQueryRepository(driver)
	return &repositories
}
