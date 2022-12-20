package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapContactInputToEntity(input model.ContactInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		CreatedAt: input.CreatedAt,
		FirstName: utils.IfNotNilString(input.FirstName),
		LastName:  utils.IfNotNilString(input.LastName),
		Readonly:  utils.IfNotNilBool(input.Readonly),
		Label:     utils.IfNotNilString(input.Label),
		Title:     utils.IfNotNilString(input.Title, func() string { return input.Title.String() }),
	}
	return &contactEntity
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
