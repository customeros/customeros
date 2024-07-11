package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

const (
	JobRoleAggregateType eventstore.AggregateType = "job_role"
)

type JobRoleAggregate struct {
	*eventstore.CommonTenantIdAggregate
}

func NewJobRoleAggregateWithTenantAndID(tenant, id string) *JobRoleAggregate {
	jobAggregate := JobRoleAggregate{}
	jobAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(JobRoleAggregateType, tenant, id)
	return &jobAggregate
}
