package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

func MapContactUpdateInputToContactFields(input model.ContactUpdateInput) data_fields.ContactFields {
	fields := data_fields.ContactFields{
		Name:            input.Name,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Prefix:          input.Prefix,
		Description:     input.Description,
		Timezone:        input.Timezone,
		ProfilePhotoUrl: input.ProfilePhotoURL,
		Username:        input.Username,
		Source:          utils.StringPtr(neo4jentity.DataSourceOpenline.String()),
		AppSource:       utils.StringPtr(constants.AppSourceCustomerOsApi),
	}
	return fields
}

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
		AppSource:       utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &contactEntity
}

func MapCustomerContactInputToEntity(input model.CustomerContactInput) *neo4jentity.ContactEntity {
	contactEntity := neo4jentity.ContactEntity{
		CreatedAt:   utils.IfNotNilTimeWithDefault(input.CreatedAt, utils.Now()),
		Name:        utils.IfNotNilString(input.Name),
		FirstName:   utils.IfNotNilString(input.FirstName),
		LastName:    utils.IfNotNilString(input.LastName),
		Prefix:      utils.IfNotNilString(input.Prefix),
		Description: utils.IfNotNilString(input.Description),
		Timezone:    utils.IfNotNilString(input.Timezone),
		Source:      neo4jentity.DataSourceOpenline,
		AppSource:   utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &contactEntity
}

func MapEntityToContact(contact *neo4jentity.ContactEntity) *model.Contact {
	return &model.Contact{
		Metadata: &model.Metadata{
			ID:          contact.Id,
			Created:     contact.CreatedAt,
			LastUpdated: contact.UpdatedAt,
			Source:      MapDataSourceToModel(contact.Source),
			AppSource:   contact.AppSource,
			Version:     contact.AggregateVersion,
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
		AppSource:       utils.StringPtr(contact.AppSource),
		EnrichDetails:   prepareContactEnrichDetails(contact.EnrichDetails),
	}
}

func MapEntitiesToContacts(contactEntities *neo4jentity.ContactEntities) []*model.Contact {
	var contacts []*model.Contact
	for _, contactEntity := range *contactEntities {
		contacts = append(contacts, MapEntityToContact(&contactEntity))
	}
	return contacts
}

func prepareContactEnrichDetails(enrichDetails neo4jentity.ContactEnrichDetails) *model.EnrichDetails {
	output := model.EnrichDetails{
		RequestedAt: enrichDetails.EnrichRequestedAt,
		EnrichedAt:  enrichDetails.EnrichedAt,
		FailedAt:    enrichDetails.EnrichFailedAt,

		EmailRequestedAt: enrichDetails.FindWorkEmailWithBetterContactRequestedAt,
		EmailEnrichedAt:  enrichDetails.FindWorkEmailWithBetterContactCompletedAt,
		EmailFound:       enrichDetails.FindWorkEmailWithBetterContactFound,

		MobilePhoneFound:       enrichDetails.FindMobilePhoneWithBetterContactFound,
		MobilePhoneRequestedAt: enrichDetails.FindMobilePhoneWithBetterContactRequestedAt,
		MobilePhoneEnrichedAt:  enrichDetails.FindMobilePhoneWithBetterContactCompletedAt,
	}
	if enrichDetails.EnrichedAt == nil && enrichDetails.EnrichFailedAt == nil && enrichDetails.EnrichRequestedAt != nil {
		// if requested is older than 1 min, remove it
		if time.Since(*enrichDetails.EnrichRequestedAt) > time.Minute {
			output.RequestedAt = nil
		}
	}
	if output.EmailEnrichedAt == nil && output.EmailFound == nil && output.EmailRequestedAt != nil {
		// if requested is older than 10 min, set found to false
		if time.Since(*output.EmailRequestedAt) > 10*time.Minute {
			output.EmailRequestedAt = nil
			output.EmailFound = utils.BoolPtr(false)
		}
	}
	if output.MobilePhoneEnrichedAt == nil && output.MobilePhoneFound == nil && output.MobilePhoneRequestedAt != nil {
		// if requested is older than 10 min, set found to false
		if time.Since(*output.MobilePhoneRequestedAt) > 10*time.Minute {
			output.MobilePhoneRequestedAt = nil
			output.MobilePhoneFound = utils.BoolPtr(false)
		}
	}
	return &output
}
