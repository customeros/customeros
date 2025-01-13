package sli

import (
	"context"
	"fmt"
	"strconv"
	"time"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

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

type serviceLineItemService struct {
	log      logger.Logger
	events   *events.EventsService
	neo4j    *neoRepo.Repositories
	contract interfaces.ContractService
}

func NewServiceLineItemService(log logger.Logger, events *events.EventsService, neo4j *neoRepo.Repositories, contract interfaces.ContractService) interfaces.ServiceLineItemService {
	return &serviceLineItemService{
		log:      log,
		events:   events,
		neo4j:    neo4j,
		contract: contract,
	}
}

func (s *serviceLineItemService) SetContractService(contract interfaces.ContractService) {
	s.contract = contract
}

func (s *serviceLineItemService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *serviceLineItemService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, dataFields data_fields.SLIFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	sliId := ""
	priceChanged := false
	quantityChanged := false

	if utils.IfNotNilString(id) == "" {
		createFlow = true
		span.LogKV("flow", "create")

		// prepare missing fields
		if dataFields.CreatedAt == nil || dataFields.CreatedAt.IsZero() {
			dataFields.CreatedAt = utils.NowPtr()
		}
		if dataFields.StartedAt == nil || dataFields.StartedAt.IsZero() {
			dataFields.StartedAt = utils.NowPtr()
		}
		if utils.IfNotNilString(dataFields.Source) == "" {
			dataFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
		if utils.IfNotNilString(dataFields.AppSource) == "" {
			dataFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		if utils.IfNotNilFloat64(dataFields.TaxRate) < 0 {
			dataFields.TaxRate = utils.Float64Ptr(0)
		}
		dataFields.TaxRate = utils.Float64Ptr(utils.TruncateFloat64(utils.IfNotNilFloat64(dataFields.TaxRate), 2))

		// validate data
		if utils.IfNotNilInt64(dataFields.Quantity) < 0 {
			err = errors.New("quantity cannot be negative")
			tracing.TraceErr(span, err)
			return "", err
		}

		if dataFields.EndedAt != nil && dataFields.EndedAt.Before(*dataFields.StartedAt) {
			err = errors.New("endedAt cannot be before startedAt")
			tracing.TraceErr(span, err)
			return "", err
		}

		if utils.IfNotNilString(dataFields.ContractId) == "" {
			err = errors.New("contractId is required")
			tracing.TraceErr(span, err)
			return "", err
		}

		if utils.IfNotNilString(dataFields.ContractId) != "" {
			exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, *dataFields.ContractId, model.NodeLabelContract)
			if err != nil || !exists {
				err = errors.New("contract not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		// generate id
		sliId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelServiceLineItem)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if utils.IfNotNilString(dataFields.ParentId) == "" {
			dataFields.ParentId = utils.StringPtr(sliId)
		}
	} else {
		span.LogKV("flow", "update")
		sliId = *id

		sliEntity, err := s.GetById(ctx, sliId)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}

		// reset non-updatable fields
		if sliEntity.Billed != neo4jenum.BilledTypeNone {
			dataFields.BilledType = utils.ToPtr(sliEntity.Billed)
		}
		if utils.IfNotNilFloat64(dataFields.TaxRate) < 0 {
			dataFields.TaxRate = utils.Float64Ptr(0)
		}
		dataFields.TaxRate = utils.Float64Ptr(utils.TruncateFloat64(utils.IfNotNilFloat64(dataFields.TaxRate), 2))

		// validate data
		if sliEntity.Canceled {
			err = errors.New("service line item is canceled")
			tracing.TraceErr(span, err)
			return "", err
		}

		// validate data
		if utils.IfNotNilInt64(dataFields.Quantity) < 0 {
			err = errors.New("quantity cannot be negative")
			tracing.TraceErr(span, err)
			return "", err
		}

		priceChanged = dataFields.Price != nil && sliEntity.Price != *dataFields.Price
		quantityChanged = dataFields.Quantity != nil && sliEntity.Quantity != *dataFields.Quantity
	}
	tracing.TagEntity(span, sliId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if createFlow {
			err := s.neo4j.ServiceLineItemWriteRepository.CreateForContract(ctx, txWithPostCommit.Tx, tenant, sliId, dataFields)
			if err != nil {
				s.log.Errorf("error creating service line item %s: %s", sliId, err.Error())
				return nil, err
			}
		} else {
			err := s.neo4j.ServiceLineItemWriteRepository.Update(ctx, txWithPostCommit.Tx, tenant, sliId, dataFields)
			if err != nil {
				s.log.Errorf("error updating service line item %s: %s", sliId, err.Error())
				return nil, err
			}
		}

		err := s.neo4j.ServiceLineItemWriteRepository.AdjustEndDates(ctx, txWithPostCommit.Tx, tenant, utils.IfNotNilString(dataFields.ParentId))
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while adjusting end dates for service line item %s: %s", sliId, err.Error())
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			contractDbNode, err := s.neo4j.ContractReadRepository.GetContractByServiceLineItemId(ctx, tenant, sliId)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

			err = s.contract.UpdateActiveRenewalOpportunityArr(ctx, *dataFields.ContractId)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			err = s.contract.RecalculateContractLtv(ctx, *dataFields.ContractId)
			if err != nil {
				tracing.TraceErr(span, err)
			}

			if createFlow {
				if dataFields.BilledType != nil && utils.IfNotNilString(dataFields.BilledType.String()) != "" {
					name := "Unnamed service"
					if utils.IfNotNilString(dataFields.Name) != "" {
						name = *dataFields.Name
					}

					userName := ""
					userDbNode, err := s.neo4j.UserReadRepository.GetUserById(ctx, tenant, common.GetUserIdFromContext(ctx))
					if err != nil {
						tracing.TraceErr(span, err)
					}
					if userDbNode != nil {
						userEntity := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
						userName = userEntity.GetFullName()
					}
					extraActionProperties := map[string]interface{}{
						"comments": utils.IfNotNilString(dataFields.Comments),
					}
					cycle := getBillingCycleNamingConvention(dataFields.BilledType.String())

					metadataBilledType, err := utils.ToJson(SLIActionMetadata{
						UserName:        userName,
						ServiceName:     name,
						BilledType:      dataFields.BilledType.String(),
						Quantity:        utils.IfNotNilInt64(dataFields.Quantity),
						Price:           utils.IfNotNilFloat64(dataFields.Price),
						Comment:         "billed type is " + dataFields.BilledType.String() + " for service " + name,
						ReasonForChange: utils.IfNotNilString(dataFields.Comments),
						StartedAt:       dataFields.StartedAt,
						Currency:        contractEntity.Currency.String(),
					})
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("Failed to serialize billed type metadata: %s", err.Error())
						return err
					}
					if dataFields.BilledType.IsRecurrent() {
						message := userName + " added a recurring service to " + contractEntity.Name + ": " + name + " at " + strconv.FormatInt(utils.IfNotNilInt64(dataFields.Quantity), 10) + " x " + fmt.Sprintf("%.2f", utils.IfNotNilFloat64(dataFields.Price)) + "/" + cycle + " starting with " + dataFields.StartedAt.Format("2006-01-02")
						_, err = s.neo4j.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, enum.ActionServiceLineItemBilledTypeRecurringCreated, message, metadataBilledType, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
						if err != nil {
							tracing.TraceErr(span, err)
							s.log.Errorf("Failed creating recurring billed type service line item created action for contract %s: %s", contractEntity.Id, err.Error())
						}
					}
					if *dataFields.BilledType == neo4jenum.BilledTypeOnce {
						message := userName + " added a one time service to " + contractEntity.Name + ": " + name + " at " + fmt.Sprintf("%.2f", utils.IfNotNilFloat64(dataFields.Price)) + " starting with " + dataFields.StartedAt.Format("2006-01-02")
						_, err = s.neo4j.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, enum.ActionServiceLineItemBilledTypeOnceCreated, message, metadataBilledType, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
						if err != nil {
							tracing.TraceErr(span, err)
							s.log.Errorf("Failed creating once billed type service line item created action for contract %s: %s", contractEntity.Id, err.Error())
						}
					}
					if *dataFields.BilledType == neo4jenum.BilledTypeUsage {
						message := userName + " added a per use service to " + contractEntity.Name + ": " + name + " at " + fmt.Sprintf("%.4f", utils.IfNotNilFloat64(dataFields.Price)) + " starting with " + dataFields.StartedAt.Format("2006-01-02")
						_, err = s.neo4j.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, enum.ActionServiceLineItemBilledTypeUsageCreated, message, metadataBilledType, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
						if err != nil {
							tracing.TraceErr(span, err)
							s.log.Errorf("Failed creating per use billed type service line item created action for contract %s: %s", contractEntity.Id, err.Error())
						}
					}
				}

				err = s.events.Publisher.PublishEvent(ctx, sliId, model.SERVICE_LINE_ITEM, dto.CreateServiceLineItem{dataFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateServiceLineItem for SLI"))
				}
				err = s.events.Publisher.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.CreateServiceLineItem{dataFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateServiceLineItem for Contract"))
				}

				s.events.Publisher.PublishEventCompleted(ctx, tenant, sliId, model.SERVICE_LINE_ITEM, utils.NewEventCompletedDetails().WithCreate())
				s.events.Publisher.PublishEventCompleted(ctx, tenant, contractEntity.Id, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())
			} else {
				name := "Unnamed service"
				if utils.IfNotNilString(dataFields.Name) != "" {
					name = *dataFields.Name
				}

				userName := ""
				userDbNode, err := s.neo4j.UserReadRepository.GetUserById(ctx, tenant, common.GetUserIdFromContext(ctx))
				if err != nil {
					tracing.TraceErr(span, err)
				}
				if userDbNode != nil {
					userEntity := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
					userName = userEntity.GetFullName()
				}
				extraActionProperties := map[string]interface{}{
					"comments": utils.IfNotNilString(dataFields.Comments),
				}
				cycle := getBillingCycleNamingConvention(dataFields.BilledType.String())

				serviceLineItemEntity, err := s.GetById(ctx, sliId)
				if err != nil {
					tracing.TraceErr(span, err)
					return err
				}

				actionPriceMetadata := SLIActionMetadata{
					UserName:        userName,
					ServiceName:     serviceLineItemEntity.Name,
					Price:           utils.IfNotNilFloat64(dataFields.Price),
					PreviousPrice:   serviceLineItemEntity.Price,
					BilledType:      serviceLineItemEntity.Billed.String(),
					Quantity:        serviceLineItemEntity.Quantity,
					Comment:         "price changed is " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " for service " + name,
					ReasonForChange: utils.IfNotNilString(dataFields.Comments),
					Currency:        contractEntity.Currency.String(),
				}
				actionQuantityMetadata := SLIActionMetadata{
					UserName:         userName,
					ServiceName:      serviceLineItemEntity.Name,
					PreviousQuantity: serviceLineItemEntity.Quantity,
					Quantity:         utils.IfNotNilInt64(dataFields.Quantity),
					Price:            serviceLineItemEntity.Price,
					BilledType:       serviceLineItemEntity.Billed.String(),
					Comment:          "quantity changed is " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " for service " + name,
					ReasonForChange:  utils.IfNotNilString(dataFields.Comments),
					Currency:         contractEntity.Currency.String(),
				}
				if dataFields.StartedAt != nil {
					actionPriceMetadata.StartedAt = dataFields.StartedAt
					actionQuantityMetadata.StartedAt = dataFields.StartedAt
				}
				metadataPrice, err := utils.ToJson(actionPriceMetadata)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Failed to serialize price metadata: %s", err.Error())
					return errors.Wrap(err, "Failed to serialize price metadata")
				}
				metadataQuantity, err := utils.ToJson(actionQuantityMetadata)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Failed to serialize quantity metadata: %s", err.Error())
					return errors.Wrap(err, "Failed to serialize quantity metadata")
				}
				oldCycle := getBillingCycleNamingConvention(serviceLineItemEntity.Billed.String())

				if priceChanged && dataFields.BilledType != nil && dataFields.BilledType.IsRecurrent() {
					message := ""
					if utils.IfNotNilFloat64(dataFields.Price) > serviceLineItemEntity.Price {
						message = userName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", utils.IfNotNilFloat64(dataFields.Price)) + "/" + cycle
					}
					if utils.IfNotNilFloat64(dataFields.Price) < serviceLineItemEntity.Price {
						message = userName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + "/" + oldCycle + " to " + fmt.Sprintf("%.2f", utils.IfNotNilFloat64(dataFields.Price)) + "/" + cycle
					}
					_, err = s.neo4j.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, enum.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractEntity.Id, err.Error())
					}
				}

				if priceChanged && dataFields.BilledType != nil && *dataFields.BilledType == neo4jenum.BilledTypeOnce {
					message := ""
					if utils.IfNotNilFloat64(dataFields.Price) > serviceLineItemEntity.Price {
						message = userName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.2f", utils.IfNotNilFloat64(dataFields.Price))
					}
					if utils.IfNotNilFloat64(dataFields.Price) < serviceLineItemEntity.Price {
						message = userName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.2f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.2f", utils.IfNotNilFloat64(dataFields.Price))
					}
					_, err = s.neo4j.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, enum.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractEntity.Id, err.Error())
					}
				}
				if priceChanged && *dataFields.BilledType == neo4jenum.BilledTypeUsage {
					message := ""
					if utils.IfNotNilFloat64(dataFields.Price) > serviceLineItemEntity.Price {
						message = userName + " retroactively increased the price for " + name + " from " + fmt.Sprintf("%.4f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.4f", utils.IfNotNilFloat64(dataFields.Price))
					}
					if utils.IfNotNilFloat64(dataFields.Price) < serviceLineItemEntity.Price {
						message = userName + " retroactively decreased the price for " + name + " from " + fmt.Sprintf("%.4f", serviceLineItemEntity.Price) + " to " + fmt.Sprintf("%.4f", utils.IfNotNilFloat64(dataFields.Price))
					}
					_, err = s.neo4j.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, enum.ActionServiceLineItemPriceUpdated, message, metadataPrice, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("Failed creating price update action for contract service line item %s: %s", contractEntity.Id, err.Error())
					}
				}

				if quantityChanged {
					message := ""
					if utils.IfNotNilInt64(dataFields.Quantity) > serviceLineItemEntity.Quantity {
						message = userName + " retroactively increased the quantity of " + name + " from " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " to " + strconv.FormatInt(utils.IfNotNilInt64(dataFields.Quantity), 10)
					}
					if utils.IfNotNilInt64(dataFields.Quantity) < serviceLineItemEntity.Quantity {
						message = userName + " retroactively decreased the quantity of " + name + " from " + strconv.FormatInt(serviceLineItemEntity.Quantity, 10) + " to " + strconv.FormatInt(utils.IfNotNilInt64(dataFields.Quantity), 10)
					}
					_, err = s.neo4j.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, enum.ActionServiceLineItemQuantityUpdated, message, metadataQuantity, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
					if err != nil {
						tracing.TraceErr(span, err)
						s.log.Errorf("Failed creating quantity update action for contract service line item %s: %s", contractEntity.Id, err.Error())
					}
				}

				err = s.events.Publisher.PublishEvent(ctx, sliId, model.SERVICE_LINE_ITEM, dto.UpdateServiceLineItem{dataFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateServiceLineItem for SLI"))
				}
				err = s.events.Publisher.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.UpdateServiceLineItem{dataFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateServiceLineItem for Contract"))
				}

				s.events.Publisher.PublishEventCompleted(ctx, tenant, sliId, model.SERVICE_LINE_ITEM, utils.NewEventCompletedDetails().WithUpdate())
				s.events.Publisher.PublishEventCompleted(ctx, tenant, contractEntity.Id, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createFlow {
		span.LogFields(log.Bool("response.commentCreated", true))
	} else {
		span.LogFields(log.Bool("response.commentCreated", true))
	}

	return sliId, nil
}

func (s *serviceLineItemService) GetById(ctx context.Context, serviceLineItemId string) (*neo4jentity.ServiceLineItemEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	if sliDbNode, err := s.neo4j.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, common.GetTenantFromContext(ctx), serviceLineItemId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("service line item with id {%s} not found", serviceLineItemId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode), nil
	}
}

