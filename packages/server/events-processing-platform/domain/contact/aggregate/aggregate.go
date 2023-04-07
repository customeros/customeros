package aggregate

import (
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
	*eventstore.AggregateBase
	Contact *models.Contact
}

func NewContactAggregateWithTenantAndID(tenant, id string) *ContactAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewContactAggregate()
	aggregate.SetID(tenant + "-" + id)
	return aggregate
}

func NewContactAggregate() *ContactAggregate {
	contactAggregate := &ContactAggregate{Contact: models.NewContact()}
	base := eventstore.NewAggregateBase(contactAggregate.When)
	base.SetType(ContactAggregateType)
	contactAggregate.AggregateBase = base
	return contactAggregate
}

func (contactAggregate *ContactAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {

	case events.ContactCreated:
		return contactAggregate.onContactCreated(event)
	case events.ContactUpdated:
		return contactAggregate.onContactUpdated(event)
	case events.ContactPhoneNumberLinked:
		return contactAggregate.onPhoneNumberLinked(event)
	case events.ContactEmailLinked:
		return contactAggregate.onEmailLinked(event)

	default:
		return eventstore.ErrInvalidEventType
	}
}

func (a *ContactAggregate) onContactCreated(event eventstore.Event) error {
	var eventData events.ContactCreatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.FirstName = eventData.FirstName
	a.Contact.LastName = eventData.LastName
	a.Contact.Name = eventData.Name
	a.Contact.Prefix = eventData.Prefix
	a.Contact.Source = commonModels.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Contact.CreatedAt = eventData.CreatedAt
	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *ContactAggregate) onContactUpdated(event eventstore.Event) error {
	var eventData events.ContactUpdatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.Source.SourceOfTruth = eventData.SourceOfTruth
	a.Contact.UpdatedAt = eventData.UpdatedAt
	a.Contact.FirstName = eventData.FirstName
	a.Contact.LastName = eventData.LastName
	a.Contact.Name = eventData.Name
	a.Contact.Prefix = eventData.Prefix
	return nil
}

func (a *ContactAggregate) onPhoneNumberLinked(event eventstore.Event) error {
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

func (a *ContactAggregate) onEmailLinked(event eventstore.Event) error {
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
