package aggregate

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	EmailAggregateType eventstore.AggregateType = "email"
)

type EmailAggregate struct {
	*eventstore.AggregateBase
	Email *models.Email
}

func NewEmailAggregateWithTenantAndID(tenant, id string) *EmailAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewEmailAggregate()
	aggregate.SetID(tenant + "-" + id)
	return aggregate
}

func NewEmailAggregate() *EmailAggregate {
	emailAggregate := &EmailAggregate{Email: models.NewEmail()}
	base := eventstore.NewAggregateBase(emailAggregate.When)
	base.SetType(EmailAggregateType)
	emailAggregate.AggregateBase = base
	return emailAggregate
}

func (a *EmailAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {

	case events.EmailCreatedV1:
		return a.onEmailCreated(event)
	case events.EmailUpdatedV1:
		return a.onEmailUpdated(event)

	default:
		return eventstore.ErrInvalidEventType
	}
}

func (a *EmailAggregate) onEmailCreated(event eventstore.Event) error {
	var eventData events.EmailCreatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Email.RawEmail = eventData.RawEmail
	a.Email.Source = common_models.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Email.CreatedAt = eventData.CreatedAt
	a.Email.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *EmailAggregate) onEmailUpdated(event eventstore.Event) error {
	var eventData events.EmailUpdatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Email.Source.SourceOfTruth = eventData.SourceOfTruth
	a.Email.UpdatedAt = eventData.UpdatedAt
	return nil
}