func (s *serviceLineItemService) GetServiceLineItemsByParentId(ctx context.Context, sliParentId string) (*neo4jentity.ServiceLineItemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetServiceLineItemsByParentId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("sliParentId", sliParentId))

	serviceLineItems, err := s.neo4j.ServiceLineItemReadRepository.GetServiceLineItemsByParentId(ctx, common.GetTenantFromContext(ctx), sliParentId)
	if err != nil {
		return nil, err
	}
	serviceLineItemEntities := make(neo4jentity.ServiceLineItemEntities, 0, len(serviceLineItems))
	for _, v := range serviceLineItems {
		serviceLineItemEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(v)
		serviceLineItemEntities = append(serviceLineItemEntities, *serviceLineItemEntity)
	}
	return &serviceLineItemEntities, nil
}

func (s *serviceLineItemService) GetServiceLineItemsForContract(ctx context.Context, contractId string) (*neo4jentity.ServiceLineItemEntities, error) {
	return s.GetServiceLineItemsForContracts(ctx, []string{contractId})
}

func (s *serviceLineItemService) GetServiceLineItemsForContracts(ctx context.Context, contractIDs []string) (*neo4jentity.ServiceLineItemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetServiceLineItemsForContracts")
	defer span.Finish()
	span.LogFields(log.Object("contractIDs", contractIDs))

	serviceLineItems, err := s.neo4j.ServiceLineItemReadRepository.GetServiceLineItemsForContracts(ctx, common.GetTenantFromContext(ctx), contractIDs)
	if err != nil {
		return nil, err
	}
	serviceLineItemEntities := make(neo4jentity.ServiceLineItemEntities, 0, len(serviceLineItems))
	for _, v := range serviceLineItems {
		serviceLineItemEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(v.Node)
		serviceLineItemEntity.DataloaderKey = v.LinkedNodeId
		serviceLineItemEntities = append(serviceLineItemEntities, *serviceLineItemEntity)
	}
	return &serviceLineItemEntities, nil
}

