package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

const (
	ContactAggregateType eventstore.AggregateType = "contact"
)

type Contact struct {
}

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

func (a *ContactAggregate) When(evt eventstore.Event) error {
	return nil
}
