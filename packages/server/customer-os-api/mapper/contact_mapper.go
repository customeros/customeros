package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapContactInputToEntity(input model.ContactInput) *entity.ContactEntity {
	contactEntity := new(entity.ContactEntity)
	contactEntity.CreatedAt = input.CreatedAt
	if input.FirstName != nil {
		contactEntity.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		contactEntity.LastName = *input.LastName
	}
	if input.Readonly != nil {
		contactEntity.Readonly = *input.Readonly
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
	return contactEntity
}

func MapContactUpdateInputToEntity(input model.ContactUpdateInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		Id: input.ID,
	}
	if input.FirstName != nil {
		contactEntity.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		contactEntity.LastName = *input.LastName
	}
	if input.Readonly != nil {
		contactEntity.Readonly = *input.Readonly
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
	return &model.Contact{
		ID:        contact.Id,
		Title:     &title,
		FirstName: utils.StringPtr(contact.FirstName),
		LastName:  utils.StringPtr(contact.LastName),
		Label:     utils.StringPtr(contact.Label),
		Notes:     utils.StringPtr(contact.Notes),
		Readonly:  contact.Readonly,
		CreatedAt: *contact.CreatedAt,
	}
}

func MapEntitiesToContacts(contactEntities *entity.ContactEntities) []*model.Contact {
	var contacts []*model.Contact
	for _, contactEntity := range *contactEntities {
		contacts = append(contacts, MapEntityToContact(&contactEntity))
	}
	return contacts
}
