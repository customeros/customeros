package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/events"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OpportunityEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewOpportunityEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *OpportunityEventHandler {
	return &OpportunityEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

type ActionLikelihoodMetadata struct {
	Likelihood string `json:"likelihood"`
	Reason     string `json:"reason"`
}

func (h *OpportunityEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OpportunityCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	data := neo4jrepository.OpportunityCreateFields{
		OrganizationId: eventData.OrganizationId,
		CreatedAt:      eventData.CreatedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:        helper.GetSource(eventData.Source.Source),
			SourceOfTruth: helper.GetSource(eventData.Source.Source),
			AppSource:     helper.GetAppSource(eventData.Source.AppSource),
		},
		Name:              eventData.Name,
		MaxAmount:         eventData.MaxAmount,
		InternalType:      eventData.InternalType,
		ExternalType:      eventData.ExternalType,
		InternalStage:     eventData.InternalStage,
		ExternalStage:     eventData.ExternalStage,
		EstimatedClosedAt: eventData.EstimatedClosedAt,
		GeneralNotes:      eventData.GeneralNotes,
		NextSteps:         eventData.NextSteps,
		CreatedByUserId:   eventData.CreatedByUserId,
		LikelihoodRate:    eventData.LikelihoodRate,
	}
	if eventData.Currency != "" {
		data.Currency = neo4jenum.DecodeCurrency(eventData.Currency)
	} else {
		tenantSettingsDbNode, err := h.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, eventData.Tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while getting tenant settings for tenant %s: %s", eventData.Tenant, err.Error())
		}
		if tenantSettingsDbNode != nil {
			tenantSettings := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)
			data.Currency = tenantSettings.BaseCurrency
		}
	}
	err := h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.CreateForOrganization(ctx, eventData.Tenant, opportunityId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	if eventData.OwnerUserId != "" {
		err = h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, nil, eventData.Tenant, opportunityId, eventData.OwnerUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while replacing owner of opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
	}

	if eventData.ExternalSystem.Available() {
		externalSystemData := neo4jmodel.ExternalSystem{
			ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
			ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
			ExternalId:       eventData.ExternalSystem.ExternalId,
			ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
			ExternalSource:   eventData.ExternalSystem.ExternalSource,
			SyncDate:         eventData.ExternalSystem.SyncDate,
		}
		err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, opportunityId, model.NodeLabelOpportunity, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while linking opportunity %s with external system %s: %s", opportunityId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return nil
		}
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.OPPORTUNITY.String(), opportunityId, h.grpcClients, utils.NewEventCompletedDetails().WithCreate())

	return nil
}

