package graph

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

type ActionLikelihoodMetadata struct {
	Likelihood string `json:"likelihood"`
	Reason     string `json:"reason"`
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
	data := neo4jrepository.OpportunityCreateFields{
		OrganizationId: eventData.OrganizationId,
		CreatedAt:      eventData.CreatedAt,
		UpdatedAt:      eventData.UpdatedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source.Source),
			SourceOfTruth: helper.GetSource(eventData.Source.Source),
			AppSource:     helper.GetAppSource(eventData.Source.AppSource),
		},
		Name:              eventData.Name,
		Amount:            eventData.Amount,
		InternalType:      eventData.InternalType,
		ExternalType:      eventData.ExternalType,
		InternalStage:     eventData.InternalStage,
		ExternalStage:     eventData.ExternalStage,
		EstimatedClosedAt: eventData.EstimatedClosedAt,
		GeneralNotes:      eventData.GeneralNotes,
		NextSteps:         eventData.NextSteps,
		CreatedByUserId:   eventData.CreatedByUserId,
	}
	err := h.repositories.Neo4jRepositories.OpportunityWriteRepository.CreateForOrganization(ctx, eventData.Tenant, opportunityId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	if eventData.OwnerUserId != "" {
		err = h.repositories.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, eventData.Tenant, opportunityId, eventData.OwnerUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while replacing owner of opportunity %s: %s", opportunityId, err.Error())
			return err
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
		err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, opportunityId, neo4jutil.NodeLabelOpportunity, externalSystemData)
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

	// check if active renewal opportunity already exists for this contract
	opportunotyDbNode, err := h.repositories.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting renewal opportunity for contract %s: %s", eventData.ContractId, err.Error())
		return nil
	}
	if opportunotyDbNode != nil {
		opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunotyDbNode)
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
		UpdatedAt:  eventData.UpdatedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source.Source),
			SourceOfTruth: helper.GetSource(eventData.Source.Source),
			AppSource:     helper.GetAppSource(eventData.Source.AppSource),
		},
		InternalType:      eventData.InternalType,
		InternalStage:     eventData.InternalStage,
		RenewalLikelihood: eventData.RenewalLikelihood,
		RenewalApproved:   eventData.RenewalApproved,
	}
	err = h.repositories.Neo4jRepositories.OpportunityWriteRepository.CreateRenewal(ctx, eventData.Tenant, opportunityId, data)
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
	err := h.repositories.Neo4jRepositories.OpportunityWriteRepository.UpdateNextRenewalDate(ctx, eventData.Tenant, opportunityId, eventData.UpdatedAt, eventData.RenewedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating next cycle date for opportunity %s: %s", opportunityId, err.Error())
	}

	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
	}
	if contractDbNode != nil {
		contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
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

	return nil
}

