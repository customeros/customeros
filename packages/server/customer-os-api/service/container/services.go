package container

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type Services struct {
	ContactService               service.ContactService
	CompanyService               service.CompanyService
	ContactGroupService          service.ContactGroupService
	CustomFieldService           service.CustomFieldService
	PhoneNumberService           service.PhoneNumberService
	EmailService                 service.EmailService
	UserService                  service.UserService
	FieldSetService              service.FieldSetService
	EntityDefinitionService      service.EntityDefinitionService
	FieldSetDefinitionService    service.FieldSetDefinitionService
	CustomFieldDefinitionService service.CustomFieldDefinitionService
	ConversationService          service.ConversationService
	ContactTypeService           service.ContactTypeService
	ActionsService               service.ActionsService
	NoteService                  service.NoteService
	ContactRoleService           service.ContactRoleService
	AddressService               service.AddressService
}

func InitServices(driver *neo4j.Driver) *Services {
	repositories := repository.InitRepos(driver)

	return &Services{
		ContactService:               service.NewContactService(repositories),
		CompanyService:               service.NewCompanyService(repositories),
		ContactGroupService:          service.NewContactGroupService(repositories),
		CustomFieldService:           service.NewCustomFieldService(repositories),
		PhoneNumberService:           service.NewPhoneNumberService(repositories),
		EmailService:                 service.NewEmailService(repositories),
		UserService:                  service.NewUserService(repositories),
		FieldSetService:              service.NewFieldSetService(repositories),
		EntityDefinitionService:      service.NewEntityDefinitionService(repositories),
		FieldSetDefinitionService:    service.NewFieldSetDefinitionService(repositories),
		CustomFieldDefinitionService: service.NewCustomFieldDefinitionService(repositories),
		ConversationService:          service.NewConversationService(repositories),
		ContactTypeService:           service.NewContactTypeService(repositories),
		ActionsService:               service.NewActionsService(repositories),
		NoteService:                  service.NewNoteService(repositories),
		ContactRoleService:           service.NewContactRoleService(repositories),
		AddressService:               service.NewAddressService(repositories),
	}
}
