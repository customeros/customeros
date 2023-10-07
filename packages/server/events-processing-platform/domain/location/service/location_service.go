package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	location_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	grpcErrors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
)

type locationService struct {
	location_grpc_service.UnimplementedLocationGrpcServiceServer
	log              logger.Logger
	repositories     *repository.Repositories
	locationCommands *command_handler.LocationCommands
}

func NewLocationService(log logger.Logger, repositories *repository.Repositories, locationCommands *command_handler.LocationCommands) *locationService {
	return &locationService{
		log:              log,
		repositories:     repositories,
		locationCommands: locationCommands,
	}
}

func (s *locationService) UpsertLocation(ctx context.Context, request *location_grpc_service.UpsertLocationGrpcRequest) (*location_grpc_service.LocationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "LocationService.UpsertLocation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)

	locationId := request.Id
	locationId = utils.NewUUIDIfEmpty(locationId)

	sourceFields := common_models.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.Source = utils.StringFirstNonEmpty(sourceFields.Source, request.Source)
	sourceFields.SourceOfTruth = utils.StringFirstNonEmpty(sourceFields.SourceOfTruth, request.SourceOfTruth)
	sourceFields.AppSource = utils.StringFirstNonEmpty(sourceFields.AppSource, request.AppSource)

	addressFields := models.LocationAddressFields{
		Country:      request.Country,
		Region:       request.Region,
		District:     request.District,
		Locality:     request.Locality,
		Street:       request.Street,
		Address1:     request.AddressLine1,
		Address2:     request.AddressLine2,
		Zip:          request.ZipCode,
		AddressType:  request.AddressType,
		HouseNumber:  request.HouseNumber,
		PostalCode:   request.PostalCode,
		PlusFour:     request.PlusFour,
		Commercial:   request.Commercial,
		Predirection: request.Predirection,
		Latitude:     utils.ParseStringToFloat(request.Latitude),
		Longitude:    utils.ParseStringToFloat(request.Longitude),
	}

	cmd := command.NewUpsertLocationCommand(locationId, request.Tenant, request.LoggedInUserId, request.Name, request.RawAddress, addressFields, sourceFields,
		utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.locationCommands.UpsertLocation.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(UpsertLocation.Handle) tenant:{%s}, location id: {%s}, err: {%v}", request.Tenant, locationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(Upserted location): {%s}", locationId)

	return &location_grpc_service.LocationIdGrpcResponse{Id: locationId}, nil
}

func (s *locationService) errResponse(err error) error {
	return grpcErrors.ErrResponse(err)
}
