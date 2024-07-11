package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
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

	eventId := eventstore.GetAggregateWithTenantAndIdObjectID(evt.AggregateID, aggregate.JobRoleAggregateType, eventData.Tenant)
	data := neo4jrepository.JobRoleCreateFields{
		Description: utils.IfNotNilString(eventData.Description),
		JobTitle:    eventData.JobTitle,
		StartedAt:   eventData.StartedAt,
		EndedAt:     eventData.EndedAt,
		SourceFields: neo4jmodel.Source{
			Source:        eventData.Source,
			SourceOfTruth: eventData.SourceOfTruth,
			AppSource:     eventData.AppSource,
		},
		CreatedAt: eventData.CreatedAt,
	}
	err := h.Repositories.Neo4jRepositories.JobRoleWriteRepository.CreateJobRole(ctx, eventData.Tenant, eventId, data)
	return err
}
