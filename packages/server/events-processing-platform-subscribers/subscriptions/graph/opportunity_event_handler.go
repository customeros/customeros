package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
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
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
		err = h.services.CommonServices.ContractService.RefreshContractStatus(ctx, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshContractStatus failed: %s", err.Error())
		}
	}

	h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)

	h.services.CommonServices.RabbitMQService.PublishEventCompleted(ctx, eventData.Tenant, opportunityId, model.OPPORTUNITY, utils.NewEventCompletedDetails().WithUpdate())

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
	opportunityDbNode, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, nil, eventData.Tenant, opportunityId)
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

	h.services.CommonServices.RabbitMQService.PublishEventCompleted(ctx, eventData.Tenant, opportunityId, model.OPPORTUNITY, utils.NewEventCompletedDetails().WithUpdate())

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

	opportunityDbNode, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, nil, eventData.Tenant, opportunityId)
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
			_, err = h.services.CommonServices.OpportunityService.Save(ctx, nil, &opportunityId, &data_fields.OpportunityFields{
				ExternalStage: utils.ToPtr(""),
			})
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("error in UpdateOpportunity: %v", err.Error())
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
			if organizationEntity.Relationship == neo4jenum.OrganizationRelationshipProspect && organizationEntity.Stage == neo4jenum.Engaged {
				_, err = h.services.CommonServices.OrganizationService.Save(ctx, nil, &organizationEntity.ID, data_fields.OrganizationFields{
					Stage: utils.ToPtr(neo4jenum.Target),
				})
				if err != nil {
					tracing.TraceErr(span, err)
				}
			}
		}
	}

	h.services.CommonServices.RabbitMQService.PublishEventCompleted(ctx, eventData.Tenant, opportunityId, model.OPPORTUNITY, utils.NewEventCompletedDetails().WithUpdate())

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
