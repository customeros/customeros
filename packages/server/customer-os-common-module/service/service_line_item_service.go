package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strconv"
	"time"
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

type ServiceLineItemService interface {
	GetById(ctx context.Context, id string) (*neo4jentity.ServiceLineItemEntity, error)
	GetServiceLineItemsByParentId(ctx context.Context, sliParentId string) (*neo4jentity.ServiceLineItemEntities, error)
	GetServiceLineItemsForContract(ctx context.Context, contractId string) (*neo4jentity.ServiceLineItemEntities, error)
	GetServiceLineItemsForContracts(ctx context.Context, contractIds []string) (*neo4jentity.ServiceLineItemEntities, error)
	GetServiceLineItemsForInvoiceLines(ctx context.Context, invoiceLineIds []string) (*neo4jentity.ServiceLineItemEntities, error)
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, dataFields data_fields.SLIFields) (string, error)
	PauseServiceLineItem(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error
	ResumeServiceLineItem(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error
}

type serviceLineItemService struct {
	log      logger.Logger
	services *Services
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
			exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, *dataFields.ContractId, model.NodeLabelContract)
			if err != nil || !exists {
				err = errors.New("contract not found")
				tracing.TraceErr(span, err)
				return "", err
			}
		}

		// generate id
		sliId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelServiceLineItem)
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

		// validate service line item exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, sliId, model.NodeLabelServiceLineItem)
		if err != nil || !exists {
			err = errors.New("comment not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, sliId)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		if createFlow {
			// created SLI in neo4j
			err := s.services.Neo4jRepositories.ServiceLineItemWriteRepository.CreateForContract(ctx, txWithPostCommit.Tx, tenant, sliId, dataFields)
			if err != nil {
				s.log.Errorf("error creating service line item %s: %s", sliId, err.Error())
				return nil, err
			}
			err = s.services.Neo4jRepositories.ServiceLineItemWriteRepository.AdjustEndDates(ctx, txWithPostCommit.Tx, tenant, *dataFields.ParentId)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error while adjusting end dates for service line item %s: %s", sliId, err.Error())
				return nil, err
			}

		} else {
			// TODO implement update
			//err := s.services.Neo4jRepositories.CommentWriteRepository.Update(ctx, txWithPostCommit.Tx, tenant, sliId, commentFields)
			//if err != nil {
			//	s.log.Errorf("Error while updating service line item %s: %s", sliId, err.Error())
			//	return nil, err
			//}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if createFlow {
				err = s.services.ContractService.UpdateActiveRenewalOpportunityArr(ctx, *dataFields.ContractId)
				if err != nil {
					tracing.TraceErr(span, err)
				}
				err = s.services.ContractService.RecalculateContractLtv(ctx, *dataFields.ContractId)
				if err != nil {
					tracing.TraceErr(span, err)
				}
				contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, *dataFields.ContractId)
				if err != nil {
					tracing.TraceErr(span, err)
				}
				contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

				if dataFields.BilledType != nil && utils.IfNotNilString(dataFields.BilledType.String()) != "" {
					name := "Unnamed service"
					if utils.IfNotNilString(dataFields.Name) != "" {
						name = *dataFields.Name
					}

					userName := ""
					userDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetUserById(ctx, tenant, common.GetUserIdFromContext(ctx))
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
					}
					if err == nil {
						if *dataFields.BilledType == neo4jenum.BilledTypeAnnually || *dataFields.BilledType == neo4jenum.BilledTypeQuarterly || *dataFields.BilledType == neo4jenum.BilledTypeMonthly {
							message := userName + " added a recurring service to " + contractEntity.Name + ": " + name + " at " + strconv.FormatInt(utils.IfNotNilInt64(dataFields.Quantity), 10) + " x " + fmt.Sprintf("%.2f", utils.IfNotNilFloat64(dataFields.Price)) + "/" + cycle + " starting with " + dataFields.StartedAt.Format("2006-01-02")
							_, err = s.services.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, neo4jenum.ActionServiceLineItemBilledTypeRecurringCreated, message, metadataBilledType, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
							if err != nil {
								tracing.TraceErr(span, err)
								s.log.Errorf("Failed creating recurring billed type service line item created action for contract %s: %s", contractEntity.Id, err.Error())
							}
						}
						if *dataFields.BilledType == neo4jenum.BilledTypeOnce {
							message := userName + " added a one time service to " + contractEntity.Name + ": " + name + " at " + fmt.Sprintf("%.2f", utils.IfNotNilFloat64(dataFields.Price)) + " starting with " + dataFields.StartedAt.Format("2006-01-02")
							_, err = s.services.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, neo4jenum.ActionServiceLineItemBilledTypeOnceCreated, message, metadataBilledType, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
							if err != nil {
								tracing.TraceErr(span, err)
								s.log.Errorf("Failed creating once billed type service line item created action for contract %s: %s", contractEntity.Id, err.Error())
							}
						}
						if *dataFields.BilledType == neo4jenum.BilledTypeUsage {
							message := userName + " added a per use service to " + contractEntity.Name + ": " + name + " at " + fmt.Sprintf("%.4f", utils.IfNotNilFloat64(dataFields.Price)) + " starting with " + dataFields.StartedAt.Format("2006-01-02")
							_, err = s.services.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, tenant, contractEntity.Id, model.CONTRACT, neo4jenum.ActionServiceLineItemBilledTypeUsageCreated, message, metadataBilledType, utils.Now(), common.GetAppSourceFromContext(ctx), extraActionProperties)
							if err != nil {
								tracing.TraceErr(span, err)
								s.log.Errorf("Failed creating per use billed type service line item created action for contract %s: %s", contractEntity.Id, err.Error())
							}
						}
					}
				}

				err = s.services.RabbitMQService.PublishEvent(ctx, sliId, model.SERVICE_LINE_ITEM, dto.CreateServiceLineItem{dataFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateServiceLineItem for SLI"))
				}
				err = s.services.RabbitMQService.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.CreateServiceLineItem{dataFields})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateServiceLineItem for Contract"))
				}

				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, sliId, model.SERVICE_LINE_ITEM, utils.NewEventCompletedDetails().WithCreate())
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, contractEntity.Id, model.CONTRACT, utils.NewEventCompletedDetails().WithUpdate())
			} else {

			}
			//	// send events
			//	if createFlow {
			//		err := s.services.RabbitMQService.PublishEvent(ctx, commentId, model.COMMENT, dto.CreateComment{commentFields})
			//		if err != nil {
			//			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateComment"))
			//		}
			//		s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, commentId, model.COMMENT, utils.NewEventCompletedDetails().WithCreate())
			//	} else {
			// TODO implement update
			//		err := s.services.RabbitMQService.PublishEvent(ctx, commentId, model.COMMENT, dto.UpdateComment{commentFields})
			//		if err != nil {
			//			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateComment"))
			//		}
			//		if common.GetTenantFromContext(ctx) != constants.AppSourceCustomerOsApi {
			//			s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, commentId, model.COMMENT, utils.NewEventCompletedDetails().WithUpdate())
			//		}
			//	}
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