func (h *OpportunityEventHandler) OnCreateRenewal(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnCreateRenewal")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData opportunityevent.OpportunityCreateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	// check if active renewal opportunity already exists for this contract
	opportunityDbNode, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting renewal opportunity for contract %s: %s", eventData.ContractId, err.Error())
		return nil
	}
	if opportunityDbNode != nil {
		opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
		if opportunity.RenewalDetails.RenewedAt != nil && opportunity.RenewalDetails.RenewedAt.After(utils.Now()) {
			span.LogFields(log.String("result", "active renewal opportunity already exists, skip creation"))
			h.log.Infof("active renewal opportunity already exists for contract %s", eventData.ContractId)
			return nil
		}
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	data := neo4jrepository.RenewalOpportunityCreateFields{
		ContractId: eventData.ContractId,
		CreatedAt:  eventData.CreatedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:        helper.GetSource(eventData.Source.Source),
			SourceOfTruth: helper.GetSource(eventData.Source.Source),
			AppSource:     helper.GetAppSource(eventData.Source.AppSource),
		},
		InternalType:        eventData.InternalType,
		InternalStage:       eventData.InternalStage,
		RenewalLikelihood:   eventData.RenewalLikelihood,
		RenewalApproved:     eventData.RenewalApproved,
		RenewedAt:           eventData.RenewedAt,
		RenewalAdjustedRate: eventData.RenewalAdjustedRate,
	}
	newOpportunityCreated, err := h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.CreateRenewal(ctx, eventData.Tenant, opportunityId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while saving renewal opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	if newOpportunityCreated {
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityRenewDateAndArr(ctx, eventData.Tenant, eventData.ContractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
	} else {
		// Mark event store stream for deletion
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*emptypb.Empty](func() (*emptypb.Empty, error) {
			return h.grpcClients.EventStoreClient.DeleteEventStoreStream(ctx, &eventstorepb.DeleteEventStoreStreamRequest{
				Tenant: eventData.Tenant,
				Type:   constants.AggregateTypeOpportunity,
				Id:     opportunityId,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("DeleteEventStoreStream failed: %v", err.Error())
		}
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.OPPORTUNITY.String(), opportunityId, h.grpcClients, utils.NewEventCompletedDetails().WithCreate())

	return nil
}

func (h *OpportunityEventHandler) OnUpdateNextCycleDate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnUpdateNextCycleDate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData opportunityevent.OpportunityUpdateNextCycleDateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.UpdateNextRenewalDate(ctx, eventData.Tenant, opportunityId, eventData.RenewedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating next cycle date for opportunity %s: %s", opportunityId, err.Error())
	}

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
	}
	if contractDbNode != nil {
		contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityLikelihood(ctx, eventData.Tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractEntity.Id, err.Error())
		}

		// refresh contract status
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*contractpb.ContractIdGrpcResponse](func() (*contractpb.ContractIdGrpcResponse, error) {
			return h.grpcClients.ContractClient.RefreshContractStatus(ctx, &contractpb.RefreshContractStatusGrpcRequest{
				Tenant:    eventData.Tenant,
				Id:        contractEntity.Id,
				AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshContractStatus failed: %s", err.Error())
		}
	}

	h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)

	utils.EventCompleted(ctx, eventData.Tenant, model.OPPORTUNITY.String(), opportunityId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OpportunityEventHandler) sendEventToUpdateOrganizationRenewalSummary(ctx context.Context, tenant, opportunityId string, span opentracing.Span) {
	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
		return
	}
	if organizationDbNode == nil {
		return
	}
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.RefreshRenewalSummaryGrpcRequest{
			Tenant:         tenant,
			OrganizationId: organization.ID,
			AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("RefreshRenewalSummary failed: %v", err.Error())
	}
}

func (h *OpportunityEventHandler) sendEventToUpdateOrganizationArr(ctx context.Context, tenant, opportunityId string, span opentracing.Span) {
	// if amount changed, recalculate organization combined ARR forecast
	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
		return
	}
	if organizationDbNode == nil {
		return
	}
	organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         tenant,
			OrganizationId: organization.ID,
			AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("RefreshArr failed: %v", err.Error())
	}
}

