package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

	h.services.CommonServices.RabbitMQService.PublishEventCompleted(ctx, eventData.Tenant, contractId, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}
