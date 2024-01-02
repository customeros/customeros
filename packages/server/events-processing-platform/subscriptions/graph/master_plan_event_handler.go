package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type MasterPlanEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewMasterPlanEventHandler(log logger.Logger, repositories *repository.Repositories) *MasterPlanEventHandler {
	return &MasterPlanEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *MasterPlanEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.MasterPlanCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	masterPlanId := aggregate.GetMasterPlanObjectID(evt.GetAggregateID(), eventData.Tenant)
	source := helper.GetSource(eventData.SourceFields.Source)
	appSource := helper.GetAppSource(eventData.SourceFields.AppSource)
	err := h.repositories.Neo4jRepositories.MasterPlanWriteRepository.Create(ctx, eventData.Tenant, masterPlanId, eventData.Name,
		source, appSource, eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving master plan %s: %s", masterPlanId, err.Error())
		return err
	}
	return err
}
