package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
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
	return &contactEntity
}

func MapContactUpdateInputToEntity(input model.ContactUpdateInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		Id:        input.ID,
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
	return &contactEntity
}

func MapEntityToContact(contact *entity.ContactEntity) *model.Contact {
	var title = model.PersonTitle(contact.Title)
	if !title.IsValid() {
		title = model.PersonTitleMr
	}
	return &model.Contact{
		ID:        contact.Id,
		Title:     &title,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Label:     utils.StringPtr(contact.Label),
		Notes:     utils.StringPtr(contact.Notes),
		CreatedAt: contact.CreatedAt,
	}
}

func MapEntitiesToContacts(contactEntities *entity.ContactEntities) []*model.Contact {
	var contacts []*model.Contact
	for _, contactEntity := range *contactEntities {
		contacts = append(contacts, MapEntityToContact(&contactEntity))
	}
	return contacts
}
