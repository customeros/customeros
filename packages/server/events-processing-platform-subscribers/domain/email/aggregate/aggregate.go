package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	EmailAggregateType eventstore.AggregateType = "email"
)

type EmailAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Email *models.Email
}

func NewEmailAggregateWithTenantAndID(tenant, id string) *EmailAggregate {
	emailAggregate := EmailAggregate{}
	emailAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(EmailAggregateType, tenant, id)
	emailAggregate.SetWhen(emailAggregate.When)
	emailAggregate.Email = &models.Email{}
	emailAggregate.Tenant = tenant

	return &emailAggregate
}

func (a *EmailAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {

	case events.EmailCreateV1:
		return a.onEmailCreate(event)
	case events.EmailUpdateV1:
		return a.onEmailUpdated(event)
	case events.EmailValidationFailedV1:
		return a.OnEmailFailedValidation(event)
	case events.EmailValidatedV1:
		return a.OnEmailValidated(event)

	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *EmailAggregate) onEmailCreate(event eventstore.Event) error {
	var eventData events.EmailCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Email.RawEmail = eventData.RawEmail
	if eventData.SourceFields.Available() {
		a.Email.Source = eventData.SourceFields
	} else {
		a.Email.Source.Source = eventData.Source
		a.Email.Source.SourceOfTruth = eventData.SourceOfTruth
		a.Email.Source.AppSource = eventData.AppSource
	}
	a.Email.CreatedAt = eventData.CreatedAt
	a.Email.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *EmailAggregate) onEmailUpdated(event eventstore.Event) error {
	var eventData events.EmailUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.Email.Source.SourceOfTruth = eventData.Source
	}
	a.Email.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *EmailAggregate) OnEmailFailedValidation(event eventstore.Event) error {
	var eventData events.EmailFailedValidationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Email.EmailValidation.ValidationError = eventData.ValidationError
	return nil
}

func (a *EmailAggregate) OnEmailValidated(event eventstore.Event) error {
	var eventData events.EmailValidatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Email.Email = eventData.EmailAddress
	a.Email.EmailValidation.IsReachable = eventData.IsReachable
	a.Email.EmailValidation.ValidationError = eventData.ValidationError
	a.Email.EmailValidation.AcceptsMail = eventData.AcceptsMail
	a.Email.EmailValidation.CanConnectSmtp = eventData.CanConnectSmtp
	a.Email.EmailValidation.HasFullInbox = eventData.HasFullInbox
	a.Email.EmailValidation.IsCatchAll = eventData.IsCatchAll
	a.Email.EmailValidation.IsDeliverable = eventData.IsDeliverable
	a.Email.EmailValidation.IsDisabled = eventData.IsDisabled
	a.Email.EmailValidation.Domain = eventData.Domain
	a.Email.EmailValidation.IsValidSyntax = eventData.IsValidSyntax
	a.Email.EmailValidation.Username = eventData.Username
	return nil
}
