package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitycmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ContractEventHandler struct {
	log                 logger.Logger
	repositories        *repository.Repositories
	opportunityCommands *opportunitycmdhandler.CommandHandlers
}

func (h *ContractEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnCreate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.ContractCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	contractId := aggregate.GetContractObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.ContractRepository.CreateForOrganization(ctx, eventData.Tenant, contractId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving contract %s: %s", contractId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, contractId, constants.NodeLabel_Contract, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while linking contract %s with external system %s: %s", contractId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	if model.IsFrequencyBasedRenewalCycle(eventData.RenewalCycle) {
		err = h.opportunityCommands.CreateRenewalOpportunity.Handle(ctx, opportunitycmd.NewCreateRenewalOpportunityCommand("", eventData.Tenant, "", contractId, "", eventData.Source, nil, nil))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CreateRenewalOpportunity failed: %v", err.Error())
		}
	}

	return nil
}

func (h *ContractEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnCreate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.ContractUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contractId := aggregate.GetContractObjectID(evt.GetAggregateID(), eventData.Tenant)

	contractDbNode, err := h.repositories.ContractRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntityBeforeUpdate := graph_db.MapDbNodeToContractEntity(contractDbNode)

	contractDbNode, err = h.repositories.ContractRepository.UpdateAndReturn(ctx, eventData.Tenant, contractId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s: %s", contractId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, contractId, constants.NodeLabel_Contract, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link contract %s with external system %s: %s", contractId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	contractEntityAfterUpdate := graph_db.MapDbNodeToContractEntity(contractDbNode)
	if contractEntityAfterUpdate == nil {
		err = errors.New(fmt.Sprintf("contract {%s} not found after update", contractId))
		tracing.TraceErr(span, err)
		h.log.Error(err)
		return nil
	}

	if contractEntityBeforeUpdate.RenewalCycle == "" && model.IsFrequencyBasedRenewalCycle(contractEntityAfterUpdate.RenewalCycle) {
		err = h.opportunityCommands.CreateRenewalOpportunity.Handle(ctx, opportunitycmd.NewCreateRenewalOpportunityCommand("", eventData.Tenant, "", contractId, "", commonmodel.Source{Source: eventData.Source}, nil, nil))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CreateRenewalOpportunity failed: %v", err.Error())
		}
	} else {
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
		err = contractHandler.UpdateRenewalArr(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
			return nil
		}
	}
	if contractEntityBeforeUpdate.RenewalCycle != "" &&
		contractEntityBeforeUpdate.RenewalCycle != contractEntityAfterUpdate.RenewalCycle {
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
		err = contractHandler.UpdateRenewalNextCycleDate(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating contract's {%s} renewal date: %s", contractId, err.Error())
			return nil
		}
	}

	return nil
}