func (h *OpportunityEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OpportunityUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)

	opportunityDbNode, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting opportunity %s: %s", opportunityId, err.Error())
		return err
	}
	opportunityBeforeUpdate := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)

	data := neo4jrepository.OpportunityUpdateFields{
		Source:                  eventData.Source,
		Name:                    eventData.Name,
		Amount:                  eventData.Amount,
		MaxAmount:               eventData.MaxAmount,
		ExternalStage:           eventData.ExternalStage,
		ExternalType:            eventData.ExternalType,
		EstimatedClosedAt:       eventData.EstimatedClosedAt,
		InternalStage:           eventData.InternalStage,
		LikelihoodRate:          eventData.LikelihoodRate,
		NextSteps:               eventData.NextSteps,
		Currency:                neo4jenum.DecodeCurrency(eventData.Currency),
		UpdateName:              eventData.UpdateName(),
		UpdateAmount:            eventData.UpdateAmount(),
		UpdateMaxAmount:         eventData.UpdateMaxAmount(),
		UpdateExternalStage:     eventData.UpdateExternalStage(),
		UpdateExternalType:      eventData.UpdateExternalType(),
		UpdateEstimatedClosedAt: eventData.UpdateEstimatedClosedAt(),
		UpdateInternalStage:     eventData.UpdateInternalStage(),
		UpdateCurrency:          eventData.UpdateCurrency(),
		UpdateLikelihoodRate:    eventData.UpdateLikelihoodRate(),
		UpdateNextSteps:         eventData.UpdateNextSteps(),
	}
	err = h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.Update(ctx, eventData.Tenant, opportunityId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	//on service
	if eventData.UpdateOwnerUserId() {
		if eventData.OwnerUserId != "" {
			err = h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, nil, eventData.Tenant, opportunityId, eventData.OwnerUserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error while replacing owner of opportunity %s: %s", opportunityId, err.Error())
			}
		} else {
			err = h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.RemoveOwner(ctx, nil, eventData.Tenant, opportunityId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error while removing owner of opportunity %s: %s", opportunityId, err.Error())
			}
		}
	}

	//NOT on service
	if eventData.ExternalSystem.Available() {
		externalSystemData := neo4jmodel.ExternalSystem{
			ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
			ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
			ExternalId:       eventData.ExternalSystem.ExternalId,
			ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
			ExternalSource:   eventData.ExternalSystem.ExternalSource,
			SyncDate:         eventData.ExternalSystem.SyncDate,
		}
		err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, opportunityId, model.NodeLabelOpportunity, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while linking opportunity %s with external system %s: %s", opportunityId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	opportunityDbNode, err = h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting opportunity %s: %s", opportunityId, err.Error())
		return err
	}
	opportunityAfterUpdate := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)

	//on service
	if opportunityBeforeUpdate.InternalStage != opportunityAfterUpdate.InternalStage || opportunityBeforeUpdate.ExternalStage != opportunityAfterUpdate.ExternalStage {
		err = h.services.CommonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, eventData.Tenant, model.NodeLabelOpportunity, opportunityId, string(neo4jentity.OpportunityPropertyStageUpdatedAt), utils.NowPtr())
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while updating opportunity %s: %s", opportunityId, err.Error())
		}
	}

	//on service
	// if amount changed, recalculate organization combined ARR forecast
	if (eventData.UpdateAmount() || eventData.UpdateMaxAmount()) && opportunityBeforeUpdate.InternalType == neo4jenum.OpportunityInternalTypeRenewal {
		h.sendEventToUpdateOrganizationArr(ctx, eventData.Tenant, opportunityId, span)
	}

	if eventData.AppSource != constants.AppSourceCustomerOsApi {
		utils.EventCompleted(ctx, eventData.Tenant, model.OPPORTUNITY.String(), opportunityId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
	}
	return nil
}