func (s *serviceLineItemService) GetServiceLineItemsForInvoiceLines(ctx context.Context, invoiceLineIds []string) (*neo4jentity.ServiceLineItemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetServiceLineItemsForInvoiceLines")
	defer span.Finish()
	span.LogFields(log.Object("invoiceLineIds", invoiceLineIds))

	serviceLineItems, err := s.neo4j.ServiceLineItemReadRepository.GetServiceLineItemsForInvoiceLines(ctx, common.GetTenantFromContext(ctx), invoiceLineIds)
	if err != nil {
		return nil, err
	}
	serviceLineItemEntities := make(neo4jentity.ServiceLineItemEntities, 0, len(serviceLineItems))
	for _, v := range serviceLineItems {
		serviceLineItemEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(v.Node)
		serviceLineItemEntity.DataloaderKey = v.LinkedNodeId
		serviceLineItemEntities = append(serviceLineItemEntities, *serviceLineItemEntity)
	}
	return &serviceLineItemEntities, nil
}

func (s *serviceLineItemService) Pause(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Pause")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, serviceLineItemId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate SLI exists
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, serviceLineItemId, model.NodeLabelServiceLineItem)
	if err != nil || !exists {
		err = errors.New("service line item not found")
		tracing.TraceErr(span, err)
		return err
	}

	// get contract for service line item
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractByServiceLineItemId(ctx, tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		err := s.neo4j.CommonWriteRepository.UpdateBoolProperty(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelServiceLineItem, serviceLineItemId, string(neo4jentity.SLIPropertyPaused), true)
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err := s.events.Publisher.PublishEvent(ctx, serviceLineItemId, model.SERVICE_LINE_ITEM, dto.PauseServiceLineItem{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message PauseServiceLineItem"))
			}
			err = s.events.Publisher.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.PauseServiceLineItem{ServiceLineItemId: serviceLineItemId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message PauseServiceLineItem for contract"))
			}

			s.events.Publisher.PublishEventCompleted(ctx, tenant, contractEntity.Id, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())

			return nil
		})
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *serviceLineItemService) Resume(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Resume")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, serviceLineItemId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate SLI exists
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, serviceLineItemId, model.NodeLabelServiceLineItem)
	if err != nil || !exists {
		err = errors.New("service line item not found")
		tracing.TraceErr(span, err)
		return err
	}

	// get contract for service line item
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractByServiceLineItemId(ctx, tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		err := s.neo4j.CommonWriteRepository.UpdateBoolProperty(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelServiceLineItem, serviceLineItemId, string(neo4jentity.SLIPropertyPaused), false)
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err := s.events.Publisher.PublishEvent(ctx, serviceLineItemId, model.SERVICE_LINE_ITEM, dto.ResumeServiceLineItem{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message ResumeServiceLineItem"))
			}
			err = s.events.Publisher.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.ResumeServiceLineItem{ServiceLineItemId: serviceLineItemId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message ResumeServiceLineItem for contract"))
			}

			s.events.Publisher.PublishEventCompleted(ctx, tenant, contractEntity.Id, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())

			return nil
		})
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func getBillingCycleNamingConvention(billedType string) string {
	switch billedType {
	case neo4jenum.BilledTypeAnnually.String():
		return "year"
	case neo4jenum.BilledTypeQuarterly.String():
		return "quarter"
	case neo4jenum.BilledTypeMonthly.String():
		return "month"
	default:
		return ""
	}
}

