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
	"strings"
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
	Price         float64 `json:"price"`
	PreviousPrice float64 `json:"previousPrice"`
}
type ActionQuantityMetadata struct {
	UserName         string `json:"user-name"`
	ServiceName      string `json:"service-name"`
	Quantity         int64  `json:"quantity"`
	PreviousQuantity int64  `json:"previousQuantity"`
}
type ActionBilledTypeMetadata struct {
	UserName           string `json:"user-name"`
	ServiceName        string `json:"service-name"`
	BilledType         string `json:"billedType"`
	PreviousBilledType string `json:"previousBilledType"`
}

func (h *ServiceLineItemEventHandler) OnCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnCreate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

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
		sliDbNode, err := h.repositories.ServiceLineItemRepository.GetLatestServiceLineItemByParentId(ctx, eventData.Tenant, eventData.ParentId, eventData.StartedAt)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while getting latest service line item with parent id %s: %s", eventData.ParentId, err.Error())
		}
		if sliDbNode != nil {
			previousServiceLineItem := graph_db.MapDbNodeToServiceLineItemEntity(*sliDbNode)
			previousPrice = previousServiceLineItem.Price
			previousQuantity = previousServiceLineItem.Quantity
			previousBilled = previousServiceLineItem.Billed
		}
	}
	err := h.repositories.ServiceLineItemRepository.CreateForContract(ctx, eventData.Tenant, serviceLineItemId, eventData, isNewVersionForExistingSLI, previousQuantity, previousPrice, previousBilled)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
	err = contractHandler.UpdateRenewalArr(ctx, eventData.Tenant, eventData.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while updating renewal opportunity for contract %s: %s", eventData.ContractId, err.Error())
		return nil
	}

	return nil
}

func (h *ServiceLineItemEventHandler) OnUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnUpdate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)
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
		err = contractHandler.UpdateRenewalArr(ctx, eventData.Tenant, contract.Id)
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
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to serialize billed type metadata: %s", err.Error())
		return errors.Wrap(err, "Failed to serialize billed type metadata")
	}

	if priceChanged {
		if eventData.Price > serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " / " + strings.ToLower(serviceLineItemEntity.Billed) + " to " + fmt.Sprintf("%.2f", eventData.Price) + " / " + strings.ToLower(eventData.Billed)
		}
		if eventData.Price < serviceLineItemEntity.Price {
			message = userEntity.FirstName + " " + userEntity.LastName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " / " + strings.ToLower(serviceLineItemEntity.Billed) + " to " + fmt.Sprintf("%.2f", eventData.Price) + " / " + strings.ToLower(eventData.Billed)
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
		message = userEntity.FirstName + " " + userEntity.LastName + " changed the billing cycle for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " / " + strings.ToLower(serviceLineItemEntity.Billed) + " to " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " / " + strings.ToLower(eventData.Billed)
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
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData event.ServiceLineItemDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	serviceLineItemId := aggregate.GetServiceLineItemObjectID(evt.GetAggregateID(), eventData.Tenant)

	contractDbNode, err := h.repositories.ContractRepository.GetContractByServiceLineItemId(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while getting contract for service line item %s: %s", serviceLineItemId, err.Error())
		return nil
	}

	err = h.repositories.ServiceLineItemRepository.Delete(ctx, eventData.Tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting service line item %s: %s", serviceLineItemId, err.Error())
		return err
	}

	if contractDbNode != nil {
		contract := graph_db.MapDbNodeToContractEntity(contractDbNode)
		contractHandler := contracthandler.NewContractHandler(h.log, h.repositories, h.opportunityCommands)
		err = contractHandler.UpdateRenewalArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contract.Id, err.Error())
			return nil
		}
	}

	return nil
}

func (h *ServiceLineItemEventHandler) OnClose(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemEventHandler.OnClose")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

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
		err = contractHandler.UpdateRenewalArr(ctx, eventData.Tenant, contract.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("error while updating renewal opportunity for contract %s: %s", contract.Id, err.Error())
			return nil
		}
	}

	return nil
}
