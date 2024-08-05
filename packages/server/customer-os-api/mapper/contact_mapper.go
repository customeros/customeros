package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapContactInputToEntity(input model.ContactInput) *neo4jentity.ContactEntity {
	contactEntity := neo4jentity.ContactEntity{
		CreatedAt:       utils.IfNotNilTimeWithDefault(input.CreatedAt, utils.Now()),
		FirstName:       utils.IfNotNilString(input.FirstName),
		LastName:        utils.IfNotNilString(input.LastName),
		Name:            utils.IfNotNilString(input.Name),
		Prefix:          utils.IfNotNilString(input.Prefix),
		Description:     utils.IfNotNilString(input.Description),
		Timezone:        utils.IfNotNilString(input.Timezone),
		ProfilePhotoUrl: utils.IfNotNilString(input.ProfilePhotoURL),
		Username:        utils.IfNotNilString(input.Username),
		Source:          neo4jentity.DataSourceOpenline,
		SourceOfTruth:   neo4jentity.DataSourceOpenline,
		AppSource:       utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &contactEntity
}

func MapCustomerContactInputToEntity(input model.CustomerContactInput) *neo4jentity.ContactEntity {
	contactEntity := neo4jentity.ContactEntity{
		CreatedAt:     utils.IfNotNilTimeWithDefault(input.CreatedAt, utils.Now()),
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

func MapEntityToContact(contact *neo4jentity.ContactEntity) *model.Contact {
	return &model.Contact{
		Metadata: &model.Metadata{
			ID:            contact.Id,
			Created:       contact.CreatedAt,
			LastUpdated:   contact.UpdatedAt,
			Source:        MapDataSourceToModel(contact.Source),
			SourceOfTruth: MapDataSourceToModel(contact.SourceOfTruth),
			AppSource:     contact.AppSource,
			Version:       contact.AggregateVersion,
		},
		ID:              contact.Id,
		Prefix:          utils.StringPtr(contact.Prefix),
		Name:            utils.StringPtr(contact.Name),
		FirstName:       utils.StringPtr(contact.FirstName),
		LastName:        utils.StringPtr(contact.LastName),
		Description:     utils.StringPtr(contact.Description),
		Timezone:        utils.StringPtr(contact.Timezone),
		ProfilePhotoURL: utils.StringPtr(contact.ProfilePhotoUrl),
		Username:        utils.StringPtr(contact.Username),
		Hide:            utils.BoolPtr(contact.Hide),
		CreatedAt:       contact.CreatedAt,
		UpdatedAt:       contact.UpdatedAt,
		Source:          MapDataSourceToModel(contact.Source),
		SourceOfTruth:   MapDataSourceToModel(contact.SourceOfTruth),
		AppSource:       utils.StringPtr(contact.AppSource),
	}
}

func MapEntitiesToContacts(contactEntities *neo4jentity.ContactEntities) []*model.Contact {
	var contacts []*model.Contact
	for _, contactEntity := range *contactEntities {
		contacts = append(contacts, MapEntityToContact(&contactEntity))
	}
	return contacts
}
