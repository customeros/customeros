package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapContactInputToEntity(input model.ContactInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		CreatedAt: input.CreatedAt,
		FirstName: utils.IfNotNilString(input.FirstName),
		LastName:  utils.IfNotNilString(input.LastName),
		Title:     utils.IfNotNilString(input.Title, func() string { return input.Title.String() }),
	}
	return &contactEntity
}

func MapContactUpdateInputToEntity(input model.ContactUpdateInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		Id:        input.ID,
		FirstName: utils.IfNotNilString(input.FirstName),
		LastName:  utils.IfNotNilString(input.LastName),
		Title:     utils.IfNotNilString(input.Title, func() string { return input.Title.String() }),
	}
	return &contactEntity
}

func MapEntityToContact(contact *entity.ContactEntity) *model.Contact {
	var title = model.PersonTitle(contact.Title)
	return &model.Contact{
		ID:            contact.Id,
		Title:         &title,
		Name:          utils.StringPtr(contact.Name),
		FirstName:     utils.StringPtr(contact.FirstName),
		LastName:      utils.StringPtr(contact.LastName),
		CreatedAt:     *contact.CreatedAt,
		UpdatedAt:     contact.UpdatedAt,
		Source:        MapDataSourceToModel(contact.Source),
		SourceOfTruth: MapDataSourceToModel(contact.Source),
		AppSource:     utils.StringPtr(contact.AppSource),
	}
}

func MapEntitiesToContacts(contactEntities *entity.ContactEntities) []*model.Contact {
	var contacts []*model.Contact
	for _, contactEntity := range *contactEntities {
		contacts = append(contacts, MapEntityToContact(&contactEntity))
	}
	return contacts
}
