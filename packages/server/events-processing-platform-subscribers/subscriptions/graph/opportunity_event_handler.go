package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/opportunity"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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
		err = h.services.CommonServices.ContractService.UpdateActiveRenewalOpportunityLikelihood(ctx, eventData.Tenant, contractEntity.Id)
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
