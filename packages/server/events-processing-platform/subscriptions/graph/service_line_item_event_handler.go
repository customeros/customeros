package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	contracthandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions/contract"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strconv"
)

type ServiceLineItemEventHandler struct {
	log                 logger.Logger
	repositories        *repository.Repositories
	opportunityCommands *opportunitycmdhandler.CommandHandlers
}
type userMetadata struct {
	UserId string `json:"user-id"`
}

type ActionPriceMetadata struct {
	UserName      string  `json:"user-name"`
	ServiceName   string  `json:"service-name"`
	Quantity      int64   `json:"quantity"`
	Price         float64 `json:"price"`
	PreviousPrice float64 `json:"previousPrice"`
	BilledType    string  `json:"billedType"`
	Comment       string  `json:"comment"`
}
type ActionQuantityMetadata struct {
	UserName         string  `json:"user-name"`
	ServiceName      string  `json:"service-name"`
	Quantity         int64   `json:"quantity"`
	PreviousQuantity int64   `json:"previousQuantity"`
	Price            float64 `json:"price"`
	BilledType       string  `json:"billedType"`
	Comment          string  `json:"comment"`
}
type ActionBilledTypeMetadata struct {
	UserName           string  `json:"user-name"`
	ServiceName        string  `json:"service-name"`
	Price              float64 `json:"price"`
	Quantity           int64   `json:"quantity"`
	BilledType         string  `json:"billedType"`
	PreviousBilledType string  `json:"previousBilledType"`
	Comment            string  `json:"comment"`
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
	var userEntity *entity.UserEntity
	var message string
	var name string
	var priceChanged bool
	var quantityChanged bool
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
	if isNewVersionForExistingSLI {
		//get the previous service line item to get the previous price and quantity
		sliDbNode, err := h.repositories.ServiceLineItemRepository.GetServiceLineItemById(ctx, eventData.Tenant, eventData.ParentId)
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
		}
	}
	err := h.repositories.ServiceLineItemRepository.CreateForContract(ctx, eventData.Tenant, serviceLineItemId, eventData, isNewVersionForExistingSLI, previousQuantity, previousPrice, previousBilled)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}
	serviceLineItemDbNode, err := h.repositories.ServiceLineItemRepository.GetServiceLineItemById(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting service line item by id %s: %s", serviceLineItemId, err.Error())
		return err
	}
	serviceLineItemEntity := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)

	contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
	err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating renewal opportunity for contract %s: %s", eventData.ContractId, err.Error())
		return nil
	}
	contractDbNode, err := h.repositories.ContractRepository.GetContractById(ctx, eventData.Tenant, eventData.ContractId)
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
			user, err = h.repositories.UserRepository.GetUser(ctx, eventData.Tenant, usrMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to get user for service line item %s with userid %s", serviceLineItemId, usrMetadata.UserId)
			}
		}
		userEntity = graph_db.MapDbNodeToUserEntity(*user)
	}
	if eventData.Name == "" {
		name = serviceLineItemEntity.Name
	}
	if serviceLineItemEntity.Name != "" {
		name = serviceLineItemEntity.Name
	}
	if name == "" {
		name = "unnamed service"
	}

	metadataPrice, err := utils.ToJson(ActionPriceMetadata{
		UserName:      userEntity.FirstName + " " + userEntity.LastName,
		ServiceName:   name,
		Quantity:      eventData.Quantity,
		BilledType:    eventData.Billed,
		PreviousPrice: previousPrice,
		Price:         eventData.Price,
		Comment:       "price is " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " for service " + name,
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
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize quantity metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize quantity metadata")
	}
	metadataBilledType, err := utils.ToJson(ActionBilledTypeMetadata{
		UserName:    userEntity.FirstName + " " + userEntity.LastName,
		ServiceName: name,
		BilledType:  eventData.Billed,
		Quantity:    eventData.Quantity,
		Price:       eventData.Price,
		Comment:     "billed type is " + serviceLineItemEntity.Billed + " for service " + name,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize billed type metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize billed type metadata")
	}

	var oldCycle string
	switch serviceLineItemEntity.Billed {
	case model.AnnuallyBilled.String():
		oldCycle = "year"
	case model.QuarterlyBilled.String():
		oldCycle = "quarter"
	case model.MonthlyBilled.String():
		oldCycle = "month"
	}

	var cycle string
	switch eventData.Billed {
	case model.AnnuallyBilled.String():
		cycle = "year"
	case model.QuarterlyBilled.String():
		cycle = "quarter"
	case model.MonthlyBilled.String():
		cycle = "month"
	}
	if !isNewVersionForExistingSLI {
		if serviceLineItemEntity.Billed == model.AnnuallyBilled.String() || serviceLineItemEntity.Billed == model.QuarterlyBilled.String() || serviceLineItemEntity.Billed == model.MonthlyBilled.String() {
			message = userEntity.FirstName + " " + userEntity.LastName + " added a recurring service to " + contractEntity.Name + ": " + name + " at " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " x " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle
			_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, eventData.ContractId, entity.CONTRACT, entity.ActionServiceLineItemBilledTypeRecurringCreated, message, metadataBilledType, utils.Now())
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating recurring billed type service line item created action for contract %s: %s", eventData.ContractId, err.Error())
			}
		}
		if serviceLineItemEntity.Billed == model.OnceBilled.String() {
			message = userEntity.FirstName + " " + userEntity.LastName + " added an one time service to " + contractEntity.Name + ": " + name + " at " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price)
			_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, eventData.ContractId, entity.CONTRACT, entity.ActionServiceLineItemBilledTypeOnceCreated, message, metadataBilledType, utils.Now())
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating once billed type service line item created action for contract %s: %s", eventData.ContractId, err.Error())
			}
		}
		if serviceLineItemEntity.Billed == model.UsageBilled.String() {
			message = userEntity.FirstName + " " + userEntity.LastName + " added a per use service to " + contractEntity.Name + ": " + name + " at " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price)
			_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, eventData.ContractId, entity.CONTRACT, entity.ActionServiceLineItemBilledTypeUsageCreated, message, metadataBilledType, utils.Now())
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating per use billed type service line item created action for contract %s: %s", eventData.ContractId, err.Error())
			}
		}
	}
	if isNewVersionForExistingSLI {
		if priceChanged && (eventData.Billed == model.AnnuallyBilled.String() || eventData.Billed == model.QuarterlyBilled.String() || eventData.Billed == model.MonthlyBilled.String()) {
			if eventData.Price > previousPrice {
				message = userEntity.FirstName + " " + userEntity.LastName + " increased the price for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
			}
			if eventData.Price < previousPrice {
				message = userEntity.FirstName + " " + userEntity.LastName + " decreased the price for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
			}
			_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractEntity.Id, entity.CONTRACT, entity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now())
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractEntity.Id, err.Error())
			}
		}

		if priceChanged && (eventData.Billed == model.OnceBilled.String() || eventData.Billed == model.UsageBilled.String()) {
			if eventData.Price > previousPrice {
				message = userEntity.FirstName + " " + userEntity.LastName + " increased the price for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + " to " + fmt.Sprintf("%.2f", eventData.Price)
			}
			if eventData.Price < serviceLineItemEntity.Price {
				message = userEntity.FirstName + " " + userEntity.LastName + " decreased the price for " + name + " from " + fmt.Sprintf("%.2f", previousPrice) + " to " + fmt.Sprintf("%.2f", eventData.Price)
			}
			_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractEntity.Id, entity.CONTRACT, entity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now())
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
			_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractEntity.Id, entity.CONTRACT, entity.ActionServiceLineItemQuantityUpdated, message, metadataQuantity, utils.Now())
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating quantity update action for contract service line item %s: %s", contractEntity.Id, err.Error())
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
	var userEntity *entity.UserEntity
	var name string
	var message string
	var eventData event.ServiceLineItemUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)
	serviceLineItemDbNode, err := h.repositories.ServiceLineItemRepository.GetServiceLineItemById(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	serviceLineItemEntity := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	//we will use the following booleans below to check if the price, quantity, billed type has changed
	priceChanged := serviceLineItemEntity.Price != eventData.Price
	quantityChanged := serviceLineItemEntity.Quantity != eventData.Quantity
	billedTypeChanged := serviceLineItemEntity.Billed != eventData.Billed

	err = h.repositories.ServiceLineItemRepository.Update(ctx, eventData.Tenant, serviceLineItemId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while updating service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	contractDbNode, err := h.repositories.ContractRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	if contractDbNode != nil {
		contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
		contractId = contract.Id
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
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
		name = "unnamed service"
	}
	// get user
	usrMetadata := userMetadata{}
	if err = json.Unmarshal(evt.Metadata, &usrMetadata); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Unmarshal")
	} else {
		if usrMetadata.UserId != "" {
			user, err = h.repositories.UserRepository.GetUser(ctx, eventData.Tenant, usrMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to get user for service line item %s with userid %s", serviceLineItemId, usrMetadata.UserId)
			}
		}
		userEntity = graph_db.MapDbNodeToUserEntity(*user)
	}

	metadataPrice, err := utils.ToJson(ActionPriceMetadata{
		UserName:      userEntity.FirstName + " " + userEntity.LastName,
		ServiceName:   serviceLineItemEntity.Name,
		Price:         eventData.Price,
		PreviousPrice: serviceLineItemEntity.Price,
		BilledType:    serviceLineItemEntity.Billed,
		Quantity:      serviceLineItemEntity.Quantity,
		Comment:       "price changed is " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " for service " + name,
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
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize billed type metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize billed type metadata")
	}
	var oldCycle string
	switch serviceLineItemEntity.Billed {
	case model.AnnuallyBilled.String():
		oldCycle = "year"
	case model.QuarterlyBilled.String():
		oldCycle = "quarter"
	case model.MonthlyBilled.String():
		oldCycle = "month"
	}

	var cycle string
	switch eventData.Billed {
	case model.AnnuallyBilled.String():
		cycle = "year"
	case model.QuarterlyBilled.String():
		cycle = "quarter"
	case model.MonthlyBilled.String():
		cycle = "month"
	}

	if priceChanged && (eventData.Billed == model.AnnuallyBilled.String() || eventData.Billed == model.QuarterlyBilled.String() || eventData.Billed == model.MonthlyBilled.String()) {
		if eventData.Price > serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", eventData.Price) + "/" + cycle
		}
		_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractId, entity.CONTRACT, entity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now())
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractId, err.Error())
		}
	}

	if priceChanged && (eventData.Billed == model.OnceBilled.String() || eventData.Billed == model.UsageBilled.String()) {
		if eventData.Price > serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.2f", eventData.Price)
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.2f", eventData.Price)
		}
		_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractId, entity.CONTRACT, entity.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now())
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
		_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractId, entity.CONTRACT, entity.ActionServiceLineItemQuantityUpdated, message, metadataQuantity, utils.Now())
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed creating quantity update action for contract service line item %s: %s", contractId, err.Error())
		}
	}
	if billedTypeChanged && serviceLineItemEntity.Billed != "" {
		message = userEntity.FirstName + " " + userEntity.LastName + " changed the billing cycle for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + cycle
		_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contractId, entity.CONTRACT, entity.ActionServiceLineItemBilledTypeUpdated, message, metadataBilledType, utils.Now())
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
	var userEntity *entity.UserEntity
	var serviceLineItemName string
	var contractName string

	var eventData event.ServiceLineItemDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)
	serviceLineItemDbNode, err := h.repositories.ServiceLineItemRepository.GetServiceLineItemById(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	serviceLineItemEntity := graph_db.MapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode)
	if serviceLineItemEntity.Name != "" {
		serviceLineItemName = serviceLineItemEntity.Name
	} else {
		serviceLineItemName = "unnamed service"
	}

	// get user
	usrMetadata := userMetadata{}
	if err := json.Unmarshal(evt.Metadata, &usrMetadata); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Unmarshal")
	} else {
		if usrMetadata.UserId != "" {
			user, err = h.repositories.UserRepository.GetUser(ctx, eventData.Tenant, usrMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to get user for service line item %s with userid %s", serviceLineItemId, usrMetadata.UserId)
			}
		}
		userEntity = graph_db.MapDbNodeToUserEntity(*user)
	}

	contractDbNode, err := h.repositories.ContractRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
	if contract.Name != "" {
		contractName = contract.Name
	} else {
		contractName = "unnamed contract"
	}

	err = h.repositories.ServiceLineItemRepository.Delete(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	if contractDbNode != nil {
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
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

	_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, contract.Id, entity.CONTRACT, entity.ActionServiceLineItemRemoved, message, metadata, utils.Now())
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
	err := h.repositories.ServiceLineItemRepository.Close(ctx, eventData.Tenant, serviceLineItemId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while closing service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	contractDbNode, err := h.repositories.ContractRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}
	if contractDbNode != nil {
		contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
		err = contractHandler.UpdateActiveRenewalOpportunityArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contract.Id, err.Error())
			return nil
		}
	}

	return nil
}
