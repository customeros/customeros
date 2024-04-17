package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ServiceLineItemService interface {
	GetById(ctx context.Context, id string) (*neo4jentity.ServiceLineItemEntity, error)
	GetServiceLineItemsByParentId(ctx context.Context, sliParentId string) (*neo4jentity.ServiceLineItemEntities, error)
	GetServiceLineItemsForContract(ctx context.Context, contractId string) (*neo4jentity.ServiceLineItemEntities, error)
	GetServiceLineItemsForContracts(ctx context.Context, contractIds []string) (*neo4jentity.ServiceLineItemEntities, error)
}

type serviceLineItemService struct {
	log      logger.Logger
	services *Services
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
