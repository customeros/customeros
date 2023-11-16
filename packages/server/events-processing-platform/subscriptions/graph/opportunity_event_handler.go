package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OpportunityEventHandler struct {
	log                 logger.Logger
	repositories        *repository.Repositories
	opportunityCommands *opportunitycmdhandler.CommandHandlers
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

func (h *OpportunityEventHandler) OnCreateRenewal(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnCreateRenewal")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityCreateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.OpportunityRepository.CreateRenewalOpportunity(ctx, eventData.Tenant, opportunityId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while saving renewal opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
	err = contractHandler.UpdateRenewalNextCycleDate(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
		return nil
	}
	err = contractHandler.UpdateRenewalArrForecast(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
		return nil
	}

	return nil
}

func (h *OpportunityEventHandler) OnUpdateNextCycleDate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnUpdateNextCycleDate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityUpdateNextCycleDateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.OpportunityRepository.UpdateNextCycleDate(ctx, eventData.Tenant, opportunityId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating next cycle date for opportunity %s: %s", opportunityId, err.Error())
	}

	return nil
}

func (h *OpportunityEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnUpdate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.OpportunityRepository.Update(ctx, eventData.Tenant, opportunityId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving opportunity %s: %s", opportunityId, err.Error())
		return err
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
