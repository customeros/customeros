package mapper

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapContactInputToEntity(input model.ContactInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		CreatedAt:       input.CreatedAt,
		FirstName:       utils.IfNotNilString(input.FirstName),
		LastName:        utils.IfNotNilString(input.LastName),
		Name:            utils.IfNotNilString(input.Name),
		Prefix:          utils.IfNotNilString(input.Prefix),
		Description:     utils.IfNotNilString(input.Description),
		Timezone:        utils.IfNotNilString(input.Timezone),
		ProfilePhotoUrl: utils.IfNotNilString(input.ProfilePhotoURL),
		Source:          neo4jentity.DataSourceOpenline,
		SourceOfTruth:   neo4jentity.DataSourceOpenline,
		AppSource:       utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &contactEntity
}

func MapCustomerContactInputToEntity(input model.CustomerContactInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		CreatedAt:     input.CreatedAt,
		Name:          utils.IfNotNilString(input.Name),
		FirstName:     utils.IfNotNilString(input.FirstName),
		LastName:      utils.IfNotNilString(input.LastName),
		Prefix:        utils.IfNotNilString(input.Prefix),
		Description:   utils.IfNotNilString(input.Description),
		Timezone:      utils.IfNotNilString(input.Timezone),
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &contactEntity
}

func MapContactUpdateInputToEntity(input model.ContactUpdateInput) *entity.ContactEntity {
	contactEntity := entity.ContactEntity{
		Id:              input.ID,
		Name:            utils.IfNotNilString(input.Name),
		FirstName:       utils.IfNotNilString(input.FirstName),
		LastName:        utils.IfNotNilString(input.LastName),
		Prefix:          utils.IfNotNilString(input.Prefix),
		Description:     utils.IfNotNilString(input.Description),
		Timezone:        utils.IfNotNilString(input.Timezone),
		ProfilePhotoUrl: utils.IfNotNilString(input.ProfilePhotoURL),
		SourceOfTruth:   neo4jentity.DataSourceOpenline,
	}
	return &contactEntity
}

func MapEntityToContact(contact *entity.ContactEntity) *model.Contact {
	return &model.Contact{
		ID:              contact.Id,
		Prefix:          utils.StringPtr(contact.Prefix),
		Name:            utils.StringPtr(contact.Name),
		FirstName:       utils.StringPtr(contact.FirstName),
		LastName:        utils.StringPtr(contact.LastName),
		Description:     utils.StringPtr(contact.Description),
		Timezone:        utils.StringPtr(contact.Timezone),
		ProfilePhotoURL: utils.StringPtr(contact.ProfilePhotoUrl),
		CreatedAt:       *contact.CreatedAt,
		UpdatedAt:       contact.UpdatedAt,
		Source:          MapDataSourceToModel(contact.Source),
		SourceOfTruth:   MapDataSourceToModel(contact.SourceOfTruth),
		AppSource:       utils.StringPtr(contact.AppSource),
	}
}

func MapEntitiesToContacts(contactEntities *entity.ContactEntities) []*model.Contact {
	var contacts []*model.Contact
	for _, contactEntity := range *contactEntities {
		contacts = append(contacts, MapEntityToContact(&contactEntity))
	}
	return contacts
}
