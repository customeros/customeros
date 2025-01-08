package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LocationService interface {
	GetById(ctx context.Context, locationId string) (*neo4jentity.LocationEntity, error)
}

type locationService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewLocationService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) LocationService {
	return &locationService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *locationService) GetById(ctx context.Context, locationId string) (*neo4jentity.LocationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("locationId", locationId))

	locationNode, err := s.repositories.LocationRepository.GetById(ctx, locationId)
	if err != nil {
		return nil, err
	}

	return neo4jmapper.MapDbNodeToLocationEntity(locationNode), nil
}
