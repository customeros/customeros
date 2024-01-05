package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strconv"
)

type ServiceLineItemEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewServiceLineItemEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *ServiceLineItemEventHandler {
	return &ServiceLineItemEventHandler{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

type userMetadata struct {
	UserId string `json:"user-id"`
}

type ActionPriceMetadata struct {
	UserName        string  `json:"user-name"`
	ServiceName     string  `json:"service-name"`
	Quantity        int64   `json:"quantity"`
	Price           float64 `json:"price"`
	PreviousPrice   float64 `json:"previousPrice"`
	BilledType      string  `json:"billedType"`
	Comment         string  `json:"comment"`
	ReasonForChange string  `json:"reasonForChange"`
}
type ActionQuantityMetadata struct {
	UserName         string  `json:"user-name"`
	ServiceName      string  `json:"service-name"`
	Quantity         int64   `json:"quantity"`
	PreviousQuantity int64   `json:"previousQuantity"`
	Price            float64 `json:"price"`
	BilledType       string  `json:"billedType"`
	Comment          string  `json:"comment"`
	ReasonForChange  string  `json:"reasonForChange"`
}
type ActionBilledTypeMetadata struct {
	UserName           string  `json:"user-name"`
	ServiceName        string  `json:"service-name"`
	Price              float64 `json:"price"`
	Quantity           int64   `json:"quantity"`
	BilledType         string  `json:"billedType"`
	PreviousBilledType string  `json:"previousBilledType"`
	Comment            string  `json:"comment"`
	ReasonForChange    string  `json:"reasonForChange"`
}
type ActionServiceLineItemRemovedMetadata struct {
	UserName    string `json:"user-name"`
	ServiceName string `json:"service-name"`
	Comment     string `json:"comment"`
}

func (h *ServiceLineItemEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)
	var user *dbtype.Node
	var userEntity *neo4jentity.UserEntity
	var message string
	var name string
	var priceChanged bool
	var quantityChanged bool
	var billedTypeChanged bool
	var eventData event.ServiceLineItemCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)

	isNewVersionForExistingSLI := serviceLineItemId != eventData.ParentId && eventData.ParentId != ""
	previousPrice := float64(0)
	previousQuantity := int64(0)
	previousBilled := ""
	reasonForChange := eventData.Comments
	if isNewVersionForExistingSLI {
		//get the previous service line item to get the previous price and quantity
		sliDbNode, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, eventData.Tenant, eventData.ParentId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting latest service line item with parent id %s: %s", eventData.ParentId, err.Error())
		}
		if sliDbNode != nil {
			previousServiceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*sliDbNode)
			previousPrice = previousServiceLineItem.Price
			previousQuantity = previousServiceLineItem.Quantity
			previousBilled = previousServiceLineItem.Billed
			//use the booleans below to create the appropriate action message
			priceChanged = previousServiceLineItem.Price != eventData.Price
			quantityChanged = previousServiceLineItem.Quantity != eventData.Quantity
			billedTypeChanged = previousServiceLineItem.Billed != eventData.Billed
		}
	}
	data := neo4jrepository.ServiceLineItemCreateFields{
		IsNewVersionForExistingSLI: isNewVersionForExistingSLI,
		PreviousQuantity:           previousQuantity,
		PreviousPrice:              previousPrice,
		PreviousBilled:             previousBilled,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source.Source),
			SourceOfTruth: helper.GetSourceOfTruth(eventData.Source.SourceOfTruth),
			AppSource:     helper.GetAppSource(eventData.Source.AppSource),
		},
		ContractId: eventData.ContractId,
		ParentId:   eventData.ParentId,
		CreatedAt:  eventData.CreatedAt,
		UpdatedAt:  eventData.UpdatedAt,
		StartedAt:  eventData.StartedAt,
		EndedAt:    eventData.EndedAt,
		Price:      eventData.Price,
		Quantity:   eventData.Quantity,
		Name:       eventData.Name,
		Billed:     eventData.Billed,
		Comments:   eventData.Comments,
	}
	err := h.repositories.Neo4jRepositories.ServiceLineItemWriteRepository.CreateForContract(ctx, eventData.Tenant, serviceLineItemId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}
	serviceLineItemDbNode, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting service line item by id %s: %s", serviceLineItemId, err.Error())
		return err
	}
	serviceLineItemEntity := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)

	contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
	err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating renewal opportunity for contract %s: %s", eventData.ContractId, err.Error())
		return nil
	}
	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := graph_db.MapDbNodeToContractEntity(contractDbNode)

	// get user
	usrMetadata := userMetadata{}
	if err = json.Unmarshal(evt.Metadata, &usrMetadata); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Unmarshal")
	} else {
		if usrMetadata.UserId != "" {
			user, err = h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, eventData.Tenant, usrMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to get user for service line item %s with userid %s", serviceLineItemId, usrMetadata.UserId)
			}
		}
		userEntity = neo4jmapper.MapDbNodeToUserEntity(user)
	}
	if eventData.Name == "" {
		name = serviceLineItemEntity.Name
	}
	if serviceLineItemEntity.Name != "" {
		name = serviceLineItemEntity.Name
	}
	if name == "" {
		name = "Unnamed service"
	}

	metadataPrice, err := utils.ToJson(ActionPriceMetadata{
		UserName:        userEntity.FirstName + " " + userEntity.LastName,
		ServiceName:     name,
		Quantity:        eventData.Quantity,
		BilledType:      eventData.Billed,
		PreviousPrice:   previousPrice,
		Price:           eventData.Price,
		Comment:         "price is " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " for service " + name,
		ReasonForChange: reasonForChange,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize price metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize price metadata")
	}
	metadataQuantity, err := utils.ToJson(ActionQuantityMetadata{
		UserName:         userEntity.FirstName + " " + userEntity.LastName,
		ServiceName:      name,
		Price:            eventData.Price,
		PreviousQuantity: previousQuantity,
		Quantity:         eventData.Quantity,
		BilledType:       eventData.Billed,
		Comment:          "quantity is " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " for service " + name,
		ReasonForChange:  reasonForChange,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize quantity metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize quantity metadata")
	}
	metadataBilledType, err := utils.ToJson(ActionBilledTypeMetadata{
		UserName:           userEntity.FirstName + " " + userEntity.LastName,
		ServiceName:        name,
		BilledType:         eventData.Billed,
		PreviousBilledType: previousBilled,
		Quantity:           eventData.Quantity,
		Price:              eventData.Price,
		Comment:            "billed type is " + serviceLineItemEntity.Billed + " for service " + name,
		ReasonForChange:    reasonForChange,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize billed type metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize billed type metadata")
	}
	extraActionProperties := map[string]interface{}{
		"comments": reasonForChange,
	}
	cycle := getBillingCycleNamingConvention(eventData.Billed)
	previousCycle := getBillingCycleNamingConvention(previousBilled)
	if previousCycle == "" {
		previousCycle = getBillingCycleNamingConvention(serviceLineItemEntity.Billed)
	}

	if !isNewVersionForExistingSLI {
		if serviceLineItemEntity.Billed == model.AnnuallyBilled.String() || serviceLineItemEntity.Billed == model.QuarterlyBilled.String() || serviceLineItemEntity.Billed == model.MonthlyBilled.String() {
			message = userEntity.FirstName + " " + userEntity.LastName + " added a recurring service to " + contractEntity.Name + ": " + name + " at " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " x " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + cycle
			_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, eventData.ContractId, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemBilledTypeRecurringCreated, message, metadataBilledType, utils.Now(), extraActionProperties)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating recurring billed type service line item created action for contract %s: %s", eventData.ContractId, err.Error())
			}
		}
		if serviceLineItemEntity.Billed == model.OnceBilled.String() {
			message = userEntity.FirstName + " " + userEntity.LastName + " added an one time service to " + contractEntity.Name + ": " + name + " at " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price)
			_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, eventData.ContractId, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemBilledTypeOnceCreated, message, metadataBilledType, utils.Now(), extraActionProperties)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating once billed type service line item created action for contract %s: %s", eventData.ContractId, err.Error())
			}
		}
		if serviceLineItemEntity.Billed == model.UsageBilled.String() {
			message = userEntity.FirstName + " " + userEntity.LastName + " added a per use service to " + contractEntity.Name + ": " + name + " at " + fmt.Sprintf("%.4f", serviceLineItemEntity.Price)
			_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, eventData.ContractId, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemBilledTypeUsageCreated, message, metadataBilledType, utils.Now(), extraActionProperties)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating per use billed type service line item created action for contract %s: %s", eventData.ContractId, err.Error())
			}
		}
	}
	if isNewVersionForExistingSLI {
		if priceChanged && (eventData.Billed == model.AnnuallyBilled.String() || eventData.Billed == model.QuarterlyBilled.String() || eventData.Billed == model.MonthlyBilled.String()) {
			if eventData.Price > previousPrice {
				message = userEntity.FirstName + " " + userEntity.LastName + " increased the price for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + "/" + previousCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
			}
			if eventData.Price < previousPrice {
				message = userEntity.FirstName + " " + userEntity.LastName + " decreased the price for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + "/" + previousCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
			}
			_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractEntity.Id, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), extraActionProperties)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractEntity.Id, err.Error())
			}
		}

		if priceChanged && eventData.Billed == model.OnceBilled.String() {
			if eventData.Price > previousPrice {
				message = userEntity.FirstName + " " + userEntity.LastName + " increased the price for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + " to " + fmt.Sprintf("%.2f", eventData.Price)
			}
			if eventData.Price < serviceLineItemEntity.Price {
				message = userEntity.FirstName + " " + userEntity.LastName + " decreased the price for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + " to " + fmt.Sprintf("%.2f", eventData.Price)
			}
			_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractEntity.Id, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), extraActionProperties)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractEntity.Id, err.Error())
			}
		}
		if priceChanged && eventData.Billed == model.UsageBilled.String() {
			if eventData.Price > previousPrice {
				message = userEntity.FirstName + " " + userEntity.LastName + " increased the price for " + name + " from " + fmt.Sprintf("%.4f", previousPrice) + " to " + fmt.Sprintf("%.4f", eventData.Price)
			}
			if eventData.Price < serviceLineItemEntity.Price {
				message = userEntity.FirstName + " " + userEntity.LastName + " decreased the price for " + name + " from " + fmt.Sprintf("%.4f", previousPrice) + " to " + fmt.Sprintf("%.4f", eventData.Price)
			}
			_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractEntity.Id, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), extraActionProperties)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractEntity.Id, err.Error())
			}
		}
		if quantityChanged {
			if eventData.Quantity > previousQuantity {
				message = userEntity.FirstName + " " + userEntity.LastName + " increased the quantity of " + name + " from " + strconv.FormatInt(previousQuantity, 10) + " to " + strconv.FormatInt(eventData.Quantity, 10)
			}
			if eventData.Quantity < previousQuantity {
				message = userEntity.FirstName + " " + userEntity.LastName + " decreased the quantity of " + name + " from " + strconv.FormatInt(previousQuantity, 10) + " to " + strconv.FormatInt(eventData.Quantity, 10)
			}
			_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractEntity.Id, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemQuantityUpdated, message, metadataQuantity, utils.Now(), extraActionProperties)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating quantity update action for contract service line item %s: %s", contractEntity.Id, err.Error())
			}
		}
		if billedTypeChanged && previousBilled != "" {
			message = userEntity.FirstName + " " + userEntity.LastName + " changed the billing cycle for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + "/" + previousCycle + " to " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + cycle
			_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractEntity.Id, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemBilledTypeUpdated, message, metadataBilledType, utils.Now(), extraActionProperties)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating billed type update action for contract service line item %s: %s", contractEntity.Id, err.Error())
			}
		}
	}

	return nil
}

