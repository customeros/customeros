package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OrderService interface {
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*entity.OrderEntities, error)
	mapDbNodeToOrderEntity(node *dbtype.Node) *entity.OrderEntity
}

type orderService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewOrderService(log logger.Logger, repositories *repository.Repositories) OrderService {
	return &orderService{
		log:          log,
		repositories: repositories,
	}
}

func (s *orderService) GetAllForOrganizations(ctx context.Context, organizationIds []string) (*entity.OrderEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderService.GetAllForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("linkedEntityIds", organizationIds))

	dbNodeList, err := s.repositories.Neo4jRepositories.OrderReadRepository.GetAllForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}

	entityList := make(entity.OrderEntities, 0, len(dbNodeList))

	for _, v := range dbNodeList {
		entity := s.mapDbNodeToOrderEntity(v.Node)
		entity.DataloaderKey = v.LinkedNodeId
		entityList = append(entityList, *entity)
	}

	return &entityList, nil
}

func (s *orderService) mapDbNodeToOrderEntity(node *dbtype.Node) *entity.OrderEntity {
	props := utils.GetPropsFromNode(*node)
	return &entity.OrderEntity{
		Id:          utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:   utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:   utils.GetTimePropOrEpochStart(props, "updatedAt"),
		ConfirmedAt: utils.GetTimePropOrNil(props, "confirmedAt"),
		PaidAt:      utils.GetTimePropOrNil(props, "paidAt"),
		FulfilledAt: utils.GetTimePropOrNil(props, "fulfilledAt"),
		CancelledAt: utils.GetTimePropOrNil(props, "cancelledAt"),
		SourceFields: entity.SourceFields{
			Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
			SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
			AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		},
	}
}
