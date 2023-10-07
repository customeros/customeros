package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	PhoneNumberAggregateType eventstore.AggregateType = "phone_number"
)

type PhoneNumberAggregate struct {
	*aggregate.CommonTenantIdAggregate
	PhoneNumber *models.PhoneNumber
}

func NewPhoneNumberAggregateWithTenantAndID(tenant, id string) *PhoneNumberAggregate {
	phoneNumberAggregate := PhoneNumberAggregate{}
	phoneNumberAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(PhoneNumberAggregateType, tenant, id)
	phoneNumberAggregate.SetWhen(phoneNumberAggregate.When)
	phoneNumberAggregate.PhoneNumber = &models.PhoneNumber{}
	phoneNumberAggregate.Tenant = tenant

	return &phoneNumberAggregate
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
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *PhoneNumberAggregate) onPhoneNumberCreate(event eventstore.Event) error {
	var eventData events.PhoneNumberCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.PhoneNumber.RawPhoneNumber = eventData.RawPhoneNumber
	a.PhoneNumber.CreatedAt = eventData.CreatedAt
	a.PhoneNumber.UpdatedAt = eventData.UpdatedAt
	if eventData.SourceFields.Available() {
		a.PhoneNumber.Source = eventData.SourceFields
	} else {
		a.PhoneNumber.Source.Source = eventData.Source
		a.PhoneNumber.Source.SourceOfTruth = eventData.SourceOfTruth
		a.PhoneNumber.Source.AppSource = eventData.AppSource
	}
	return nil
}

func (a *PhoneNumberAggregate) onPhoneNumberUpdate(event eventstore.Event) error {
	var eventData events.PhoneNumberUpdatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.PhoneNumber.Source.SourceOfTruth = eventData.Source
	}
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
