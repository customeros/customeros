package aggregate

import (
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/models"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
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
	return nil
}
