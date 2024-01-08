package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
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
	data := neo4jrepository.ContractCreateFields{
		OrganizationId:   eventData.OrganizationId,
		Name:             eventData.Name,
		ContractUrl:      eventData.ContractUrl,
		CreatedByUserId:  eventData.CreatedByUserId,
		ServiceStartedAt: eventData.ServiceStartedAt,
		SignedAt:         eventData.SignedAt,
		RenewalCycle:     eventData.RenewalCycle,
		RenewalPeriods:   eventData.RenewalPeriods,
		Status:           eventData.Status,
		CreatedAt:        eventData.CreatedAt,
		UpdatedAt:        eventData.UpdatedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source.Source),
			AppSource:     helper.GetAppSource(eventData.Source.AppSource),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.Source.Source),
		},
	}
	err := h.repositories.Neo4jRepositories.ContractWriteRepository.CreateForOrganization(ctx, eventData.Tenant, contractId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving contract %s: %s", contractId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, contractId, neo4jentity.NodeLabelContract, eventData.ExternalSystem)
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

	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	beforeUpdateContractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)

	data := neo4jrepository.ContractUpdateFields{
		Name:             eventData.Name,
		ContractUrl:      eventData.ContractUrl,
		ServiceStartedAt: eventData.ServiceStartedAt,
		Status:           eventData.Status,
		Source:           helper.GetSource(eventData.Source),
		RenewalPeriods:   eventData.RenewalPeriods,
		RenewalCycle:     eventData.RenewalCycle,
		UpdatedAt:        eventData.UpdatedAt,
		SignedAt:         eventData.SignedAt,
		EndedAt:          eventData.EndedAt,
	}
	updatedContractDbNode, err := h.repositories.Neo4jRepositories.ContractWriteRepository.UpdateAndReturn(ctx, eventData.Tenant, contractId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s: %s", contractId, err.Error())
		return err
	}
	afterUpdateContractEntity := graph_db.MapDbNodeToContractEntity(updatedContractDbNode)

	if eventData.ExternalSystem.Available() {
		err = h.repositories.ExternalSystemRepository.LinkWithEntity(ctx, eventData.Tenant, contractId, neo4jentity.NodeLabelContract, eventData.ExternalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link contract %s with external system %s: %s", contractId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	if beforeUpdateContractEntity.RenewalCycle != "" && afterUpdateContractEntity.RenewalCycle == "" {
		err = h.repositories.Neo4jRepositories.ContractWriteRepository.SuspendActiveRenewalOpportunity(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while suspending renewal opportunity for contract %s: %s", contractId, err.Error())
		}
		organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
			return nil
		}
		if organizationDbNode == nil {
			h.log.Errorf("Organization not found for contract %s", contractId)
			return nil
		}
		organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

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
			err = h.repositories.Neo4jRepositories.ContractWriteRepository.ActivateSuspendedRenewalOpportunity(ctx, eventData.Tenant, contractId)
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

	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)

	if model.IsFrequencyBasedRenewalCycle(contractEntity.RenewalCycle) {
		currentRenewalOpportunityDbNode, err := h.repositories.Neo4jRepositories.OpportunityReadRepository.GetOpenRenewalOpportunityForContract(ctx, eventData.Tenant, contractId)
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

	_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.Create(ctx, eventData.Tenant, contractId, neo4jentity.CONTRACT, neo4jentity.ActionContractRenewed, message, metadata, utils.Now())
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

	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)
	//we will use this boolean below to check if the status has changed
	statusChanged := contractEntity.Status != eventData.Status

	err = h.repositories.Neo4jRepositories.ContractWriteRepository.UpdateStatus(ctx, eventData.Tenant, contractId, eventData.Status, eventData.ServiceStartedAt, eventData.EndedAt)
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
	_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, contractId, neo4jentity.CONTRACT, neo4jentity.ActionContractStatusUpdated, message, metadata, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating status update action for contract %s: %s", contractId, err.Error())
	}
}

func (h *ContractEventHandler) startOnboardingIfEligible(ctx context.Context, tenant, contractId string, span opentracing.Span) {
	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	if contractDbNode == nil {
		return
	}
	contractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)

	if contractEntity.IsEligibleToStartOnboarding() {
		organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while getting organization for contract %s: %s", contractEntity.Id, err.Error())
			return
		}
		if organizationDbNode == nil {
			return
		}
		organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
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
