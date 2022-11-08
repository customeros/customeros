package container

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/service"
)

type ServiceContainer struct {
	ContactService         service.ContactService
	CompanyPositionService service.CompanyPositionService
	ContactGroupService    service.ContactGroupService
	TextCustomFieldService service.TextCustomFieldService
	PhoneNumberService     service.PhoneNumberService
	EmailService           service.EmailService
	UserService            service.UserService
	FieldSetService        service.FieldSetService
}

func InitServices(driver *neo4j.Driver) *ServiceContainer {
	return &ServiceContainer{
		ContactService:         service.NewContactService(driver),
		CompanyPositionService: service.NewCompanyPositionService(driver),
		ContactGroupService:    service.NewContactGroupService(driver),
		TextCustomFieldService: service.NewTextCustomFieldService(driver),
		PhoneNumberService:     service.NewPhoneNumberService(driver),
		EmailService:           service.NewEmailService(driver),
		UserService:            service.NewUserService(driver),
		FieldSetService:        service.NewFieldSetService(driver),
	}
}
