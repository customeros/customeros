package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OpportunityEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func (h *OpportunityEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnCreate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.OpportunityRepository.CreateForOrganization(ctx, eventData.Tenant, opportunityId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	if eventData.OwnerUserId != "" {
		err = h.repositories.OpportunityRepository.ReplaceOwner(ctx, eventData.Tenant, opportunityId, eventData.OwnerUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while replacing owner of opportunity %s: %s", opportunityId, err.Error())
			return err
		}
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, opportunityId, constants.NodeLabel_Opportunity, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while linking opportunity %s with external system %s: %s", opportunityId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	return nil
}
