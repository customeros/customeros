package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	location_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	models_common "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	grpcErrors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

type locationService struct {
	location_grpc_service.UnimplementedLocationGrpcServiceServer
	log              logger.Logger
	repositories     *repository.Repositories
	locationCommands *commands.LocationCommands
}

func NewLocationService(log logger.Logger, repositories *repository.Repositories, locationCommands *commands.LocationCommands) *locationService {
	return &locationService{
		log:              log,
		repositories:     repositories,
		locationCommands: locationCommands,
	}
}

func (s *locationService) UpsertLocation(ctx context.Context, request *location_grpc_service.UpsertLocationGrpcRequest) (*location_grpc_service.LocationIdGrpcResponse, error) {
	objectID := request.Id
	utils.TimestampProtoToTime(request.CreatedAt)

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

	source := models_common.Source{
		Source:        request.Source,
		SourceOfTruth: request.SourceOfTruth,
		AppSource:     request.AppSource,
	}

	command := commands.NewUpsertLocationCommand(objectID, request.Tenant, request.Name, request.RawAddress, addressFields, source, utils.TimestampProtoToTime(request.CreatedAt), utils.TimestampProtoToTime(request.UpdatedAt))
	if err := s.locationCommands.UpsertLocation.Handle(ctx, command); err != nil {
		s.log.Errorf("(UpsertLocation.Handle) tenant:{%s}, location ID: {%s}, err: {%v}", request.Tenant, objectID, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(Upserted location): {%s}", objectID)

	return &location_grpc_service.LocationIdGrpcResponse{Id: objectID}, nil
}

func (s *locationService) errResponse(err error) error {
	return grpcErrors.ErrResponse(err)
}
