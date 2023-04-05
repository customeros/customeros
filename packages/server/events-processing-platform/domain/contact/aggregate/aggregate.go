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
