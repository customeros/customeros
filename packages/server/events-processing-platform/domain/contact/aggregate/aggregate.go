package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	ContactAggregateType eventstore.AggregateType = "contact"
)

type ContactAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Contact *models.Contact
}

func NewContactAggregateWithTenantAndID(tenant, id string) *ContactAggregate {
	contactAggregate := ContactAggregate{}
	contactAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(ContactAggregateType, tenant, id)
	contactAggregate.SetWhen(contactAggregate.When)
	contactAggregate.Contact = &models.Contact{}
	contactAggregate.Tenant = tenant
	return &contactAggregate
}

func (a *ContactAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {

	case events.ContactCreateV1:
		return a.onContactCreate(event)
	case events.ContactUpdateV1:
		return a.onContactUpdate(event)
	case events.ContactPhoneNumberLinkV1:
		return a.onPhoneNumberLink(event)
	case events.ContactEmailLinkV1:
		return a.onEmailLink(event)
	case events.ContactLocationLinkV1:
		return a.onLocationLink(event)
	case events.ContactOrganizationLinkV1:
		return a.onOrganizationLink(event)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *ContactAggregate) onContactCreate(event eventstore.Event) error {
	var eventData events.ContactCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.FirstName = eventData.FirstName
	a.Contact.LastName = eventData.LastName
	a.Contact.Prefix = eventData.Prefix
	a.Contact.Name = eventData.Name
	a.Contact.Description = eventData.Description
	a.Contact.Timezone = eventData.Timezone
	a.Contact.ProfilePhotoUrl = eventData.ProfilePhotoUrl
	a.Contact.Source = cmnmod.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Contact.CreatedAt = eventData.CreatedAt
	a.Contact.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Contact.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *ContactAggregate) onContactUpdate(event eventstore.Event) error {
	var eventData events.ContactUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if eventData.Source != a.Contact.Source.SourceOfTruth && a.Contact.Source.SourceOfTruth == constants.SourceOpenline {
		if a.Contact.Name == "" {
			a.Contact.Name = eventData.Name
		}
		if a.Contact.FirstName == "" {
			a.Contact.FirstName = eventData.FirstName
		}
		if a.Contact.LastName == "" {
			a.Contact.LastName = eventData.LastName
		}
		if a.Contact.Timezone == "" {
			a.Contact.Timezone = eventData.Timezone
		}
		if a.Contact.ProfilePhotoUrl == "" {
			a.Contact.ProfilePhotoUrl = eventData.ProfilePhotoUrl
		}
		if a.Contact.Prefix == "" {
			a.Contact.Prefix = eventData.Prefix
		}
	} else {
		a.Contact.Name = eventData.Name
		a.Contact.FirstName = eventData.FirstName
		a.Contact.LastName = eventData.LastName
		a.Contact.Prefix = eventData.Prefix
		a.Contact.Description = eventData.Description
		a.Contact.Timezone = eventData.Timezone
		a.Contact.ProfilePhotoUrl = eventData.ProfilePhotoUrl
	}
	a.Contact.UpdatedAt = eventData.UpdatedAt
	if eventData.Source == constants.SourceOpenline {
		a.Contact.Source.SourceOfTruth = eventData.Source
	}

	if eventData.ExternalSystem.Available() {
		found := false
		for _, externalSystem := range a.Contact.ExternalSystems {
			if externalSystem.ExternalSystemId == eventData.ExternalSystem.ExternalSystemId &&
				externalSystem.ExternalId == eventData.ExternalSystem.ExternalId {
				found = true
				externalSystem.ExternalUrl = eventData.ExternalSystem.ExternalUrl
				externalSystem.SyncDate = eventData.ExternalSystem.SyncDate
				externalSystem.ExternalSource = eventData.ExternalSystem.ExternalSource
				if eventData.ExternalSystem.ExternalIdSecond != "" {
					externalSystem.ExternalIdSecond = eventData.ExternalSystem.ExternalIdSecond
				}
			}
		}
		if !found {
			a.Contact.ExternalSystems = append(a.Contact.ExternalSystems, eventData.ExternalSystem)
		}
	}

	return nil
}

func (a *ContactAggregate) onPhoneNumberLink(event eventstore.Event) error {
	var eventData events.ContactLinkPhoneNumberEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.PhoneNumbers == nil {
		a.Contact.PhoneNumbers = make(map[string]models.ContactPhoneNumber)
	}
	a.Contact.PhoneNumbers[eventData.PhoneNumberId] = models.ContactPhoneNumber{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *ContactAggregate) onEmailLink(event eventstore.Event) error {
	var eventData events.ContactLinkEmailEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.Emails == nil {
		a.Contact.Emails = make(map[string]models.ContactEmail)
	}
	a.Contact.Emails[eventData.EmailId] = models.ContactEmail{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *ContactAggregate) onLocationLink(event eventstore.Event) error {
	var eventData events.ContactLinkLocationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.Locations = utils.AddToListIfNotExists(a.Contact.Locations, eventData.LocationId)
	return nil
}

func (a *ContactAggregate) onOrganizationLink(event eventstore.Event) error {
	var eventData events.ContactLinkWithOrganizationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.JobRolesByOrganization == nil {
		a.Contact.JobRolesByOrganization = make(map[string]models.JobRole)
	}
	jobRole, found := a.Contact.JobRolesByOrganization[eventData.OrganizationId]
	if !found {
		a.Contact.JobRolesByOrganization[eventData.OrganizationId] = models.JobRole{
			JobTitle:    eventData.JobTitle,
			Primary:     eventData.Primary,
			Description: eventData.Description,
			StartedAt:   eventData.StartedAt,
			EndedAt:     eventData.EndedAt,
			Source: cmnmod.Source{
				Source:        eventData.SourceFields.Source,
				SourceOfTruth: eventData.SourceFields.SourceOfTruth,
				AppSource:     eventData.SourceFields.AppSource,
			},
			CreatedAt: eventData.CreatedAt,
		}
	} else {
		if eventData.SourceFields.Source != jobRole.Source.SourceOfTruth && jobRole.Source.SourceOfTruth == constants.SourceOpenline {
			if jobRole.JobTitle == "" {
				jobRole.JobTitle = eventData.JobTitle
			}
			if jobRole.Description == "" {
				jobRole.Description = eventData.Description
			}
			if jobRole.StartedAt == nil {
				jobRole.StartedAt = eventData.StartedAt
			}
			if jobRole.EndedAt == nil {
				jobRole.EndedAt = eventData.EndedAt
			}
		} else {
			jobRole.JobTitle = eventData.JobTitle
			jobRole.Primary = eventData.Primary
			jobRole.Description = eventData.Description
			jobRole.StartedAt = eventData.StartedAt
			jobRole.EndedAt = eventData.EndedAt
		}
	}

	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}
