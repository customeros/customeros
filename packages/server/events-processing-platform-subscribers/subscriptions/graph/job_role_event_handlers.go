package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type JobRoleEventHandler struct {
	services *service.Services
}

func NewJobRoleEventHandler(services *service.Services) *JobRoleEventHandler {
	return &JobRoleEventHandler{services: services}
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
	data := neo4jrepository.JobRoleFields{
		Description: utils.IfNotNilString(eventData.Description),
		JobTitle:    eventData.JobTitle,
		StartedAt:   eventData.StartedAt,
		EndedAt:     eventData.EndedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:        eventData.Source,
			SourceOfTruth: eventData.SourceOfTruth,
			AppSource:     eventData.AppSource,
		},
	}
	err := h.services.CommonServices.Neo4jRepositories.JobRoleWriteRepository.CreateJobRole(ctx, eventData.Tenant, eventId, data)
	return err
}