func NewServiceLineItemService(log logger.Logger, services *Services) ServiceLineItemService {
	return &serviceLineItemService{
		log:      log,
		services: services,
	}
}

func (s *serviceLineItemService) GetById(ctx context.Context, serviceLineItemId string) (*neo4jentity.ServiceLineItemEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	if sliDbNode, err := s.services.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, common.GetTenantFromContext(ctx), serviceLineItemId); err != nil {
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

	serviceLineItems, err := s.services.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemsByParentId(ctx, common.GetTenantFromContext(ctx), sliParentId)
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

	serviceLineItems, err := s.services.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemsForContracts(ctx, common.GetTenantFromContext(ctx), contractIDs)
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

	serviceLineItems, err := s.services.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemsForInvoiceLines(ctx, common.GetTenantFromContext(ctx), invoiceLineIds)
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

func (s *serviceLineItemService) PauseServiceLineItem(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.PauseServiceLineItem")
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
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, serviceLineItemId, model.NodeLabelServiceLineItem)
	if err != nil || !exists {
		err = errors.New("service line item not found")
		tracing.TraceErr(span, err)
		return err
	}

	// get contract for service line item
	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		err := s.services.Neo4jRepositories.CommonWriteRepository.UpdateBoolProperty(ctx, nil, tenant, model.NodeLabelServiceLineItem, serviceLineItemId, string(neo4jentity.SLIPropertyPaused), true)
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err := s.services.RabbitMQService.PublishEvent(ctx, serviceLineItemId, model.SERVICE_LINE_ITEM, dto.PauseServiceLineItem{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message PauseServiceLineItem"))
			}
			err = s.services.RabbitMQService.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.PauseServiceLineItem{ServiceLineItemId: serviceLineItemId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message PauseServiceLineItem for contract"))
			}
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

func (s *serviceLineItemService) ResumeServiceLineItem(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.ResumeServiceLineItem")
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
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, serviceLineItemId, model.NodeLabelServiceLineItem)
	if err != nil || !exists {
		err = errors.New("service line item not found")
		tracing.TraceErr(span, err)
		return err
	}

	// get contract for service line item
	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractByServiceLineItemId(ctx, tenant, serviceLineItemId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		err := s.services.Neo4jRepositories.CommonWriteRepository.UpdateBoolProperty(ctx, nil, tenant, model.NodeLabelServiceLineItem, serviceLineItemId, string(neo4jentity.SLIPropertyPaused), false)
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err := s.services.RabbitMQService.PublishEvent(ctx, serviceLineItemId, model.SERVICE_LINE_ITEM, dto.ResumeServiceLineItem{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message ResumeServiceLineItem"))
			}
			err = s.services.RabbitMQService.PublishEvent(ctx, contractEntity.Id, model.CONTRACT, dto.ResumeServiceLineItem{ServiceLineItemId: serviceLineItemId})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message ResumeServiceLineItem for contract"))
			}
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
