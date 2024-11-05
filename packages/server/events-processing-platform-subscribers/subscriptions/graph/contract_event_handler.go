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
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"math"
	"time"
)

type ActionStatusMetadata struct {
	Status       string `json:"status"`
	ContractName string `json:"contract-name"`
	Comment      string `json:"comment"`
}

type ContractEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewContractEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *ContractEventHandler {
	return &ContractEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
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
		OrganizationId:         eventData.OrganizationId,
		Name:                   eventData.Name,
		ContractUrl:            eventData.ContractUrl,
		CreatedByUserId:        eventData.CreatedByUserId,
		ServiceStartedAt:       eventData.ServiceStartedAt,
		SignedAt:               eventData.SignedAt,
		LengthInMonths:         eventData.LengthInMonths,
		Status:                 eventData.Status,
		CreatedAt:              eventData.CreatedAt,
		BillingCycleInMonths:   eventData.BillingCycleInMonths,
		Currency:               neo4jenum.DecodeCurrency(eventData.Currency),
		InvoicingStartDate:     eventData.InvoicingStartDate,
		InvoicingEnabled:       eventData.InvoicingEnabled,
		PayOnline:              eventData.PayOnline,
		PayAutomatically:       eventData.PayAutomatically,
		AutoRenew:              eventData.AutoRenew,
		CanPayWithCard:         eventData.CanPayWithCard,
		CanPayWithDirectDebit:  eventData.CanPayWithDirectDebit,
		CanPayWithBankTransfer: eventData.CanPayWithBankTransfer,
		Check:                  eventData.Check,
		DueDays:                eventData.DueDays,
		Country:                eventData.Country,
		Approved:               eventData.Approved,
		SourceFields: neo4jmodel.SourceFields{
			Source:        helper.GetSource(eventData.Source.Source),
			AppSource:     helper.GetAppSource(eventData.Source.AppSource),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.Source.Source),
		},
	}
	err := h.services.CommonServices.Neo4jRepositories.ContractWriteRepository.CreateForOrganizationOld(ctx, eventData.Tenant, contractId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving contract %s: %s", contractId, err.Error())
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
		err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, contractId, model.NodeLabelContract, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while linking contract %s with external system %s: %s", contractId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	_, _, err = h.updateStatus(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
		return err
	}

	if eventData.LengthInMonths > 0 {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return h.grpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
				Tenant:     eventData.Tenant,
				ContractId: contractId,
				SourceFields: &commonpb.SourceFields{
					Source:    eventData.Source.Source,
					AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CreateRenewalOpportunity failed: %s", err.Error())
		}
	}

	h.startOnboardingIfEligible(ctx, eventData.Tenant, contractId, span)

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTRACT.String(), contractId, h.grpcClients, utils.NewEventCompletedDetails().WithCreate())

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

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	beforeUpdateContractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	data := neo4jrepository.ContractUpdateFields{
		Name:                         eventData.Name,
		ContractUrl:                  eventData.ContractUrl,
		ServiceStartedAt:             eventData.ServiceStartedAt,
		Source:                       helper.GetSource(eventData.Source),
		LengthInMonths:               eventData.LengthInMonths,
		SignedAt:                     eventData.SignedAt,
		EndedAt:                      eventData.EndedAt,
		BillingCycleInMonths:         eventData.BillingCycleInMonths,
		Currency:                     neo4jenum.DecodeCurrency(eventData.Currency),
		InvoicingStartDate:           eventData.InvoicingStartDate,
		AddressLine1:                 eventData.AddressLine1,
		AddressLine2:                 eventData.AddressLine2,
		Locality:                     eventData.Locality,
		Country:                      eventData.Country,
		Region:                       eventData.Region,
		Zip:                          eventData.Zip,
		OrganizationLegalName:        eventData.OrganizationLegalName,
		InvoiceEmail:                 eventData.InvoiceEmail,
		InvoiceEmailCC:               eventData.InvoiceEmailCC,
		InvoiceEmailBCC:              eventData.InvoiceEmailBCC,
		InvoiceNote:                  eventData.InvoiceNote,
		NextInvoiceDate:              eventData.NextInvoiceDate,
		InvoicingEnabled:             eventData.InvoicingEnabled,
		PayOnline:                    eventData.PayOnline,
		PayAutomatically:             eventData.PayAutomatically,
		CanPayWithCard:               eventData.CanPayWithCard,
		CanPayWithDirectDebit:        eventData.CanPayWithDirectDebit,
		CanPayWithBankTransfer:       eventData.CanPayWithBankTransfer,
		AutoRenew:                    eventData.AutoRenew,
		Check:                        eventData.Check,
		DueDays:                      eventData.DueDays,
		Approved:                     eventData.Approved,
		UpdateName:                   eventData.UpdateName(),
		UpdateContractUrl:            eventData.UpdateContractUrl(),
		UpdateServiceStartedAt:       eventData.UpdateServiceStartedAt(),
		UpdateSignedAt:               eventData.UpdateSignedAt(),
		UpdateEndedAt:                eventData.UpdateEndedAt(),
		UpdateInvoicingStartDate:     eventData.UpdateInvoicingStartDate(),
		UpdateLengthInMonths:         eventData.UpdateLengthInMonths(),
		UpdateBillingCycleInMonths:   eventData.UpdateBillingCycleInMonths(),
		UpdateCurrency:               eventData.UpdateCurrency(),
		UpdateAddressLine1:           eventData.UpdateAddressLine1(),
		UpdateAddressLine2:           eventData.UpdateAddressLine2(),
		UpdateLocality:               eventData.UpdateLocality(),
		UpdateCountry:                eventData.UpdateCountry(),
		UpdateRegion:                 eventData.UpdateRegion(),
		UpdateZip:                    eventData.UpdateZip(),
		UpdateOrganizationLegalName:  eventData.UpdateOrganizationLegalName(),
		UpdateInvoiceEmail:           eventData.UpdateInvoiceEmail(),
		UpdateInvoiceEmailCC:         eventData.UpdateInvoiceEmailCC(),
		UpdateInvoiceEmailBCC:        eventData.UpdateInvoiceEmailBCC(),
		UpdateInvoiceNote:            eventData.UpdateInvoiceNote(),
		UpdateNextInvoiceDate:        eventData.UpdateNextInvoiceDate(),
		UpdateCanPayWithCard:         eventData.UpdateCanPayWithCard(),
		UpdateCanPayWithDirectDebit:  eventData.UpdateCanPayWithDirectDebit(),
		UpdateCanPayWithBankTransfer: eventData.UpdateCanPayWithBankTransfer(),
		UpdateInvoicingEnabled:       eventData.UpdateInvoicingEnabled(),
		UpdatePayOnline:              eventData.UpdatePayOnline(),
		UpdatePayAutomatically:       eventData.UpdatePayAutomatically(),
		UpdateAutoRenew:              eventData.UpdateAutoRenew(),
		UpdateCheck:                  eventData.UpdateCheck(),
		UpdateDueDays:                eventData.UpdateDueDays(),
		UpdateApproved:               eventData.UpdateApproved(),
	}
	err = h.services.CommonServices.Neo4jRepositories.ContractWriteRepository.UpdateContractOld(ctx, eventData.Tenant, contractId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s: %s", contractId, err.Error())
		return err
	}
	_, statusChanged, err := h.updateStatus(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
		return err
	}

	if statusChanged {
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateOrganizationRelationship(ctx, eventData.Tenant, contractId, statusChanged)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while updating organization relationship for contract %s: %s", contractId, err.Error())
		}
	}

	contractDbNode, err = h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	afterUpdateContractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if eventData.ExternalSystem.Available() {
		externalSystemData := neo4jmodel.ExternalSystem{
			ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
			ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
			ExternalId:       eventData.ExternalSystem.ExternalId,
			ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
			ExternalSource:   eventData.ExternalSystem.ExternalSource,
			SyncDate:         eventData.ExternalSystem.SyncDate,
		}
		err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, eventData.Tenant, contractId, model.NodeLabelContract, externalSystemData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while link contract %s with external system %s: %s", contractId, eventData.ExternalSystem.ExternalSystemId, err.Error())
			return err
		}
	}

	if beforeUpdateContractEntity.LengthInMonths > 0 && afterUpdateContractEntity.LengthInMonths == 0 {
		err = h.services.CommonServices.Neo4jRepositories.ContractWriteRepository.SuspendActiveRenewalOpportunity(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while suspending renewal opportunity for contract %s: %s", contractId, err.Error())
		}
		organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, eventData.Tenant, contractId)
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
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return h.grpcClients.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.RefreshRenewalSummaryGrpcRequest{
				Tenant:         eventData.Tenant,
				OrganizationId: organization.ID,
				AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RefreshRenewalSummary failed: %v", err.Error())
		}
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
	} else {
		if beforeUpdateContractEntity.LengthInMonths == 0 && afterUpdateContractEntity.LengthInMonths > 0 {
			err = h.services.CommonServices.Neo4jRepositories.ContractWriteRepository.ActivateSuspendedRenewalOpportunity(ctx, eventData.Tenant, contractId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Error while activating renewal opportunity for contract %s: %s", contractId, err.Error())
			}
		}
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityRenewDateAndArr(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
		}
	}

	if beforeUpdateContractEntity.ContractStatus != afterUpdateContractEntity.ContractStatus {
		h.createActionForStatusChange(ctx, eventData.Tenant, contractId, string(afterUpdateContractEntity.ContractStatus), afterUpdateContractEntity.Name)
	}

	contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
	err = contractHandler.UpdateActiveRenewalOpportunityLikelihood(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
	}
	contractHandler.UpdateContractLtv(ctx, eventData.Tenant, contractId)

	h.startOnboardingIfEligible(ctx, eventData.Tenant, contractId, span)

	if eventData.AppSource != constants.AppSourceCustomerOsApi {
		utils.EventCompleted(ctx, eventData.Tenant, model.CONTRACT.String(), contractId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
	}
	return nil
}

