package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type HealthIndicatorService interface {
	GetAll(ctx context.Context) (*entity.HealthIndicatorEntities, error)
}

type healthIndicatorService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewHealthIndicatorService(log logger.Logger, repository *repository.Repositories) HealthIndicatorService {
	return &healthIndicatorService{
		log:          log,
		repositories: repository,
	}
}

func (s *healthIndicatorService) GetAll(ctx context.Context) (*entity.HealthIndicatorEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HealthIndicatorService.GetAll")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	healthIndicatorDbNodes, err := s.repositories.HealthIndicatorRepository.GetAll(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return nil, err
	}
	healthIndicatorEntities := make(entity.HealthIndicatorEntities, 0, len(healthIndicatorDbNodes))
	for _, dbNodePtr := range healthIndicatorDbNodes {
		healthIndicatorEntities = append(healthIndicatorEntities, *s.mapDbNodeToHealthIndicatorEntity(*dbNodePtr))
	}
	return &healthIndicatorEntities, nil
}

func (s *healthIndicatorService) mapDbNodeToHealthIndicatorEntity(node dbtype.Node) *entity.HealthIndicatorEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.HealthIndicatorEntity{
		Id:    utils.GetStringPropOrEmpty(props, "id"),
		Name:  utils.GetStringPropOrEmpty(props, "name"),
		Order: utils.GetInt64PropOrZero(props, "order"),
	}
}
