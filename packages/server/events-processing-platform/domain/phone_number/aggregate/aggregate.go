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

	case events.PhoneNumberCreateV1:
		return a.onPhoneNumberCreate(event)
	case events.PhoneNumberUpdateV1:
		return a.onPhoneNumberUpdate(event)
	case events.PhoneNumberValidationSkippedV1:
		return a.OnPhoneNumberSkippedValidation(event)
	case events.PhoneNumberValidationFailedV1:
		return a.OnPhoneNumberFailedValidation(event)
	case events.PhoneNumberValidatedV1:
		return a.OnPhoneNumberValidated(event)

	default:
		return eventstore.ErrInvalidEventType
	}
}

func (a *PhoneNumberAggregate) onPhoneNumberCreate(event eventstore.Event) error {
	var eventData events.PhoneNumberCreateEvent
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

func (a *PhoneNumberAggregate) onPhoneNumberUpdate(event eventstore.Event) error {
	var eventData events.PhoneNumberUpdatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.Source.SourceOfTruth = eventData.SourceOfTruth
	a.PhoneNumber.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *PhoneNumberAggregate) OnPhoneNumberSkippedValidation(event eventstore.Event) error {
	var eventData events.PhoneNumberSkippedValidationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.PhoneNumberValidation.SkipReason = eventData.Reason
	return nil
}

func (a *PhoneNumberAggregate) OnPhoneNumberFailedValidation(event eventstore.Event) error {
	var eventData events.PhoneNumberFailedValidationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.PhoneNumberValidation.ValidationError = eventData.ValidationError
	return nil
}

func (a *PhoneNumberAggregate) OnPhoneNumberValidated(event eventstore.Event) error {
	var eventData events.PhoneNumberValidatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.PhoneNumberValidation.ValidationError = ""
	a.PhoneNumber.PhoneNumberValidation.SkipReason = ""
	a.PhoneNumber.E164 = eventData.E164
	a.PhoneNumber.CountryCodeA2 = eventData.CountryCodeA2
	a.PhoneNumber.UpdatedAt = eventData.ValidatedAt
	return nil
}
