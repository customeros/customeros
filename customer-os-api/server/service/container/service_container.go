package container

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/service"
)

type ServiceContainer struct {
	ContactService                             service.ContactService
	ContactGroupService                        service.ContactGroupService
	ContactWithContactGroupRelationshipService service.ContactWithContactGroupRelationshipService
	TextCustomFieldService                     service.TextCustomFieldService
	PhoneNumberService                         service.PhoneNumberService
	EmailService                               service.EmailService
}

func InitServices(driver *neo4j.Driver) *ServiceContainer {
	return &ServiceContainer{
		ContactService:      service.NewContactService(driver),
		ContactGroupService: service.NewContactGroupService(driver),
		ContactWithContactGroupRelationshipService: service.NewContactWithContactGroupRelationshipService(driver),
		TextCustomFieldService:                     service.NewTextCustomFieldService(driver),
		PhoneNumberService:                         service.NewPhoneNumberService(driver),
		EmailService:                               service.NewEmailService(driver),
	}
}
