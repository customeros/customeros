package container

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type Services struct {
	ContactService service.ContactService
}

/*func InitServices(driver *neo4j.Driver) *Services {
	repositories := repository.InitRepos(driver)

	return &Services{
		ContactCommandsService:             service.NewContactCommandsService(repositories),
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
		JobRoleService:             service.NewJobRoleService(repositories),
		PlaceService:               service.NewPlaceService(repositories),
		TagService:                 service.NewTagService(repositories),
	}
}
*/
