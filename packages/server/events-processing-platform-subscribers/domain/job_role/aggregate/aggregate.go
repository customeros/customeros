package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

const (
	JobRoleAggregateType eventstore.AggregateType = "job_role"
)

type JobRoleAggregate struct {
	*aggregate.CommonTenantIdAggregate
}

func NewJobRoleAggregateWithTenantAndID(tenant, id string) *JobRoleAggregate {
	jobAggregate := JobRoleAggregate{}
	jobAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(JobRoleAggregateType, tenant, id)
	return &jobAggregate
}
