package service

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
)

type LocationService interface {
	GetById(ctx context.Context, locationId string) (*neo4jentity.LocationEntity, error)
}

type locationService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewLocationService(log logger.Logger, repositories *repository.Repositories) LocationService {
	return &locationService{
		log:          log,
		repositories: repositories,
	}
}

func (s *locationService) GetById(ctx context.Context, locationId string) (*neo4jentity.LocationEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "LocationService.GetById")
	defer spans.Finish()
	spans.LogKV("locationId", locationId)

	locationNode, err := s.repositories.LocationRepository.GetById(ctx, locationId)
	if err != nil {
		return nil, err
	}

	return neo4jmapper.MapDbNodeToLocationEntity(locationNode), nil
}
