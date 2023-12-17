package graph

import (
	"context"
	common_aggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type JobRoleEventHandler struct {
	Repositories *repository.Repositories
}

func NewJobRoleEventHandler(repositories *repository.Repositories) *JobRoleEventHandler {
	return &JobRoleEventHandler{Repositories: repositories}
}

func (h *JobRoleEventHandler) OnJobRoleCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleEventHandler.OnJobRoleCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.JobRoleCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	eventId := common_aggregate.GetAggregateWithTenantAndIdObjectID(evt.AggregateID, aggregate.JobRoleAggregateType, eventData.Tenant)
	err := h.Repositories.JobRoleRepository.CreateJobRole(ctx, eventData.Tenant, eventId, eventData)
	return err
}
