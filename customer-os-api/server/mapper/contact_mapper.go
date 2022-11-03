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
	if input.Title != nil {
		contactEntity.Title = input.Title.String()
	}
	if input.Notes != nil {
		contactEntity.Notes = *input.Notes
	}
	if input.ContactType != nil {
		contactEntity.ContactType = *input.ContactType
	}
	return &contactEntity
}

func MapEntityToContact(contact *entity.ContactEntity) *model.Contact {
	var title = model.PersonTitle(contact.Title)
	if !title.IsValid() {
		title = ""
	}
	var label = contact.Label
	var notes = contact.Notes
	var contactType = contact.ContactType
	return &model.Contact{
		ID:          contact.Id,
		Title:       &title,
		FirstName:   contact.FirstName,
		LastName:    contact.LastName,
		Label:       &label,
		Notes:       &notes,
		ContactType: &contactType,
		CreatedAt:   contact.CreatedAt,
	}
}

func MapEntitiesToContacts(contactEntities *entity.ContactEntities) []*model.Contact {
	var contacts []*model.Contact
	for _, contactEntity := range *contactEntities {
		contacts = append(contacts, MapEntityToContact(&contactEntity))
	}
	return contacts
}
