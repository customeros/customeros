package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapContactInputToEntity(input model.ContactInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	if input.Label != nil {
		contactEntity.Label = *input.Label
	}
	if input.CompanyName != nil {
		contactEntity.CompanyName = *input.CompanyName
	}
	if input.ContactType != nil {
		contactEntity.ContactType = *input.ContactType
	}
	if input.TextCustomFields != nil {
		contactEntity.TextCustomFields = *MapTextCustomFieldInputsToEntities(input.TextCustomFields)
	}
	return &contactEntity
}

func MapEntityToContact(contact *entity.ContactEntity) *model.Contact {
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
