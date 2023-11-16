package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	servicelineitempb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/service_line_item"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type ServiceLineItemService interface {
	Create(ctx context.Context, serviceLineItem *ServiceLineItemCreateData) (string, error)
	GetById(ctx context.Context, id string) (*entity.ServiceLineItemEntity, error)
	GetServiceLineItemsForContracts(ctx context.Context, contractIds []string) (*entity.ServiceLineItemEntities, error)
}
type serviceLineItemService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewServiceLineItemService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) ServiceLineItemService {
	return &serviceLineItemService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

type ServiceLineItemCreateData struct {
	ServiceLineItemEntity *entity.ServiceLineItemEntity
	ContractId            string
	ExternalReference     *entity.ExternalSystemEntity
	Source                entity.DataSource
	AppSource             string
}

func (s *serviceLineItemService) Create(ctx context.Context, serviceLineItemDetails *ServiceLineItemCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("serviceLineItemDetails", serviceLineItemDetails))

	if serviceLineItemDetails.ServiceLineItemEntity == nil {
		err := fmt.Errorf("service line item entity is nil")
		tracing.TraceErr(span, err)
		return "", err
	}

	serviceLineItemId, err := s.createServiceLineItemWithEvents(ctx, serviceLineItemDetails)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	span.LogFields(log.String("output - createdServiceLineItemId", serviceLineItemId))
	return serviceLineItemId, nil
}

func (s *serviceLineItemService) createServiceLineItemWithEvents(ctx context.Context, serviceLineItemDetails *ServiceLineItemCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.createServiceLineItemWithEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	createServiceLineItemRequest := servicelineitempb.CreateServiceLineItemGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		ContractId:     serviceLineItemDetails.ContractId,
		Name:           serviceLineItemDetails.ServiceLineItemEntity.Name,
		Quantity:       int64(serviceLineItemDetails.ServiceLineItemEntity.Quantity),
		Price:          float32(serviceLineItemDetails.ServiceLineItemEntity.Price),
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    string(serviceLineItemDetails.Source),
			AppSource: utils.StringFirstNonEmpty(serviceLineItemDetails.AppSource, constants.AppSourceCustomerOsApi),
		},
	}

	switch serviceLineItemDetails.ServiceLineItemEntity.Billed {
	case entity.BilledTypeMonthly:
		createServiceLineItemRequest.Billed = servicelineitempb.BilledType_MONTHLY_BILLED
	case entity.BilledTypeAnnually:
		createServiceLineItemRequest.Billed = servicelineitempb.BilledType_ANNUALLY_BILLED
	case entity.BilledTypeOnce:
		createServiceLineItemRequest.Billed = servicelineitempb.BilledType_ONCE_BILLED
	default:
		createServiceLineItemRequest.Billed = servicelineitempb.BilledType_MONTHLY_BILLED
	}

	response, err := s.grpcClients.ServiceLineItemClient.CreateServiceLineItem(ctx, &createServiceLineItemRequest)

	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4jAfterEventRequest; i++ {
		serviceLineItemFound, findErr := s.repositories.CommonRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), response.Id, entity.NodeLabel_ServiceLineItem)
		if serviceLineItemFound && findErr == nil {
			span.LogFields(log.Bool("serviceLineItemSavedInGraphDb", true))
			break
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}
	return response.Id, err
}

func (s *serviceLineItemService) GetById(ctx context.Context, serviceLineItemId string) (*entity.ServiceLineItemEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("serviceLineItemId", serviceLineItemId))

	if serviceLineItemDbNode, err := s.repositories.ServiceLineItemRepository.GetById(ctx, common.GetContext(ctx).Tenant, serviceLineItemId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("service line item with id {%s} not found", serviceLineItemId))
		return nil, wrappedErr
	} else {
		return s.mapDbNodeToServiceLineItemEntity(*serviceLineItemDbNode), nil
	}
}

func (s *serviceLineItemService) GetServiceLineItemsForContracts(ctx context.Context, contractIDs []string) (*entity.ServiceLineItemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemService.GetServiceLineItemsForContracts")
	defer span.Finish()
	span.LogFields(log.Object("contractIDs", contractIDs))

	serviceLineItems, err := s.repositories.ServiceLineItemRepository.GetForContracts(ctx, common.GetTenantFromContext(ctx), contractIDs)
	if err != nil {
		return nil, err
	}
	serviceLineItemEntities := make(entity.ServiceLineItemEntities, 0, len(serviceLineItems))
	for _, v := range serviceLineItems {
		serviceLineItemEntity := s.mapDbNodeToServiceLineItemEntity(*v.Node)
		serviceLineItemEntity.DataloaderKey = v.LinkedNodeId
		serviceLineItemEntities = append(serviceLineItemEntities, *serviceLineItemEntity)
	}
	return &serviceLineItemEntities, nil
}

func (s *serviceLineItemService) mapDbNodeToServiceLineItemEntity(dbNode dbtype.Node) *entity.ServiceLineItemEntity {
	props := utils.GetPropsFromNode(dbNode)
	serviceLineItem := entity.ServiceLineItemEntity{
		ID:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Billed:        entity.GetBilledType(utils.GetStringPropOrEmpty(props, "billed")),
		Price:         utils.GetFloatPropOrZero(props, "price"),
		Quantity:      utils.GetInt64PropOrZero(props, "quantity"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &serviceLineItem
}
