package service

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type Services struct {
	ContactService             ContactService
	OrganizationService        OrganizationService
	ContactGroupService        ContactGroupService
	CustomFieldService         CustomFieldService
	PhoneNumberService         PhoneNumberService
	EmailService               EmailService
	UserService                UserService
	FieldSetService            FieldSetService
	EntityTemplateService      EntityTemplateService
	FieldSetTemplateService    FieldSetTemplateService
	CustomFieldTemplateService CustomFieldTemplateService
	ConversationService        ConversationService
	OrganizationTypeService    OrganizationTypeService
	ActionsService             ActionsService
	NoteService                NoteService
	JobRoleService             JobRoleService
	PlaceService               PlaceService
	TagService                 TagService
	SearchService              SearchService
}

func InitServices(driver *neo4j.Driver) *Services {
	repositories := repository.InitRepos(driver)

	services := Services{
		ContactService:             NewContactService(repositories),
		OrganizationService:        NewOrganizationService(repositories),
		ContactGroupService:        NewContactGroupService(repositories),
		CustomFieldService:         NewCustomFieldService(repositories),
		PhoneNumberService:         NewPhoneNumberService(repositories),
		EmailService:               NewEmailService(repositories),
		UserService:                NewUserService(repositories),
		FieldSetService:            NewFieldSetService(repositories),
		EntityTemplateService:      NewEntityTemplateService(repositories),
		FieldSetTemplateService:    NewFieldSetTemplateService(repositories),
		CustomFieldTemplateService: NewCustomFieldTemplateService(repositories),
		ConversationService:        NewConversationService(repositories),
		OrganizationTypeService:    NewOrganizationTypeService(repositories),
		ActionsService:             NewActionsService(repositories),
		NoteService:                NewNoteService(repositories),
		JobRoleService:             NewJobRoleService(repositories),
		PlaceService:               NewPlaceService(repositories),
		TagService:                 NewTagService(repositories),
	}
	services.SearchService = NewSearchService(repositories, &services)

	return &services
}
