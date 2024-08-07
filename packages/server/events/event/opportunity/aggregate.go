package opportunity

import "github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"

const (
	OpportunityAggregateType eventstore.AggregateType = "opportunity"
)

type OpportunityAggregate struct {
	*eventstore.CommonTenantIdAggregate
}

func NewOpportunityAggregateWithTenantAndID(tenant, id string) *OpportunityAggregate {
	oppAggregate := OpportunityAggregate{}
	oppAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(OpportunityAggregateType, tenant, id)
	oppAggregate.Tenant = tenant

	return &oppAggregate
}
