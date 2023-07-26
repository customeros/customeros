package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
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
	a.Contact.Source = commonModels.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Contact.CreatedAt = eventData.CreatedAt
	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *ContactAggregate) onContactUpdate(event eventstore.Event) error {
	var eventData events.ContactUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.Source.SourceOfTruth = eventData.SourceOfTruth
	a.Contact.UpdatedAt = eventData.UpdatedAt
	a.Contact.FirstName = eventData.FirstName
	a.Contact.LastName = eventData.LastName
	a.Contact.Prefix = eventData.Prefix
	a.Contact.Description = eventData.Description
	a.Contact.Timezone = eventData.Timezone
	a.Contact.Name = eventData.Name
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
