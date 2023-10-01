package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
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
	*aggregate.CommonTenantIdAggregate
	User *models.User
}

func NewUserAggregateWithTenantAndID(tenant, id string) *UserAggregate {
	userAggregate := UserAggregate{}
	userAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(UserAggregateType, tenant, id)
	userAggregate.SetWhen(userAggregate.When)
	userAggregate.User = &models.User{}
	userAggregate.Tenant = tenant

	return &userAggregate
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
	case events.UserAddPlayerV1:
		return a.onAddPlayer(event)

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
	a.User.Source = eventData.SourceFields
	a.User.CreatedAt = eventData.CreatedAt
	a.User.UpdatedAt = eventData.UpdatedAt
	a.User.Timezone = eventData.Timezone
	if eventData.ExternalSystem.Available() {
		a.User.ExternalSystems = []common_models.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *UserAggregate) onUserUpdate(event eventstore.Event) error {
	var eventData events.UserUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if eventData.Source != a.User.Source.SourceOfTruth && a.User.Source.SourceOfTruth == constants.SourceOpenline {
		if a.User.Name == "" {
			a.User.Name = eventData.Name
		}
		if a.User.FirstName == "" {
			a.User.FirstName = eventData.FirstName
		}
		if a.User.LastName == "" {
			a.User.LastName = eventData.LastName
		}
		if a.User.Timezone == "" {
			a.User.Timezone = eventData.Timezone
		}
		if a.User.ProfilePhotoUrl == "" {
			a.User.ProfilePhotoUrl = eventData.ProfilePhotoUrl
		}
	} else {
		a.User.Name = eventData.Name
		a.User.FirstName = eventData.FirstName
		a.User.LastName = eventData.LastName
		a.User.Timezone = eventData.Timezone
		a.User.ProfilePhotoUrl = eventData.ProfilePhotoUrl
	}

	a.User.UpdatedAt = eventData.UpdatedAt
	a.User.Internal = eventData.Internal
	if eventData.Source == constants.SourceOpenline {
		a.User.Source.SourceOfTruth = eventData.Source
	}
	if eventData.ExternalSystem.Available() {
		found := false
		for _, externalSystem := range a.User.ExternalSystems {
			if externalSystem.ExternalSystemId == eventData.ExternalSystem.ExternalSystemId &&
				externalSystem.ExternalId == eventData.ExternalSystem.ExternalId {
				found = true
				externalSystem.ExternalUrl = eventData.ExternalSystem.ExternalUrl
				externalSystem.SyncDate = eventData.ExternalSystem.SyncDate
				externalSystem.ExternalSource = eventData.ExternalSystem.ExternalSource
				if eventData.ExternalSystem.ExternalIdSecond != "" {
					externalSystem.ExternalIdSecond = eventData.ExternalSystem.ExternalIdSecond
				}
			}
		}
		if !found {
			a.User.ExternalSystems = append(a.User.ExternalSystems, eventData.ExternalSystem)
		}
	}
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

func (a *UserAggregate) onAddPlayer(event eventstore.Event) error {
	var eventData events.UserAddPlayerInfoEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.User.Players == nil {
		a.User.Players = make([]models.PlayerInfo, 0)
	}
	found := false
	for _, player := range a.User.Players {
		if player.AuthId == eventData.AuthId && player.Provider == eventData.Provider {
			found = true
			player.IdentityId = eventData.IdentityId
		}
	}
	if !found {
		a.User.Players = append(a.User.Players, models.PlayerInfo{
			AuthId:     eventData.AuthId,
			Provider:   eventData.Provider,
			IdentityId: eventData.IdentityId,
		})
	}

	return nil
}