func (h *ContractEventHandler) OnRolloutRenewalOpportunity(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnRolloutRenewalOpportunity")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.RolloutRenewalOpportunityEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	contractId := aggregate.GetContractObjectID(evt.GetAggregateID(), eventData.Tenant)

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if contractEntity.LengthInMonths <= 0 {
		return nil
	}

	currentRenewalOpportunityDbNode, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	if currentRenewalOpportunityDbNode != nil {
		currentOpportunity := neo4jmapper.MapDbNodeToOpportunityEntity(currentRenewalOpportunityDbNode)

		err := h.services.CommonServices.OpportunityService.CloseWon(ctx, nil, eventData.Tenant, currentOpportunity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("CloseWinOpportunity failed: %s", err.Error())
			return err
		}
	}

	// Update contract LTV
	contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
	contractHandler.UpdateContractLtv(ctx, eventData.Tenant, contractId)

	// Add action in timeline
	status := "Renewed"
	metadata, err := utils.ToJson(ActionStatusMetadata{
		Status: status,
	})
	message := contractEntity.Name + " renewed"

	_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.Create(ctx, eventData.Tenant, contractId, model.CONTRACT, neo4jenum.ActionContractRenewed, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating renewed action for contract %s: %s", contractId, err.Error())
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTRACT.String(), contractId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *ContractEventHandler) createActionForStatusChange(ctx context.Context, tenant, contractId, status, contractName string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.createActionForStatusChange")
	defer span.Finish()
	var name string
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId), log.String("status", status), log.String("contractName", contractName))

	if contractName != "" {
		name = contractName
	} else {
		name = "Unnamed contract"
	}
	actionStatusMetadata := ActionStatusMetadata{
		Status:       status,
		ContractName: name,
		Comment:      name + " is now " + status,
	}
	message := ""

	switch status {
	case string(neo4jenum.ContractStatusLive):
		message = contractName + " is now live"
		actionStatusMetadata.Comment = contractName + " is now live"
	case string(neo4jenum.ContractStatusEnded):
		message = contractName + " has ended"
		actionStatusMetadata.Comment = contractName + " has ended"
	case string(neo4jenum.ContractStatusOutOfContract):
		message = contractName + " is now out of contract"
		actionStatusMetadata.Comment = contractName + " is now out of contract"
	}
	metadata, err := utils.ToJson(actionStatusMetadata)
	_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, contractId, model.CONTRACT, neo4jenum.ActionContractStatusUpdated, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating status update action for contract %s: %s", contractId, err.Error())
	}
}

