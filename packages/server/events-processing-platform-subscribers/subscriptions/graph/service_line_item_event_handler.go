package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strconv"
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

func (h *ServiceLineItemEventHandler) OnUpdateV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnUpdateV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)
	var contractId string
	var name string
	var message string
	var eventData event.ServiceLineItemUpdateEvent
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
	//we will use the following booleans below to check if the price, quantity, billed type has changed
	priceChanged := serviceLineItemEntity.Price != eventData.Price
	quantityChanged := serviceLineItemEntity.Quantity != eventData.Quantity

	data := neo4jrepository.ServiceLineItemUpdateFields{
		Price:     eventData.Price,
		Quantity:  eventData.Quantity,
		Billed:    eventData.Billed,
		Comments:  eventData.Comments,
		Name:      eventData.Name,
		Source:    helper.GetSource(eventData.Source.Source),
		VatRate:   eventData.VatRate,
		StartedAt: eventData.StartedAt,
	}
	err = h.services.CommonServices.Neo4jRepositories.ServiceLineItemWriteRepository.Update(ctx, eventData.Tenant, serviceLineItemId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	err = h.services.CommonServices.Neo4jRepositories.ServiceLineItemWriteRepository.AdjustEndDates(ctx, nil, eventData.Tenant, serviceLineItemEntity.ParentID)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while adjusting end dates for service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	contractDbNode, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
	if contractDbNode != nil {
		contractId = contractEntity.Id
		contractHandler := contracthandler.NewContractHandler(h.log, h.services, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractEntity.Id, err.Error())
			return nil
		}
		// Update contract LTV
		contractHandler.UpdateContractLtv(ctx, eventData.Tenant, contractEntity.Id)
	}
	//check to make sure the name displays correctly in the action message
	if eventData.Name == "" {
		name = serviceLineItemEntity.Name
	} else {
		name = eventData.Name
	}
	if name == "" {
		name = "Unnamed service"
	}
	// get user
	username := ""
	usrMetadata := userMetadata{}
	if err = json.Unmarshal(evt.Metadata, &usrMetadata); err != nil {
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
	actionPriceMetadata := SLIActionMetadata{
		UserName:        username,
		ServiceName:     serviceLineItemEntity.Name,
		Price:           eventData.Price,
		PreviousPrice:   serviceLineItemEntity.Price,
		BilledType:      serviceLineItemEntity.Billed.String(),
		Quantity:        serviceLineItemEntity.Quantity,
		Comment:         "price changed is " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " for service " + name,
		ReasonForChange: eventData.Comments,
		Currency:        contractEntity.Currency.String(),
	}
	actionQuantityMetadata := SLIActionMetadata{
		UserName:         username,
		ServiceName:      serviceLineItemEntity.Name,
		PreviousQuantity: serviceLineItemEntity.Quantity,
		Quantity:         eventData.Quantity,
		Price:            serviceLineItemEntity.Price,
		BilledType:       serviceLineItemEntity.Billed.String(),
		Comment:          "quantity changed is " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " for service " + name,
		ReasonForChange:  eventData.Comments,
		Currency:         contractEntity.Currency.String(),
	}
	if eventData.StartedAt != nil {
		actionPriceMetadata.StartedAt = eventData.StartedAt
		actionQuantityMetadata.StartedAt = eventData.StartedAt
	}
	metadataPrice, err := utils.ToJson(actionPriceMetadata)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize price metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize price metadata")
	}
	metadataQuantity, err := utils.ToJson(actionQuantityMetadata)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize quantity metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize quantity metadata")
	}
	oldCycle := getBillingCycleNamingConvention(serviceLineItemEntity.Billed.String())
	cycle := getBillingCycleNamingConvention(eventData.Billed)
	extraActionProperties := map[string]interface{}{
		"comments": eventData.Comments,
	}

	if priceChanged && (eventData.Billed == model.AnnuallyBilled.String() || eventData.Billed == model.QuarterlyBilled.String() || eventData.Billed == model.MonthlyBilled.String()) {
		if eventData.Price > serviceLineItemEntity.Price {
			message = username + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = username + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
		}
		_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, commonmodel.CONTRACT, neo4jenum.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers, extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractId, err.Error())
		}
	}

	if priceChanged && eventData.Billed == model.OnceBilled.String() {
		if eventData.Price > serviceLineItemEntity.Price {
			message = username + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.2f", eventData.Price)
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = username + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.2f", eventData.Price)
		}
		_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, commonmodel.CONTRACT, neo4jenum.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers, extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractId, err.Error())
		}
	}
	if priceChanged && eventData.Billed == model.UsageBilled.String() {
		if eventData.Price > serviceLineItemEntity.Price {
			message = username + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.4f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.4f", eventData.Price)
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = username + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.4f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.4f", eventData.Price)
		}
		_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, commonmodel.CONTRACT, neo4jenum.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers, extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractId, err.Error())
		}
	}

	if quantityChanged {
		if eventData.Quantity > serviceLineItemEntity.Quantity {
			message = username + " retroactively increased the quantity of " + name + " from " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " to " + strconv.FormatInt(eventData.Quantity, 10)
		}
		if eventData.Quantity < serviceLineItemEntity.Quantity {
			message = username + " retroactively decreased the quantity of " + name + " from " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " to " + strconv.FormatInt(eventData.Quantity, 10)
		}
		_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, commonmodel.CONTRACT, neo4jenum.ActionServiceLineItemQuantityUpdated, message, metadataQuantity, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers, extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating quantity update action for contract service line item %s: %s", contractId, err.Error())
		}
	}

	if eventData.Source.AppSource != constants.AppSourceCustomerOsApi {
		h.services.CommonServices.RabbitMQService.PublishEventCompleted(ctx, eventData.Tenant, serviceLineItemId, commonmodel.SERVICE_LINE_ITEM, utils.NewEventCompletedDetails().WithUpdate())
	}

	return nil
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
