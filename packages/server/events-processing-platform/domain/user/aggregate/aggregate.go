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

func (a *UserAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {

	case events.UserCreateV1:
		return a.onUserCreate(event)
	case events.UserJobRoleLinkV1:
		return a.onJobRoleLink(event)
	case events.UserUpdateV1:
		return a.onUserUpdate(event)
	case events.UserPhoneNumberLinkV1:
		return a.onPhoneNumberLink(event)
	case events.UserEmailLinkV1:
		return a.onEmailLink(event)

	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *UserAggregate) onUserCreate(event eventstore.Event) error {
	var eventData events.UserCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.User.Name = eventData.Name
	a.User.FirstName = eventData.FirstName
	a.User.LastName = eventData.LastName
	a.User.Internal = eventData.Internal
	a.User.ProfilePhotoUrl = eventData.ProfilePhotoUrl
	a.User.Source = common_models.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.User.CreatedAt = eventData.CreatedAt
	a.User.UpdatedAt = eventData.UpdatedAt
	a.User.Timezone = eventData.Timezone
	return nil
}

func (a *UserAggregate) onUserUpdate(event eventstore.Event) error {
	var eventData events.UserUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.User.Source.SourceOfTruth = eventData.SourceOfTruth
	a.User.UpdatedAt = eventData.UpdatedAt
	a.User.Name = eventData.Name
	a.User.FirstName = eventData.FirstName
	a.User.LastName = eventData.LastName
	a.User.Internal = eventData.Internal
	a.User.ProfilePhotoUrl = eventData.ProfilePhotoUrl
	a.User.Timezone = eventData.Timezone
	return nil
}

func (a *UserAggregate) onPhoneNumberLink(event eventstore.Event) error {
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

func (a *UserAggregate) onEmailLink(event eventstore.Event) error {
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

func (a *UserAggregate) onJobRoleLink(event eventstore.Event) error {
	var eventData events.UserLinkJobRoleEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.User.JobRoles == nil {
		a.User.JobRoles = make(map[string]bool)
	}
	a.User.JobRoles[eventData.JobRoleId] = true
	a.User.UpdatedAt = eventData.UpdatedAt
	return nil
}
