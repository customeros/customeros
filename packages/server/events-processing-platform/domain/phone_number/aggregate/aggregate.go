package aggregate

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	PhoneNumberAggregateType eventstore.AggregateType = "phone_number"
)

type PhoneNumberAggregate struct {
	*eventstore.AggregateBase
	PhoneNumber *models.PhoneNumber
}

func NewPhoneNumberAggregateWithTenantAndID(tenant, id string) *PhoneNumberAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewPhoneNumberAggregate()
	aggregate.SetID(tenant + "-" + id)
	return aggregate
}

func NewPhoneNumberAggregate() *PhoneNumberAggregate {
	phoneNumberAggregate := &PhoneNumberAggregate{PhoneNumber: models.NewPhoneNumber()}
	base := eventstore.NewAggregateBase(phoneNumberAggregate.When)
	base.SetType(PhoneNumberAggregateType)
	phoneNumberAggregate.AggregateBase = base
	return phoneNumberAggregate
}

func (a *PhoneNumberAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {

	case events.PhoneNumberCreated:
		return a.onPhoneNumberCreated(event)
	case events.PhoneNumberUpdated:
		return a.onPhoneNumberUpdated(event)

	default:
		return eventstore.ErrInvalidEventType
	}
}

func (a *PhoneNumberAggregate) onPhoneNumberCreated(event eventstore.Event) error {
	var eventData events.PhoneNumberCreatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.RawPhoneNumber = eventData.RawPhoneNumber
	a.PhoneNumber.Source = commonModels.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.PhoneNumber.CreatedAt = eventData.CreatedAt
	a.PhoneNumber.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *PhoneNumberAggregate) onPhoneNumberUpdated(event eventstore.Event) error {
	var eventData events.PhoneNumberUpdatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.Source.SourceOfTruth = eventData.SourceOfTruth
	a.PhoneNumber.UpdatedAt = eventData.UpdatedAt
	return nil
}
