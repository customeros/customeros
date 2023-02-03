package container

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type Services struct {
	ContactService             service.ContactService
	OrganizationService        service.OrganizationService
	ContactGroupService        service.ContactGroupService
	CustomFieldService         service.CustomFieldService
	PhoneNumberService         service.PhoneNumberService
	EmailService               service.EmailService
	UserService                service.UserService
	FieldSetService            service.FieldSetService
	EntityTemplateService      service.EntityTemplateService
	FieldSetTemplateService    service.FieldSetTemplateService
	CustomFieldTemplateService service.CustomFieldTemplateService
	ConversationService        service.ConversationService
	ContactTypeService         service.ContactTypeService
	OrganizationTypeService    service.OrganizationTypeService
	ActionsService             service.ActionsService
	NoteService                service.NoteService
	ContactRoleService         service.ContactRoleService
	PlaceService               service.PlaceService
	TagService                 service.TagService
}

func InitServices(driver *neo4j.Driver) *Services {
	repositories := repository.InitRepos(driver)

	return &Services{
		ContactService:             service.NewContactService(repositories),
		OrganizationService:        service.NewOrganizationService(repositories),
		ContactGroupService:        service.NewContactGroupService(repositories),
		CustomFieldService:         service.NewCustomFieldService(repositories),
		PhoneNumberService:         service.NewPhoneNumberService(repositories),
		EmailService:               service.NewEmailService(repositories),
		UserService:                service.NewUserService(repositories),
		FieldSetService:            service.NewFieldSetService(repositories),
		EntityTemplateService:      service.NewEntityTemplateService(repositories),
		FieldSetTemplateService:    service.NewFieldSetTemplateService(repositories),
		CustomFieldTemplateService: service.NewCustomFieldTemplateService(repositories),
		ConversationService:        service.NewConversationService(repositories),
		ContactTypeService:         service.NewContactTypeService(repositories),
		OrganizationTypeService:    service.NewOrganizationTypeService(repositories),
		ActionsService:             service.NewActionsService(repositories),
		NoteService:                service.NewNoteService(repositories),
		ContactRoleService:         service.NewContactRoleService(repositories),
		PlaceService:               service.NewPlaceService(repositories),
		TagService:                 service.NewTagService(repositories),
	}
}
