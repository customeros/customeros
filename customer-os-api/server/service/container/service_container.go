package container

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/service"
)

type ServiceContainer struct {
	ContactService               service.ContactService
	CompanyPositionService       service.CompanyService
	ContactGroupService          service.ContactGroupService
	TextCustomFieldService       service.TextCustomFieldService
	PhoneNumberService           service.PhoneNumberService
	EmailService                 service.EmailService
	UserService                  service.UserService
	FieldSetService              service.FieldSetService
	EntityDefinitionService      service.EntityDefinitionService
	FieldSetDefinitionService    service.FieldSetDefinitionService
	CustomFieldDefinitionService service.CustomFieldDefinitionService
}

func InitServices(driver *neo4j.Driver) *ServiceContainer {
	repoContainer := repository.InitRepos(driver)

	return &ServiceContainer{
		ContactService:               service.NewContactService(repoContainer),
		CompanyPositionService:       service.NewCompanyPositionService(repoContainer),
		ContactGroupService:          service.NewContactGroupService(repoContainer),
		TextCustomFieldService:       service.NewTextCustomFieldService(repoContainer),
		PhoneNumberService:           service.NewPhoneNumberService(repoContainer),
		EmailService:                 service.NewEmailService(repoContainer),
		UserService:                  service.NewUserService(repoContainer),
		FieldSetService:              service.NewFieldSetService(repoContainer),
		EntityDefinitionService:      service.NewEntityDefinitionService(repoContainer),
		FieldSetDefinitionService:    service.NewFieldSetDefinitionService(repoContainer),
		CustomFieldDefinitionService: service.NewCustomFieldDefinitionService(repoContainer),
	}
}