func (h *ServiceLineItemEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)
	var contractId string
	var user *dbtype.Node
	var userEntity *neo4jentity.UserEntity
	var name string
	var message string
	var eventData event.ServiceLineItemUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)
	serviceLineItemDbNode, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	serviceLineItemEntity := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	//we will use the following booleans below to check if the price, quantity, billed type has changed
	priceChanged := serviceLineItemEntity.Price != eventData.Price
	quantityChanged := serviceLineItemEntity.Quantity != eventData.Quantity
	billedTypeChanged := serviceLineItemEntity.Billed != eventData.Billed

	data := neo4jrepository.ServiceLineItemUpdateFields{
		Price:     eventData.Price,
		Quantity:  eventData.Quantity,
		Billed:    eventData.Billed,
		Comments:  eventData.Comments,
		Name:      eventData.Name,
		Source:    helper.GetSource(eventData.Source.Source),
		UpdatedAt: eventData.UpdatedAt,
	}
	err = h.repositories.Neo4jRepositories.ServiceLineItemWriteRepository.Update(ctx, eventData.Tenant, serviceLineItemId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	if contractDbNode != nil {
		contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
		contractId = contract.Id
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contract.Id, err.Error())
			return nil
		}
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
	usrMetadata := userMetadata{}
	if err = json.Unmarshal(evt.Metadata, &usrMetadata); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Unmarshal")
	} else {
		if usrMetadata.UserId != "" {
			user, err = h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, eventData.Tenant, usrMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to get user for service line item %s with userid %s", serviceLineItemId, usrMetadata.UserId)
			}
		}
		userEntity = neo4jmapper.MapDbNodeToUserEntity(user)
	}

	metadataPrice, err := utils.ToJson(ActionPriceMetadata{
		UserName:        userEntity.FirstName + " " + userEntity.LastName,
		ServiceName:     serviceLineItemEntity.Name,
		Price:           eventData.Price,
		PreviousPrice:   serviceLineItemEntity.Price,
		BilledType:      serviceLineItemEntity.Billed,
		Quantity:        serviceLineItemEntity.Quantity,
		Comment:         "price changed is " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " for service " + name,
		ReasonForChange: eventData.Comments,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize price metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize price metadata")
	}
	metadataQuantity, err := utils.ToJson(ActionQuantityMetadata{
		UserName:         userEntity.FirstName + " " + userEntity.LastName,
		ServiceName:      serviceLineItemEntity.Name,
		PreviousQuantity: serviceLineItemEntity.Quantity,
		Quantity:         eventData.Quantity,
		Price:            serviceLineItemEntity.Price,
		BilledType:       serviceLineItemEntity.Billed,
		Comment:          "quantity changed is " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " for service " + name,
		ReasonForChange:  eventData.Comments,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize quantity metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize quantity metadata")
	}
	metadataBilledType, err := utils.ToJson(ActionBilledTypeMetadata{
		UserName:           userEntity.FirstName + " " + userEntity.LastName,
		ServiceName:        serviceLineItemEntity.Name,
		BilledType:         eventData.Billed,
		PreviousBilledType: serviceLineItemEntity.Billed,
		Quantity:           serviceLineItemEntity.Quantity,
		Price:              serviceLineItemEntity.Price,
		Comment:            "billed type changed is " + serviceLineItemEntity.Billed + " for service " + name,
		ReasonForChange:    eventData.Comments,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize billed type metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize billed type metadata")
	}
	oldCycle := getBillingCycleNamingConvention(serviceLineItemEntity.Billed)
	cycle := getBillingCycleNamingConvention(eventData.Billed)
	extraActionProperties := map[string]interface{}{
		"comments": eventData.Comments,
	}

	if priceChanged && (eventData.Billed == model.AnnuallyBilled.String() || eventData.Billed == model.QuarterlyBilled.String() || eventData.Billed == model.MonthlyBilled.String()) {
		if eventData.Price > serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
		}
		_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractId, err.Error())
		}
	}

	if priceChanged && eventData.Billed == model.OnceBilled.String() {
		if eventData.Price > serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.2f", eventData.Price)
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.2f", eventData.Price)
		}
		_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractId, err.Error())
		}
	}
	if priceChanged && eventData.Billed == model.UsageBilled.String() {
		if eventData.Price > serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.4f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.4f", eventData.Price)
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.4f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.4f", eventData.Price)
		}
		_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractId, err.Error())
		}
	}

	if quantityChanged {
		if eventData.Quantity > serviceLineItemEntity.Quantity {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively increased the quantity of " + name + " from " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " to " + strconv.FormatInt(eventData.Quantity, 10)
		}
		if eventData.Quantity < serviceLineItemEntity.Quantity {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively decreased the quantity of " + name + " from " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " to " + strconv.FormatInt(eventData.Quantity, 10)
		}
		_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemQuantityUpdated, message, metadataQuantity, utils.Now(), extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating quantity update action for contract service line item %s: %s", contractId, err.Error())
		}
	}
	if billedTypeChanged && serviceLineItemEntity.Billed != "" {
		message = userEntity.FirstName + " " + userEntity.LastName + " changed the billing cycle for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + cycle
		_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, contractId, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemBilledTypeUpdated, message, metadataBilledType, utils.Now(), extraActionProperties)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating billed type update action for contract service line item %s: %s", contractId, err.Error())
		}
	}
	return nil
}

