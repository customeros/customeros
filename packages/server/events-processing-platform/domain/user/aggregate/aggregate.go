package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
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

	case events.UserJobRoleLinkV1:
		return a.onJobRoleLink(event)
	case events.UserPhoneNumberLinkV1:
		return a.onPhoneNumberLink(event)
	case events.UserAddRoleV1:
		return a.onAddRole(event)
	case events.UserRemoveRoleV1:
		return a.onRemoveRole(event)
	default:
		return nil
	}
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
