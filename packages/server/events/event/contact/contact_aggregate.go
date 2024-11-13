package contact

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	ContactAggregateType eventstore.AggregateType = "contact"
)

type ContactAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Contact *Contact
}

func NewContactAggregateWithTenantAndID(tenant, id string) *ContactAggregate {
	contactAggregate := ContactAggregate{}
	contactAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(ContactAggregateType, tenant, id)
	contactAggregate.SetWhen(contactAggregate.When)
	contactAggregate.Contact = &Contact{}
	contactAggregate.Tenant = tenant
	return &contactAggregate
}

func (a *ContactAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactAggregate.HandleGRPCRequest")
	defer span.Finish()

	tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
	return nil, eventstore.ErrInvalidRequestType
}

func (a *ContactAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ContactPhoneNumberLinkV1:
		return a.onPhoneNumberLink(evt)
	case event.ContactLocationLinkV1:
		return a.onLocationLink(evt)
	default:
		return nil
	}
}

func (a *ContactAggregate) onPhoneNumberLink(evt eventstore.Event) error {
	var eventData event.ContactLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.PhoneNumbers == nil {
		a.Contact.PhoneNumbers = make(map[string]ContactPhoneNumber)
	}
	a.Contact.PhoneNumbers[eventData.PhoneNumberId] = ContactPhoneNumber{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *ContactAggregate) onLocationLink(evt eventstore.Event) error {
	var eventData event.ContactLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.LocationIds = utils.AddToListIfNotExists(a.Contact.LocationIds, eventData.LocationId)
	return nil
}
