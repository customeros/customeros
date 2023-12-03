package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	opportunitycmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	organizationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strings"
)

type ActionStatusMetadata struct {
	Status string `json:"status"`
}

type ContractEventHandler struct {
	log                  logger.Logger
	repositories         *repository.Repositories
	opportunityCommands  *opportunitycmdhandler.CommandHandlers
	organizationCommands *organizationcmdhandler.CommandHandlers
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnUpdate")
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
	beforeUpdateContractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)

	updatedContractDbNode, err := h.repositories.ContractRepository.UpdateAndReturn(ctx, eventData.Tenant, contractId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s: %s", contractId, err.Error())
		return err
	}
	afterUpdateContractEntity := graph_db.MapDbNodeToContractEntity(updatedContractDbNode)

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, contractId, constants.NodeLabel_Contract, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link contract %s with external system %s: %s", contractId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	if beforeUpdateContractEntity.RenewalCycle != "" && afterUpdateContractEntity.RenewalCycle == "" {
		err = h.repositories.ContractRepository.SuspendActiveRenewalOpportunity(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while suspending renewal opportunity for contract %s: %s", contractId, err.Error())
		}
		organizationDbNode, err := h.repositories.OrganizationRepository.GetOrganizationByContractId(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
			return nil
		}
		if organizationDbNode == nil {
			h.log.Errorf("Organization not found for contract %s", contractId)
			return nil
		}
		organization := graph_db.MapDbNodeToOrganizationEntity(*organizationDbNode)

		err = h.organizationCommands.RefreshRenewalSummary.Handle(ctx, cmd.NewRefreshRenewalSummaryCommand(eventData.Tenant, organization.ID, "", constants.AppSourceEventProcessingPlatform))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("NewRefreshRenewalSummaryCommand failed: %v", err.Error())
		}
		err = h.organizationCommands.RefreshArr.Handle(ctx, cmd.NewRefreshArrCommand(eventData.Tenant, organization.ID, "", constants.AppSourceEventProcessingPlatform))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("NewRefreshArrCommand failed: %v", err.Error())
		}
	} else {
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
		err = contractHandler.UpdateRenewalArrAndNextCycleDate(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
		}
	}

	return nil
}

func (h *ContractEventHandler) OnRolloutRenewalOpportunity(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnRolloutRenewalOpportunity")
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
	contractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)

	if model.IsFrequencyBasedRenewalCycle(contractEntity.RenewalCycle) {
		currentRenewalOpportunityDbNode, err := h.repositories.OpportunityRepository.GetOpenRenewalOpportunityForContract(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		}

		if currentRenewalOpportunityDbNode != nil {
			currentOpportunity := graph_db.MapDbNodeToOpportunityEntity(currentRenewalOpportunityDbNode)
			err := h.opportunityCommands.CloseWinOpportunity.Handle(ctx, opportunitycmd.NewCloseWinOpportunityCommand(currentOpportunity.Id, eventData.Tenant, "", constants.AppSourceEventProcessingPlatform, nil, nil))
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("CloseWinOpportunity failed: %v", err.Error())
			}
		}

		err = h.opportunityCommands.CreateRenewalOpportunity.Handle(ctx, opportunitycmd.NewCreateRenewalOpportunityCommand("", eventData.Tenant, "", contractId, "", commonmodel.Source{
			Source:    constants.SourceOpenline,
			AppSource: constants.AppSourceEventProcessingPlatform,
		}, nil, nil))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CreateRenewalOpportunity failed: %v", err.Error())
		}
	}

	return nil
}

func (h *ContractEventHandler) OnUpdateStatus(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnUpdateStatus")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.ContractUpdateStatusEvent
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
	contractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)
	//we will use this boolean below to check if the status has changed
	statusChanged := contractEntity.Status != eventData.Status

	err = h.repositories.ContractRepository.UpdateStatus(ctx, eventData.Tenant, contractId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
		return nil
	}

	if eventData.Status == string(model.ContractStatusStringEnded) {
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
		err := contractHandler.UpdateRenewalNextCycleDate(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating contract's {%s} renewal date: %s", contractId, err.Error())
		}
	}
	var message string
	metadata, err := utils.ToJson(ActionStatusMetadata{
		Status: eventData.Status,
	})

	if statusChanged {
		switch eventData.Status {
		case string(model.ContractStatusStringLive):
			message = contractEntity.Name + " is now " + strings.ToLower(eventData.Status)
		case string(model.ContractStatusStringEnded):
			message = contractEntity.Name + " has " + strings.ToLower(eventData.Status)
		}
		_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractId, entity.CONTRACT, entity.ActionContractStatusUpdated, message, metadata, utils.Now())
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating status update action for contract %s: %s", contractId, err.Error())
		}
	}
	return nil
}
