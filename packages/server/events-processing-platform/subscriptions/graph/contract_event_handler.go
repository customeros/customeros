package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ActionStatusMetadata struct {
	Status       string `json:"status"`
	ContractName string `json:"contract-name"`
	Comment      string `json:"comment"`
}

type ContractEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewContractEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *ContractEventHandler {
	return &ContractEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (h *ContractEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = h.grpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
			Tenant:     eventData.Tenant,
			ContractId: contractId,
			SourceFields: &commonpb.SourceFields{
				Source:    eventData.Source.Source,
				AppSource: constants.AppSourceEventProcessingPlatform,
			},
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CreateRenewalOpportunity failed: %s", err.Error())
		}
	}

	h.startOnboardingIfEligible(ctx, eventData.Tenant, contractId, span)

	return nil
}

func (h *ContractEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = h.grpcClients.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: organization.ID,
			AppSource:      constants.AppSourceEventProcessingPlatform,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshRenewalSummary failed: %v", err.Error())
		}
		_, err = h.grpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: organization.ID,
			AppSource:      constants.AppSourceEventProcessingPlatform,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshArr failed: %v", err.Error())
		}
	} else {
		if beforeUpdateContractEntity.RenewalCycle == "" && afterUpdateContractEntity.RenewalCycle != "" {
			err = h.repositories.ContractRepository.ActivateSuspendedRenewalOpportunity(ctx, eventData.Tenant, contractId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error while activating renewal opportunity for contract %s: %s", contractId, err.Error())
			}
		}
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityRenewDateAndArr(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
		}
	}

	if beforeUpdateContractEntity.Status != afterUpdateContractEntity.Status {
		h.createActionForStatusChange(ctx, eventData.Tenant, contractId, afterUpdateContractEntity.Status, afterUpdateContractEntity.Name, span)
	}

	contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
	err = contractHandler.UpdateActiveRenewalOpportunityLikelihood(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
	}

	h.startOnboardingIfEligible(ctx, eventData.Tenant, contractId, span)

	return nil
}

func (h *ContractEventHandler) OnRolloutRenewalOpportunity(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnRolloutRenewalOpportunity")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		if currentRenewalOpportunityDbNode != nil {
			currentOpportunity := graph_db.MapDbNodeToOpportunityEntity(currentRenewalOpportunityDbNode)
			_, err = h.grpcClients.OpportunityClient.CloseWinOpportunity(ctx, &opportunitypb.CloseWinOpportunityGrpcRequest{
				Tenant:    eventData.Tenant,
				Id:        currentOpportunity.Id,
				AppSource: constants.AppSourceEventProcessingPlatform,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("CloseWinOpportunity failed: %s", err.Error())
			}
		}

		_, err = h.grpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
			Tenant:     eventData.Tenant,
			ContractId: contractId,
			SourceFields: &commonpb.SourceFields{
				Source:    constants.SourceOpenline,
				AppSource: constants.AppSourceEventProcessingPlatform,
			},
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CreateRenewalOpportunity failed: %v", err.Error())
		}
	}
	status := "Renewed"
	metadata, err := utils.ToJson(ActionStatusMetadata{
		Status: status,
	})
	message := contractEntity.Name + " renewed"

	_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractId, entity.CONTRACT, entity.ActionContractRenewed, message, metadata, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating renewed action for contract %s: %s", contractId, err.Error())
	}

	return nil
}

func (h *ContractEventHandler) OnUpdateStatus(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnUpdateStatus")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
		err := contractHandler.UpdateActiveRenewalOpportunityNextCycleDate(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating contract's {%s} renewal date: %s", contractId, err.Error())
		}
	}

	if statusChanged {
		h.createActionForStatusChange(ctx, eventData.Tenant, contractId, eventData.Status, contractEntity.Name, span)
	}

	h.startOnboardingIfEligible(ctx, eventData.Tenant, contractId, span)

	return nil
}

func (h *ContractEventHandler) createActionForStatusChange(ctx context.Context, tenant, contractId, status, contractName string, span opentracing.Span) {
	span, ctx = opentracing.StartSpanFromContext(ctx, "ContractEventHandler.createActionForStatusChange")
	defer span.Finish()
	var name string
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId), log.String("status", status), log.String("contractName", contractName))

	if contractName != "" {
		name = contractName
	} else {
		name = "Unnamed contract"
	}
	metadata, err := utils.ToJson(ActionStatusMetadata{
		Status:       status,
		ContractName: name,
		Comment:      name + " is now " + status,
	})
	message := ""

	switch status {
	case string(model.ContractStatusStringLive):
		message = contractName + " is now live"
	case string(model.ContractStatusStringEnded):
		message = contractName + " has ended"
	}
	_, err = h.repositories.ActionRepository.Create(ctx, tenant, contractId, entity.CONTRACT, entity.ActionContractStatusUpdated, message, metadata, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating status update action for contract %s: %s", contractId, err.Error())
	}
}

func (h *ContractEventHandler) startOnboardingIfEligible(ctx context.Context, tenant, contractId string, span opentracing.Span) {
	contractDbNode, err := h.repositories.ContractRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	if contractDbNode == nil {
		return
	}
	contractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)

	if contractEntity.IsEligibleToStartOnboarding() {
		organizationDbNode, err := h.repositories.OrganizationRepository.GetOrganizationByContractId(ctx, tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while getting organization for contract %s: %s", contractEntity.Id, err.Error())
			return
		}
		if organizationDbNode == nil {
			return
		}
		organization := graph_db.MapDbNodeToOrganizationEntity(*organizationDbNode)
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = h.grpcClients.OrganizationClient.UpdateOnboardingStatus(ctx, &organizationpb.UpdateOnboardingStatusGrpcRequest{
			Tenant:             tenant,
			OrganizationId:     organization.ID,
			CausedByContractId: contractEntity.Id,
			OnboardingStatus:   organizationpb.OnboardingStatus_ONBOARDING_STATUS_NOT_STARTED,
			AppSource:          constants.AppSourceEventProcessingPlatform,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("UpdateOnboardingStatus gRPC request failed: %v", err.Error())
		}
	}
}
