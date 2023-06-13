package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
)

func (a *JobRoleAggregate) CreateJobRole(ctx context.Context, createInput *commands.CreateJobRoleCommand) error {
	return aggregate.CreateEvent(ctx, "CreateJobRole", a, createInput, func() (eventstore.Event, error) {
		return events.NewJobRoleCreateEvent(a, createInput)
	})
}
