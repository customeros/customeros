package service

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ServiceContainer struct {
	ContactService                             ContactService
	ContactGroupService                        ContactGroupService
	ContactWithContactGroupRelationshipService ContactWithContactGroupRelationshipService
	TextCustomPropertyService                  TextCustomPropertyService
}

func InitServices(driver *neo4j.Driver) *ServiceContainer {
	return &ServiceContainer{
		ContactService:      NewContactService(driver),
		ContactGroupService: NewContactGroupService(driver),
		ContactWithContactGroupRelationshipService: NewContactWithContactGroupRelationshipService(driver),
		TextCustomPropertyService:                  NewTextCustomPropertyService(driver),
	}
}