func (h *ServiceLineItemEventHandler) OnDelete(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnDelete")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)
	var user *dbtype.Node
	var userEntity *neo4jentity.UserEntity
	var serviceLineItemName string
	var contractName string

	var eventData event.ServiceLineItemDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)
	serviceLineItemDbNode, err := h.repositories.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	serviceLineItemEntity := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	if serviceLineItemEntity.Name != "" {
		serviceLineItemName = serviceLineItemEntity.Name
	} else {
		serviceLineItemName = "Unnamed service"
	}

	// get user
	usrMetadata := userMetadata{}
	if err := json.Unmarshal(evt.Metadata, &usrMetadata); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Unmarshal")
	} else {
		if usrMetadata.UserId != "" {
			user, err = h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, eventData.Tenant, usrMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to get user for service line item %s with userid %s", serviceLineItemId, usrMetadata.UserId)
			}
		}
		userEntity = neo4jmapper.MapDbNodeToUserEntity(user)
	}

	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	if contract.Name != "" {
		contractName = contract.Name
	} else {
		contractName = "Unnamed contract"
	}

	err = h.repositories.Neo4jRepositories.ServiceLineItemWriteRepository.Delete(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	if contractDbNode != nil {
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contract.Id, err.Error())
			return nil
		}
	}
	metadata, err := utils.ToJson(ActionServiceLineItemRemovedMetadata{
		UserName:    userEntity.FirstName + " " + userEntity.LastName,
		ServiceName: serviceLineItemName,
		Comment:     "service line item removed is " + serviceLineItemName + " from " + contractName + " by " + userEntity.FirstName + " " + userEntity.LastName + "",
	})
	message := userEntity.FirstName + " " + userEntity.LastName + " removed " + serviceLineItemName + " from " + contractName

	_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.Create(ctx, eventData.Tenant, contract.Id, neo4jentity.CONTRACT, neo4jentity.ActionServiceLineItemRemoved, message, metadata, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed remove service line item action for contract %s: %s", contract.Id, err.Error())
	}

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
	err := h.repositories.Neo4jRepositories.ServiceLineItemWriteRepository.Close(ctx, eventData.Tenant, serviceLineItemId, eventData.UpdatedAt, eventData.EndedAt, eventData.IsCanceled)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while closing service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	contractDbNode, err := h.repositories.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	if contractDbNode != nil {
		contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.grpcClients)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contract.Id, err.Error())
			return nil
		}
	}

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
