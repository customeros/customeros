package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	commongrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	locationgrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type LocationService interface {
	GetById(ctx context.Context, locationId string) (*entity.LocationEntity, error)
	CreateLocation(ctx context.Context, locationId, externalSystem, appSource, locationName, country, region, locality, address, address2, zip string) (string, error)
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

func (s *locationService) CreateLocation(ctx context.Context, locationId, externalSystem, appSource string,
	locationName, country, region, locality, address, address2, zip string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationService.CreateLocation")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	response, err := s.grpcClients.LocationClient.UpsertLocation(ctx, &locationgrpc.UpsertLocationGrpcRequest{
		Tenant: common.GetTenantFromContext(ctx),
		Id:     locationId,
		Name:   locationName,
		SourceFields: &commongrpc.SourceFields{
			Source:    externalSystem,
			AppSource: utils.StringFirstNonEmpty(appSource, constants.AppSourceCustomerOsWebhooks),
		},
		RawAddress:   "",
		CreatedAt:    utils.ConvertTimeToTimestampPtr(utils.NowAsPtr()),
		UpdatedAt:    utils.ConvertTimeToTimestampPtr(utils.NowAsPtr()),
		Country:      country,
		Region:       region,
		Locality:     locality,
		AddressLine1: address,
		AddressLine2: address2,
		ZipCode:      zip,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing %s", err.Error())
		return "", err
	}
	// wait for neo4j to finish processing
	for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
		locationEntity, findLocationErr := s.GetById(ctx, response.Id)
		if locationEntity != nil && findLocationErr == nil {
			break
		}
		time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
	}
	span.LogFields(log.String("upsertedLocationId", response.Id))
	return response.Id, nil
}

func (s *locationService) GetById(ctx context.Context, locationId string) (*entity.LocationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("locationId", locationId))

	locationNode, err := s.repositories.LocationRepository.GetById(ctx, locationId)
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToLocationEntity(*locationNode), nil
}

func (s *locationService) mapDbNodeToLocationEntity(node dbtype.Node) *entity.LocationEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.LocationEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
}