func (h *ContractEventHandler) startOnboardingIfEligible(ctx context.Context, tenant, contractId string, span opentracing.Span) {

	// TODO temporary not eligible for all contracts
	return

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	if contractDbNode == nil {
		return
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if contractEntity.IsEligibleToStartOnboarding() {
		organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractEntity.Id)
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
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return h.grpcClients.OrganizationClient.UpdateOnboardingStatus(ctx, &organizationpb.UpdateOnboardingStatusGrpcRequest{
				Tenant:             tenant,
				OrganizationId:     organization.ID,
				CausedByContractId: contractEntity.Id,
				OnboardingStatus:   organizationpb.OnboardingStatus_ONBOARDING_STATUS_NOT_STARTED,
				AppSource:          constants.AppSourceEventProcessingPlatformSubscribers,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("UpdateOnboardingStatus gRPC request failed: %v", err.Error())
		}
	}
}

func (h *ContractEventHandler) OnDeleteV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnDeleteV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContractDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contractId := aggregate.GetContractObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	// fetch organization of the contract
	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, eventData.Tenant, contractId)
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

	err = h.services.CommonServices.Neo4jRepositories.ContractWriteRepository.SoftDelete(ctx, eventData.Tenant, contractId, eventData.UpdatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting contract %s: %s", contractId, err.Error())
		return err
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.RefreshRenewalSummaryGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: organization.ID,
			AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: organization.ID,
			AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})

	err = h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.DeletePreviewCycleInvoices(ctx, eventData.Tenant, contractId, "")
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting preview invoice for contract %s: %s", contractId, err.Error())
		return err
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTRACT.String(), contractId, h.grpcClients, utils.NewEventCompletedDetails().WithDelete())

	return nil
}

