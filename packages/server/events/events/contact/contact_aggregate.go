package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

const (
	ContactAggregateType eventstore.AggregateType = "contact"
)

type ContactAggregate struct {
	*eventstore.CommonTenantIdAggregate
}

func NewContactAggregateWithTenantAndID(tenant, id string) *ContactAggregate {
	contactAggregate := ContactAggregate{}
	contactAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(ContactAggregateType, tenant, id)
	contactAggregate.Tenant = tenant
	return &contactAggregate
}
