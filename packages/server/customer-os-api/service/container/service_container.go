package container

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type ServiceContainer struct {
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
}

func InitServices(driver *neo4j.Driver) *ServiceContainer {
	repoContainer := repository.InitRepos(driver)

	return &ServiceContainer{
		ContactService:               service.NewContactService(repoContainer),
		CompanyService:               service.NewCompanyService(repoContainer),
		ContactGroupService:          service.NewContactGroupService(repoContainer),
		CustomFieldService:           service.NewCustomFieldService(repoContainer),
		PhoneNumberService:           service.NewPhoneNumberService(repoContainer),
		EmailService:                 service.NewEmailService(repoContainer),
		UserService:                  service.NewUserService(repoContainer),
		FieldSetService:              service.NewFieldSetService(repoContainer),
		EntityDefinitionService:      service.NewEntityDefinitionService(repoContainer),
		FieldSetDefinitionService:    service.NewFieldSetDefinitionService(repoContainer),
		CustomFieldDefinitionService: service.NewCustomFieldDefinitionService(repoContainer),
		ConversationService:          service.NewConversationService(repoContainer),
		ContactTypeService:           service.NewContactTypeService(repoContainer),
	}
}