func (h *ContractEventHandler) updateStatus(ctx context.Context, tenant, contractId string) (string, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.updateStatus")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return "", false, err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	status, err := h.deriveContractStatus(ctx, tenant, *contractEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deriving contract %s status: %s", contractId, err.Error())
		return "", false, err
	}
	statusChanged := contractEntity.ContractStatus.String() != status

	err = h.services.CommonServices.Neo4jRepositories.ContractWriteRepository.UpdateStatus(ctx, tenant, contractId, status)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
		return "", false, err
	}

	return status, statusChanged, nil
}

func (h *ContractEventHandler) OnRefreshStatus(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnRefreshStatus")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContractUpdateStatusEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contractId := aggregate.GetContractObjectID(evt.GetAggregateID(), eventData.Tenant)

	status, statusChanged, err := h.updateStatus(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
		return err
	}

	if statusChanged {
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateOrganizationRelationship(ctx, eventData.Tenant, contractId, statusChanged)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while updating organization relationship for contract %s: %s", contractId, err.Error())
		}
		contractHandler.UpdateContractLtv(ctx, eventData.Tenant, contractId)
	}

	if status == neo4jenum.ContractStatusEnded.String() {
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityRenewDateAndArr(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating contract's {%s} renewal date: %s", contractId, err.Error())
		}

		err := h.services.CommonServices.Neo4jRepositories.InvoiceWriteRepository.DeletePreviewCycleInvoices(ctx, eventData.Tenant, contractId, "")
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while deleting preview invoice for contract %s: %s", contractId, err.Error())
		}
	}

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if statusChanged {
		h.createActionForStatusChange(ctx, eventData.Tenant, contractId, status, contractEntity.Name)
	}

	h.startOnboardingIfEligible(ctx, eventData.Tenant, contractId, span)

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTRACT.String(), contractId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *ContractEventHandler) deriveContractStatus(ctx context.Context, tenant string, contractEntity neo4jentity.ContractEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.deriveContractStatus")
	defer span.Finish()

	now := utils.Now()

	// If endedAt is not nil and is in the past, the contract is considered Ended.
	if contractEntity.IsEnded() {
		span.LogFields(log.String("result.status", neo4jenum.ContractStatusEnded.String()))
		return neo4jenum.ContractStatusEnded.String(), nil
	}

	// check if contract is draft
	if !contractEntity.Approved {
		span.LogFields(log.String("result.status", neo4jenum.ContractStatusDraft.String()))
		return neo4jenum.ContractStatusDraft.String(), nil
	}

	// Check contract is scheduled
	if contractEntity.ServiceStartedAt == nil || contractEntity.ServiceStartedAt.After(now) {
		span.LogFields(log.String("result.status", neo4jenum.ContractStatusScheduled.String()))
		return neo4jenum.ContractStatusScheduled.String(), nil
	}

	// Check if contract is out of contract
	if !contractEntity.AutoRenew {
		// fetch active renewal opportunity for the contract
		opportunityDbNode, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if opportunityDbNode != nil {
			opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
			if opportunityEntity.RenewalDetails.RenewedAt != nil && opportunityEntity.RenewalDetails.RenewedAt.Before(now) {
				span.LogFields(log.String("result.status", neo4jenum.ContractStatusLive.String()))
				return neo4jenum.ContractStatusOutOfContract.String(), nil
			}
		}
	}

	// Otherwise, the contract is considered Live.
	span.LogFields(log.String("result.status", neo4jenum.ContractStatusLive.String()))
	return neo4jenum.ContractStatusLive.String(), nil
}

