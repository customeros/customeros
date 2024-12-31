package graph

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"time"
)

type ServiceLineItemEventHandler struct {
	log         logger.Logger
	services    *service.Services
	grpcClients *grpc_client.Clients
}

func NewServiceLineItemEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients) *ServiceLineItemEventHandler {
	return &ServiceLineItemEventHandler{
		log:         log,
		services:    services,
		grpcClients: grpcClients,
	}
}

type userMetadata struct {
	UserId string `json:"user-id"`
}

type SLIActionMetadata struct {
	UserName         string     `json:"user-name"`
	ServiceName      string     `json:"service-name"`
	Price            float64    `json:"price"`
	Currency         string     `json:"currency"`
	Comment          string     `json:"comment"`
	ReasonForChange  string     `json:"reasonForChange"`
	StartedAt        *time.Time `json:"startedAt,omitempty"`
	BilledType       string     `json:"billedType"`
	Quantity         int64      `json:"quantity"`
	PreviousPrice    float64    `json:"previousPrice"`
	PreviousQuantity int64      `json:"previousQuantity"`
}

func (h *ServiceLineItemEventHandler) OnDeleteV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnDeleteV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)
	var serviceLineItemName string
	var contractName string

	var eventData event.ServiceLineItemDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)
	serviceLineItemDbNode, err := h.services.CommonServices.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	serviceLineItemEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(serviceLineItemDbNode)
	if serviceLineItemEntity.Name != "" {
		serviceLineItemName = serviceLineItemEntity.Name
	} else {
		serviceLineItemName = "Unnamed service"
	}

	// get user
	username := ""
	usrMetadata := userMetadata{}
	if err := json.Unmarshal(evt.Metadata, &usrMetadata); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Unmarshal")
	} else {
		if usrMetadata.UserId != "" {
			user, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetUserById(ctx, eventData.Tenant, usrMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to get user for service line item %s with userid %s", serviceLineItemId, usrMetadata.UserId)
			}
			userEntity := *neo4jmapper.MapDbNodeToUserEntity(user)
			username = userEntity.GetFullName()
		}
	}

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
	if contract.Name != "" {
		contractName = contract.Name
	} else {
		contractName = "Unnamed contract"
	}

	err = h.services.CommonServices.Neo4jRepositories.ServiceLineItemWriteRepository.Delete(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}
	err = h.services.CommonServices.Neo4jRepositories.ServiceLineItemWriteRepository.AdjustEndDates(ctx, nil, eventData.Tenant, serviceLineItemEntity.ParentID)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while adjusting end dates for service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	if contractDbNode != nil {
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contract.Id, err.Error())
			return nil
		}
		// Update contract LTV
		contractHandler.UpdateContractLtv(ctx, eventData.Tenant, contract.Id)
	}
	metadata, err := utils.ToJson(SLIActionMetadata{
		UserName:    username,
		ServiceName: serviceLineItemName,
		Comment:     "service line item removed is " + serviceLineItemName + " from " + contractName + " by " + username,
	})
	message := username + " removed " + serviceLineItemName + " from " + contractName

	_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.Create(ctx, eventData.Tenant, contract.Id, commonmodel.CONTRACT, neo4jenum.ActionServiceLineItemRemoved, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed remove service line item action for contract %s: %s", contract.Id, err.Error())
	}

	h.services.CommonServices.RabbitMQService.PublishEventCompleted(ctx, eventData.Tenant, serviceLineItemId, commonmodel.SERVICE_LINE_ITEM, utils.NewEventCompletedDetails().WithDelete())

	return nil
}

func (h *ServiceLineItemEventHandler) OnClose(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnClose")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.ServiceLineItemCloseEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.ServiceLineItemWriteRepository.Close(ctx, eventData.Tenant, serviceLineItemId, eventData.EndedAt, eventData.IsCanceled)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while closing service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	if contractDbNode != nil {
		contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contract.Id, err.Error())
			return nil
		}
		// Update contract LTV
		contractHandler.UpdateContractLtv(ctx, eventData.Tenant, contract.Id)
	}

	h.services.CommonServices.RabbitMQService.PublishEventCompleted(ctx, eventData.Tenant, serviceLineItemId, commonmodel.SERVICE_LINE_ITEM, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func getBillingCycleNamingConvention(billedType string) string {
	switch billedType {
	case model.AnnuallyBilled.String():
		return "year"
	case model.QuarterlyBilled.String():
		return "quarter"
	case model.MonthlyBilled.String():
		return "month"
	default:
		return ""
	}
}