func (s *serviceLineItemService) Delete(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Delete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, serviceLineItemId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// get contract for service line item
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractByServiceLineItemId(ctx, tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	serviceLineItemEntity, err := s.GetById(ctx, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		err := s.neo4j.ServiceLineItemWriteRepository.Delete(ctx, txWithPostCommit.Tx, tenant, serviceLineItemId)
		if err != nil {
			return nil, err
		}

		err = s.neo4j.ServiceLineItemWriteRepository.AdjustEndDates(ctx, txWithPostCommit.Tx, tenant, serviceLineItemEntity.ParentID)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while adjusting end dates for service line item %s: %s", serviceLineItemEntity.ParentID, err.Error())
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.contract.UpdateActiveRenewalOpportunityArr(ctx, contractEntity.Id)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			err = s.contract.RecalculateContractLtv(ctx, contractEntity.Id)
			if err != nil {
				tracing.TraceErr(span, err)
			}

			serviceLineItemName := "Unnamed service"
			if serviceLineItemEntity.Name != "" {
				serviceLineItemName = serviceLineItemEntity.Name
			}
			userName := ""
			userDbNode, err := s.neo4j.UserReadRepository.GetUserById(ctx, tenant, common.GetUserIdFromContext(ctx))
			if err != nil {
				tracing.TraceErr(span, err)
			}
			if userDbNode != nil {
				userEntity := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
				userName = userEntity.GetFullName()
			}
			contractName := "Unnamed contract"
			if contractEntity.Name != "" {
				contractName = contractEntity.Name
			}

			metadata, err := utils.ToJson(SLIActionMetadata{
				UserName:    userName,
				ServiceName: serviceLineItemName,
				Comment:     "service line item removed is " + serviceLineItemName + " from " + contractName + " by " + userName,
			})
			message := userName + " removed " + serviceLineItemName + " from " + contractName

			_, err = s.neo4j.ActionWriteRepository.Create(ctx, tenant, contractEntity.Id, model.CONTRACT, enum.ActionServiceLineItemRemoved, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Failed remove service line item action for contract %s: %s", contractEntity.Id, err.Error())
			}

			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err := s.events.Publisher.PublishEvent(ctx, serviceLineItemId, model.SERVICE_LINE_ITEM, dto.DeleteServiceLineItem{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message DeleteServiceLineItem"))
			}
			err = s.events.Publisher.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.DeleteServiceLineItem{ServiceLineItemId: serviceLineItemId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message DeleteServiceLineItem for contract"))
			}

			s.events.Publisher.PublishEventCompleted(ctx, tenant, serviceLineItemId, model.SERVICE_LINE_ITEM, utils.NewEventCompletedDetails().WithDelete())
			s.events.Publisher.PublishEventCompleted(ctx, tenant, contractEntity.Id, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *serviceLineItemService) Close(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string, endedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Close")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, serviceLineItemId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// get contract for service line item
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractByServiceLineItemId(ctx, tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	serviceLineItemEntity, err := s.GetById(ctx, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if serviceLineItemEntity.StartedAt.After(utils.Now()) {
		return s.Delete(ctx, txWithPostCommit, serviceLineItemId)
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		err := s.neo4j.ServiceLineItemWriteRepository.Close(ctx, txWithPostCommit.Tx, tenant, serviceLineItemId, endedAt, true)
		if err != nil {
			return nil, err
		}

		err = s.neo4j.ServiceLineItemWriteRepository.AdjustEndDates(ctx, txWithPostCommit.Tx, tenant, serviceLineItemEntity.ParentID)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while adjusting end dates for service line item %s: %s", serviceLineItemEntity.ParentID, err.Error())
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.contract.UpdateActiveRenewalOpportunityArr(ctx, contractEntity.Id)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			err = s.contract.RecalculateContractLtv(ctx, contractEntity.Id)
			if err != nil {
				tracing.TraceErr(span, err)
			}

			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err := s.events.Publisher.PublishEvent(ctx, serviceLineItemId, model.SERVICE_LINE_ITEM, dto.CloseServiceLineItem{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CloseServiceLineItem"))
			}
			err = s.events.Publisher.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.CloseServiceLineItem{ServiceLineItemId: serviceLineItemId, EndedAt: endedAt})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CloseServiceLineItem for contract"))
			}

			s.events.Publisher.PublishEventCompleted(ctx, tenant, serviceLineItemId, model.SERVICE_LINE_ITEM, utils.NewEventCompletedDetails().WithUpdate())
			s.events.Publisher.PublishEventCompleted(ctx, tenant, contractEntity.Id, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
