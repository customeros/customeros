package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapContactInputToEntity(input model.ContactInput) *entity.ContactNode {
	return &entity.ContactNode{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Label:       *input.Label,
		CompanyName: *input.CompanyName,
		ContactType: *input.ContactType,
	}
}

func MapEntityToContact(contact *entity.ContactNode) *model.Contact {
	var label = contact.Label
	var company = contact.CompanyName
	var contactType = contact.ContactType
	return &model.Contact{
		ID:          contact.Id,
		FirstName:   contact.FirstName,
		LastName:    contact.LastName,
		Label:       &label,
		CompanyName: &company,
		ContactType: &contactType,
		CreatedAt:   contact.CreatedAt,
	}
}

func MapEntitiesToContacts(contactEntities *entity.ContactNodes) []*model.Contact {
	var contacts []*model.Contact
	for _, contactEntity := range *contactEntities {
		contacts = append(contacts, MapEntityToContact(&contactEntity))
	}
	return contacts
}
