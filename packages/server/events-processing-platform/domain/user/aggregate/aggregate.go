package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

const (
	UserAggregateType eventstore.AggregateType = "user"
)

type UserAggregate struct {
	*eventstore.CommonTenantIdAggregate
	User        *models.User
	EventHashes map[string]map[string]string
}

func NewUserAggregateWithTenantAndID(tenant, id string) *UserAggregate {
	userAggregate := UserAggregate{}
	userAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(UserAggregateType, tenant, id)
	userAggregate.SetWhen(userAggregate.When)
	userAggregate.User = &models.User{}
	userAggregate.Tenant = tenant
	userAggregate.EventHashes = make(map[string]map[string]string)
	return &userAggregate
}

func (a *UserAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch request.(type) {
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
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
	case events.UserEmailLinkV1,
		events.UserEmailUnlinkV1:
		return nil
	case events.UserAddRoleV1:
		return a.onAddRole(event)
	case events.UserRemoveRoleV1:
		return a.onRemoveRole(event)
	case events.UserAddPlayerV1:
		return nil
	default:
		if strings.HasPrefix(event.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
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
	a.User.Bot = eventData.Bot
	a.User.Test = eventData.Test
	a.User.ProfilePhotoUrl = eventData.ProfilePhotoUrl
	a.User.Source = eventData.SourceFields
	a.User.CreatedAt = eventData.CreatedAt
	a.User.UpdatedAt = eventData.UpdatedAt
	a.User.Timezone = eventData.Timezone
	if eventData.ExternalSystem.Available() {
		a.User.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
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
				externalSystem.ExternalSource = eventData.ExternalSystem.ExternalSource
				externalSystem.SyncDate = eventData.ExternalSystem.SyncDate
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

func (a *UserAggregate) onAddRole(event eventstore.Event) error {
	var eventData events.UserAddRoleEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.User.Roles = utils.AddToListIfNotExists(a.User.Roles, eventData.Role)
	a.User.UpdatedAt = eventData.At
	return nil
}

func (a *UserAggregate) onRemoveRole(event eventstore.Event) error {
	var eventData events.UserRemoveRoleEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.User.Roles = utils.RemoveFromList(a.User.Roles, eventData.Role)
	a.User.UpdatedAt = eventData.At
	return nil
}