func (h *OpportunityEventHandler) sendEventToUpdateOrganizationRenewalSummary(ctx context.Context, tenant, opportunityId string, span opentracing.Span) {
	organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, tenant, opportunityId)
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

	opportunityDbNode, err := h.repositories.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting opportunity %s: %s", opportunityId, err.Error())
		return err
	}
	opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
	amountChanged := ((opportunity.Amount != eventData.Amount) && eventData.UpdateAmount()) ||
		((opportunity.MaxAmount != eventData.MaxAmount) && eventData.UpdateMaxAmount())

	data := neo4jrepository.OpportunityUpdateFields{
		UpdatedAt:       eventData.UpdatedAt,
		Source:          eventData.Source,
		Name:            eventData.Name,
		Amount:          eventData.Amount,
		MaxAmount:       eventData.MaxAmount,
		UpdateName:      eventData.UpdateName(),
		UpdateAmount:    eventData.UpdateAmount(),
		UpdateMaxAmount: eventData.UpdateMaxAmount(),
	}
	err = h.repositories.Neo4jRepositories.OpportunityWriteRepository.Update(ctx, eventData.Tenant, opportunityId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving opportunity %s: %s", opportunityId, err.Error())
		return err
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
		err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, opportunityId, neo4jutil.NodeLabelOpportunity, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while linking opportunity %s with external system %s: %s", opportunityId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	// if amount changed, recalculate organization combined ARR forecast
	if amountChanged {
		organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
		if organizationDbNode == nil {
			return nil
		}
		organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return h.grpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
				Tenant:         eventData.Tenant,
				OrganizationId: organization.ID,
				AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
			})
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
	opportunityDbNode, err := h.repositories.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, eventData.Tenant, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting opportunity %s: %s", opportunityId, err.Error())
		return err
	}
	opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
	amountChanged := eventData.UpdateAmount() && opportunity.Amount != eventData.Amount
	likelihoodChanged := eventData.UpdateRenewalLikelihood() && opportunity.RenewalDetails.RenewalLikelihood.String() != eventData.RenewalLikelihood
	setUpdatedByUserId := (amountChanged || likelihoodChanged) && eventData.UpdatedByUserId != ""
	if eventData.OwnerUserId != "" {
		err = h.repositories.Neo4jRepositories.OpportunityWriteRepository.ReplaceOwner(ctx, eventData.Tenant, opportunityId, eventData.OwnerUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while replacing owner of opportunity %s: %s", opportunityId, err.Error())
			return err
		}
	}
	data := neo4jrepository.RenewalOpportunityUpdateFields{
		UpdatedAt:               eventData.UpdatedAt,
		Source:                  helper.GetSource(eventData.Source),
		UpdatedByUserId:         eventData.UpdatedByUserId,
		SetUpdatedByUserId:      setUpdatedByUserId,
		Comments:                eventData.Comments,
		Amount:                  eventData.Amount,
		RenewalLikelihood:       eventData.RenewalLikelihood,
		RenewalApproved:         eventData.RenewalApproved,
		UpdateComments:          eventData.UpdateComments(),
		UpdateAmount:            eventData.UpdateAmount(),
		UpdateRenewalLikelihood: eventData.UpdateRenewalLikelihood(),
		UpdateRenewalApproved:   eventData.UpdateRenewalApproved(),
	}
	err = h.repositories.Neo4jRepositories.OpportunityWriteRepository.UpdateRenewal(ctx, eventData.Tenant, opportunityId, data)
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
		contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting contract for opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
		if contractDbNode == nil {
			return nil
		}
		contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
	} else if amountChanged {
		// if amount changed, recalculate organization combined ARR forecast
		organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByOpportunityId(ctx, eventData.Tenant, opportunityId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting organization for opportunity %s: %s", opportunityId, err.Error())
			return nil
		}
		if organizationDbNode == nil {
			return nil
		}
		organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return h.grpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
				Tenant:         eventData.Tenant,
				OrganizationId: organization.ID,
				AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshArr failed: %v", err.Error())
		}
	}

	// prepare action for likelihood change
	if likelihoodChanged {
		contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractByOpportunityId(ctx, eventData.Tenant, opportunityId)
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
	err := h.repositories.Neo4jRepositories.OpportunityWriteRepository.CloseWin(ctx, eventData.Tenant, opportunityId, eventData.UpdatedAt, eventData.ClosedAt)
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
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	err := h.repositories.Neo4jRepositories.OpportunityWriteRepository.CloseLoose(ctx, eventData.Tenant, opportunityId, eventData.UpdatedAt, eventData.ClosedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while closing opportunity %s: %s", opportunityId, err.Error())
		return err
	}

	h.sendEventToUpdateOrganizationRenewalSummary(ctx, eventData.Tenant, opportunityId, span)

	return nil
}

func (h *OpportunityEventHandler) saveLikelihoodChangeAction(ctx context.Context, contractId string, eventData event.OpportunityUpdateRenewalEvent, span opentracing.Span) error {
	metadata, err := utils.ToJson(ActionLikelihoodMetadata{
		Reason:     eventData.Comments,
		Likelihood: eventData.RenewalLikelihood,
	})
	userName := ""
	if eventData.UpdatedByUserId != "" {
		userDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, eventData.Tenant, eventData.UpdatedByUserId)
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
	_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, neo4jenum.CONTRACT, neo4jenum.ActionRenewalLikelihoodUpdated, message, metadata, eventData.UpdatedAt, constants.AppSourceEventProcessingPlatformSubscribers, extraActionProperties)
	return err
}