func (h *ContractEventHandler) OnRefreshLtv(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractEventHandler.OnRefreshLtv")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ContractRefreshLtvEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	contractId := aggregate.GetContractObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	ltv := 0.0
	recalculateContractLtv := true
	if !(contractEntity.ContractStatus == neo4jenum.ContractStatusLive ||
		contractEntity.ContractStatus == neo4jenum.ContractStatusOutOfContract ||
		contractEntity.ContractStatus == neo4jenum.ContractStatusEnded) {
		span.LogFields(log.String("result", fmt.Sprintf("contract status %s is not eligible for LTV calculation", contractEntity.ContractStatus)))
		recalculateContractLtv = false
	}

	if recalculateContractLtv {
		sliDbNodes, err := h.services.CommonServices.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemsForContract(ctx, eventData.Tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		var sliEntities []*neo4jentity.ServiceLineItemEntity
		for _, sliDbNode := range sliDbNodes {
			sliEntities = append(sliEntities, neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode))
		}

		// Calculate LTV

		// Step 1 calculate one times
		for _, sliEntity := range sliEntities {
			if sliEntity.IsOneTime() {
				sliLtv := float64(sliEntity.Quantity) * sliEntity.Price
				ltv += sliLtv
				span.LogFields(log.String("result.sli - ltv", fmt.Sprintf("%s - %f", sliEntity.ID, utils.TruncateFloat64(sliLtv, 2))))
			}
		}

		defaultEndDate := utils.Today()
		if contractEntity.IsEnded() && contractEntity.EndedAt != nil {
			defaultEndDate = *contractEntity.EndedAt
		}
		// Step 2 calculate recurring
		for _, sliEntity := range sliEntities {
			if sliEntity.IsRecurrent() {
				endDate := defaultEndDate
				if sliEntity.EndedAt != nil && sliEntity.EndedAt.Before(defaultEndDate) {
					endDate = *sliEntity.EndedAt
				}
				duration := calculateDuration(sliEntity.StartedAt, endDate, sliEntity.Billed)
				sliLtv := float64(sliEntity.Quantity) * sliEntity.Price * duration
				ltv += sliLtv
				span.LogFields(log.String("result.sli - ltv", fmt.Sprintf("%s - %f", sliEntity.ID, utils.TruncateFloat64(sliLtv, 2))))
			}
		}
	}

	truncatedLtv := utils.TruncateFloat64(ltv, 2)
	err = h.services.CommonServices.Neo4jRepositories.ContractWriteRepository.SetLtv(ctx, eventData.Tenant, contractId, truncatedLtv)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating contract %s ltv: %s", contractId, err.Error())
		return err
	}

	// get organization for contract
	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, eventData.Tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
		return nil
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	// request organization ltv refresh
	if organizationEntity.ID != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return h.grpcClients.OrganizationClient.RefreshDerivedData(ctx, &organizationpb.RefreshDerivedDataGrpcRequest{
				Tenant:         eventData.Tenant,
				OrganizationId: organizationEntity.ID,
				AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	utils.EventCompleted(ctx, eventData.Tenant, model.CONTRACT.String(), contractId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func calculateDuration(startedAt, endedAt time.Time, billed neo4jenum.BilledType) float64 {
	if startedAt.After(endedAt) {
		return float64(0)
	}
	durationDays := math.Abs(float64(daysBetween(startedAt, endedAt)))

	switch billed {
	case neo4jenum.BilledTypeMonthly:
		return durationDays / 30
	case neo4jenum.BilledTypeQuarterly:
		return durationDays / 90
	case neo4jenum.BilledTypeAnnually:
		return durationDays / 365
	default:
		return 0
	}
}

func daysBetween(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours() / 24)
}
