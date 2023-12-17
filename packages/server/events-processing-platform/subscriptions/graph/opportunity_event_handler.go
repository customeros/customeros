package graph

import (
	"context"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type OpportunityEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewOpportunityEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *OpportunityEventHandler {
	return &OpportunityEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (h *OpportunityEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityCreateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.OpportunityRepository.CreateRenewal(ctx, eventData.Tenant, opportunityId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while saving renewal opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
	err = contractHandler.UpdateActiveRenewalOpportunityRenewDateAndArr(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
		return nil
	}

	h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)

	return nil
}

func (h *OpportunityEventHandler) OnUpdateNextCycleDate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnUpdateNextCycleDate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

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

	contractDbNode, err := h.repositories.ContractRepository.GetContractByOpportunityId(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
	}
	if contractDbNode != nil {
		contractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityLikelihood(ctx, eventData.Tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractEntity.Id, err.Error())
		}
	}

	h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)

	return nil
}

func (h *OpportunityEventHandler) sendEventToUpdateOrganizationRenewalSummary(ctx context.Context, tenant, opportunityId string, span opentracing.Span) {
	organizationDbNode, err := h.repositories.OrganizationRepository.GetOrganizationByOpportunityId(ctx, tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
		return
	}
	if organizationDbNode == nil {
		return
	}
	organization := graph_db.MapDbNodeToOrganizationEntity(*organizationDbNode)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = h.grpcClients.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.OrganizationIdGrpcRequest{
		Tenant:         tenant,
		OrganizationId: organization.ID,
		AppSource:      constants.AppSourceEventProcessingPlatform,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("RefreshRenewalSummary failed: %v", err.Error())
	}
}

func (h *OpportunityEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)

	opportunityDbNode, err := h.repositories.OpportunityRepository.GetOpportunityById(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting opportunity %s: %s", opportunityId, err.Error())
		return err
	}
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	amountChanged := ((opportunity.Amount != eventData.Amount) && eventData.UpdateAmount()) ||
		((opportunity.MaxAmount != eventData.MaxAmount) && eventData.UpdateMaxAmount())

	err = h.repositories.OpportunityRepository.Update(ctx, eventData.Tenant, opportunityId, eventData)
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

	// if amount changed, recalculate organization combined ARR forecast
	if amountChanged {
		organizationDbNode, err := h.repositories.OrganizationRepository.GetOrganizationByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
		if organizationDbNode == nil {
			return nil
		}
		organization := graph_db.MapDbNodeToOrganizationEntity(*organizationDbNode)

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = h.grpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: organization.ID,
			AppSource:      constants.AppSourceEventProcessingPlatform,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshArr failed: %v", err.Error())
		}
	}

	return nil
}

func (h *OpportunityEventHandler) OnUpdateRenewal(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnUpdateRenewal")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityUpdateRenewalEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	opportunityDbNode, err := h.repositories.OpportunityRepository.GetOpportunityById(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting opportunity %s: %s", opportunityId, err.Error())
		return err
	}
	opportunity := graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode)
	amountChanged := eventData.UpdateAmount() && opportunity.Amount != eventData.Amount
	likelihoodChanged := eventData.UpdateRenewalLikelihood() && opportunity.RenewalDetails.RenewalLikelihood != eventData.RenewalLikelihood
	setUpdatedByUserId := (amountChanged || likelihoodChanged) && eventData.UpdatedByUserId != ""
	if eventData.OwnerUserId != "" {
		err = h.repositories.OpportunityRepository.ReplaceOwner(ctx, eventData.Tenant, opportunityId, eventData.OwnerUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while replacing owner of opportunity %s: %s", opportunityId, err.Error())
			return err
		}
	}
	err = h.repositories.OpportunityRepository.UpdateRenewal(ctx, eventData.Tenant, opportunityId, eventData, setUpdatedByUserId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	if likelihoodChanged {
		h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)
	}
	// update renewal ARR if likelihood changed but amount didn't
	if likelihoodChanged && !amountChanged {
		contractDbNode, err := h.repositories.ContractRepository.GetContractByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
		if contractDbNode == nil {
			return nil
		}
		contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
	} else if amountChanged {
		// if amount changed, recalculate organization combined ARR forecast
		organizationDbNode, err := h.repositories.OrganizationRepository.GetOrganizationByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
		if organizationDbNode == nil {
			return nil
		}
		organization := graph_db.MapDbNodeToOrganizationEntity(*organizationDbNode)

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = h.grpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: organization.ID,
			AppSource:      constants.AppSourceEventProcessingPlatform,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshArr failed: %v", err.Error())
		}
	}

	return nil
}

func (h *OpportunityEventHandler) OnCloseWin(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnCloseWin")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityCloseWinEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.OpportunityRepository.CloseWin(ctx, eventData.Tenant, opportunityId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while closing opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	return nil
}

func (h *OpportunityEventHandler) OnCloseLoose(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityEventHandler.OnCloseLoose")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.OpportunityCloseLooseEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	opportunityId := aggregate.GetOpportunityObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.repositories.OpportunityRepository.CloseLoose(ctx, eventData.Tenant, opportunityId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while closing opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)

	return nil
}
