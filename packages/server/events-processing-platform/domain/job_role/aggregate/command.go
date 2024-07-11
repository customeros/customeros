package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

func (a *JobRoleAggregate) CreateJobRole(ctx context.Context, createInput *model.CreateJobRoleCommand) error {
	return aggregate.CreateEvent(ctx, "CreateJobRole", a, createInput, func() (eventstore.Event, error) {
		return events.NewJobRoleCreateEvent(a, createInput)
	})
}
