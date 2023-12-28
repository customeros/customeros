package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	grpcerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
)

type locationService struct {
	locationpb.UnimplementedLocationGrpcServiceServer
	log              logger.Logger
	locationCommands *command_handler.CommandHandlers
}

func NewLocationService(log logger.Logger, locationCommands *command_handler.CommandHandlers) *locationService {
	return &locationService{
		log:              log,
		locationCommands: locationCommands,
	}
}

func (s *locationService) UpsertLocation(ctx context.Context, request *locationpb.UpsertLocationGrpcRequest) (*locationpb.LocationIdGrpcResponse, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "LocationService.UpsertLocation")
	defer span.Finish()
	tracing.SetServiceSpanTags(ctx, span, request.Tenant, request.LoggedInUserId)
	tracing.LogObjectAsJson(span, "request", request)

	locationId := request.Id
	locationId = utils.NewUUIDIfEmpty(locationId)

	sourceFields := commonmodel.Source{}
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
		utils.TimestampProtoToTimePtr(request.CreatedAt), utils.TimestampProtoToTimePtr(request.UpdatedAt))
	if err := s.locationCommands.UpsertLocation.Handle(ctx, cmd); err != nil {
		s.log.Errorf("(UpsertLocation.Handle) tenant:{%s}, location id: {%s}, err: {%v}", request.Tenant, locationId, err)
		return nil, s.errResponse(err)
	}

	s.log.Infof("(Upserted location): {%s}", locationId)

	return &locationpb.LocationIdGrpcResponse{Id: locationId}, nil
}

func (s *locationService) errResponse(err error) error {
	return grpcerr.ErrResponse(err)
}
