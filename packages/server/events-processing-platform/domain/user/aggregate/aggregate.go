package aggregate

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	UserAggregateType eventstore.AggregateType = "user"
)

type UserAggregate struct {
	*eventstore.AggregateBase
	User *models.User
}

func NewUserAggregateWithTenantAndID(tenant, id string) *UserAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewUserAggregate()
	aggregate.SetID(tenant + "-" + id)
	return aggregate
}

func NewUserAggregate() *UserAggregate {
	userAggregate := &UserAggregate{User: models.NewUser()}
	base := eventstore.NewAggregateBase(userAggregate.When)
	base.SetType(UserAggregateType)
	userAggregate.AggregateBase = base
	return userAggregate
}

func (userAggregate *UserAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {

	case events.UserCreatedV1:
		return userAggregate.onUserCreated(event)
	case events.UserUpdatedV1:
		return userAggregate.onUserUpdated(event)
	case events.UserPhoneNumberLinkedV1:
		return userAggregate.onPhoneNumberLinked(event)
	case events.UserEmailLinkedV1:
		return userAggregate.onEmailLinked(event)

	default:
		return eventstore.ErrInvalidEventType
	}
}

func (a *UserAggregate) onUserCreated(event eventstore.Event) error {
	var eventData events.UserCreatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.User.Name = eventData.Name
	a.User.FirstName = eventData.FirstName
	a.User.LastName = eventData.LastName
	a.User.Source = common_models.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.User.CreatedAt = eventData.CreatedAt
	a.User.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *UserAggregate) onUserUpdated(event eventstore.Event) error {
	var eventData events.UserUpdatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.User.Source.SourceOfTruth = eventData.SourceOfTruth
	a.User.UpdatedAt = eventData.UpdatedAt
	a.User.Name = eventData.Name
	a.User.FirstName = eventData.FirstName
	a.User.LastName = eventData.LastName
	return nil
}

func (a *UserAggregate) onPhoneNumberLinked(event eventstore.Event) error {
	var eventData events.UserLinkPhoneNumberEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.User.PhoneNumbers == nil {
		a.User.PhoneNumbers = make(map[string]models.UserPhoneNumber)
	}
	a.User.PhoneNumbers[eventData.PhoneNumberId] = models.UserPhoneNumber{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.User.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *UserAggregate) onEmailLinked(event eventstore.Event) error {
	var eventData events.UserLinkEmailEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.User.Emails == nil {
		a.User.Emails = make(map[string]models.UserEmail)
	}
	a.User.Emails[eventData.EmailId] = models.UserEmail{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.User.UpdatedAt = eventData.UpdatedAt
	return nil
}