func (h *OpportunityEventHandler) OnUpdateRenewal(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnUpdateRenewal")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OpportunityUpdateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	opportunityDbNode, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting opportunity %s: %s", opportunityId, err.Error())
		return err
	}
	opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
	amountChanged := eventData.UpdateAmount() && opportunity.Amount != eventData.Amount
	likelihoodChanged := eventData.UpdateRenewalLikelihood() && opportunity.RenewalDetails.RenewalLikelihood.String() != eventData.RenewalLikelihood
	adjustedRateChanged := eventData.UpdateRenewalAdjustedRate() && opportunity.RenewalDetails.RenewalAdjustedRate != eventData.RenewalAdjustedRate
	setUpdatedByUserId := (amountChanged || likelihoodChanged || adjustedRateChanged) && eventData.UpdatedByUserId != ""
	if eventData.OwnerUserId != "" {
		err = h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, nil, eventData.Tenant, opportunityId, eventData.OwnerUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while replacing owner of opportunity %s: %s", opportunityId, err.Error())
			return err
		}
	}
	data := neo4jrepository.RenewalOpportunityUpdateFields{
		UpdatedAt:                 eventData.UpdatedAt,
		Source:                    helper.GetSource(eventData.Source),
		UpdatedByUserId:           eventData.UpdatedByUserId,
		SetUpdatedByUserId:        setUpdatedByUserId,
		Comments:                  eventData.Comments,
		Amount:                    eventData.Amount,
		RenewalLikelihood:         eventData.RenewalLikelihood,
		RenewalApproved:           eventData.RenewalApproved,
		RenewedAt:                 eventData.RenewedAt,
		RenewalAdjustedRate:       eventData.RenewalAdjustedRate,
		UpdateComments:            eventData.UpdateComments(),
		UpdateAmount:              eventData.UpdateAmount(),
		UpdateRenewalLikelihood:   eventData.UpdateRenewalLikelihood(),
		UpdateRenewalApproved:     eventData.UpdateRenewalApproved(),
		UpdateRenewedAt:           eventData.UpdateRenewedAt(),
		UpdateRenewalAdjustedRate: eventData.UpdateRenewalAdjustedRate(),
	}
	err = h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.UpdateRenewal(ctx, eventData.Tenant, opportunityId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	if likelihoodChanged {
		h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)
	}
	// update renewal ARR if likelihood changed but amount didn't
	if (likelihoodChanged || adjustedRateChanged) && !amountChanged {
		contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
		if contractDbNode == nil {
			return nil
		}
		contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
	} else if amountChanged {
		h.sendEventToUpdateOrganizationArr(ctx, eventData.Tenant, opportunityId, span)
	}

	// prepare action for likelihood change
	if likelihoodChanged {
		contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
		if contractDbNode == nil {
			return nil
		}
		contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

		err = h.saveLikelihoodChangeAction(ctx, contract.Id, eventData, span)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("saveLikelihoodChangeAction failed: %v", err.Error())
		}
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.OPPORTUNITY.String(), opportunityId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OpportunityEventHandler) OnCloseLost(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnCloseLost")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData opportunityevent.OpportunityCloseLooseEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	err := h.services.CommonServices.Neo4jRepositories.OpportunityWriteRepository.CloseLost(ctx, nil, eventData.Tenant, opportunityId, eventData.ClosedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while closing opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	opportunityDbNode, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)

	//additional actions for lost opportunity

	// update organization ARR if opportunity is renewal
	if opportunity.InternalType == neo4jenum.OpportunityInternalTypeRenewal {
		h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)
		h.sendEventToUpdateOrganizationArr(ctx, eventData.Tenant, opportunityId, span)
	}

	// clean external stage
	if opportunity.InternalType == neo4jenum.OpportunityInternalTypeNBO {
		if opportunity.ExternalStage != "" {
			ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
			_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
				return h.grpcClients.OpportunityClient.UpdateOpportunity(ctx, &opportunitypb.UpdateOpportunityGrpcRequest{
					Tenant:        eventData.Tenant,
					Id:            opportunityId,
					ExternalStage: "",
					SourceFields: &commonpb.SourceFields{
						AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
						Source:    constants.SourceOpenline,
					},
					FieldsMask: []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_EXTERNAL_STAGE},
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
			}
		}
	}

	// set organization stage to target if still engaged
	if opportunity.InternalType == neo4jenum.OpportunityInternalTypeNBO {
		organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		if organizationDbNode != nil {
			organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
			// Make organization target if it's not already
			if organizationEntity.Relationship == neo4jenum.Prospect && organizationEntity.Stage == neo4jenum.Engaged {
				_, err = h.services.CommonServices.OrganizationService.Save(ctx, nil, eventData.Tenant, &organizationEntity.ID, &neo4jrepository.OrganizationSaveFields{
					Stage:       neo4jenum.Target,
					UpdateStage: true,
					SourceFields: neo4jmodel.SourceFields{
						AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
						Source:    constants.SourceOpenline,
					},
				})
				if err != nil {
					tracing.TraceErr(span, err)
				}
			}
		}
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.OPPORTUNITY.String(), opportunityId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OpportunityEventHandler) saveLikelihoodChangeAction(ctx context.Context, contractId string, eventData events.OpportunityUpdateRenewalEvent, span opentracing.Span) error {
	metadata, err := utils.ToJson(ActionLikelihoodMetadata{
		Reason:     eventData.Comments,
		Likelihood: eventData.RenewalLikelihood,
	})
	userName := ""
	if eventData.UpdatedByUserId != "" {
		userDbNode, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetUserById(ctx, eventData.Tenant, eventData.UpdatedByUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to get user %s: %s", eventData.UpdatedByUserId, err.Error())
		}
		if userDbNode != nil {
			user := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
			userName = user.GetFullName()
		}
	}
	message := fmt.Sprintf("Renewal likelihood set to %s", cases.Title(language.English).String(eventData.RenewalLikelihood))
	if userName != "" {
		message += fmt.Sprintf(" by %s", userName)
	}

	extraActionProperties := map[string]interface{}{
		"comments": eventData.Comments,
	}
	_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, model.CONTRACT, neo4jenum.ActionRenewalLikelihoodUpdated, message, metadata, eventData.UpdatedAt, constants.AppSourceEventProcessingPlatformSubscribers, extraActionProperties)
	return err
}
