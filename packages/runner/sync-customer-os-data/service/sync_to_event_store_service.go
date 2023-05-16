package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	local_utils "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	emailgrpcservice "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	location_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/location"
	phonenumbergrpcservice "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/phone_number"
	"github.com/sirupsen/logrus"
)

type SyncToEventStoreService interface {
	SyncEmails(ctx context.Context, batchSize int)
	SyncPhoneNumbers(ctx context.Context, batchSize int)
	SyncLocations(ctx context.Context, batchSize int)
}

type syncToEventStoreService struct {
	repositories *repository.Repositories
	services     *Services
	grpcClients  *grpc_client.Clients
}

func NewSyncToEventStoreService(repositories *repository.Repositories, services *Services, grpcClients *grpc_client.Clients) SyncToEventStoreService {
	return &syncToEventStoreService{
		repositories: repositories,
		services:     services,
		grpcClients:  grpcClients,
	}
}

func (s *syncToEventStoreService) SyncEmails(ctx context.Context, batchSize int) {
	logrus.Infof("start sync emails to eventstore at %v", utils.Now())
	completedCount := 0
	failedCount := 0

	completedCount, failedCount, _ = s.upsertEmailsIntoEventStore(ctx, batchSize)

	logrus.Infof("completed %v and faled %v emails upserting to eventstore at %v", completedCount, failedCount, utils.Now())
}

func (s *syncToEventStoreService) upsertEmailsIntoEventStore(ctx context.Context, batchSize int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.EmailRepository.GetAllCrossTenantsWithRawEmail(ctx, batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.EmailClient.UpsertEmail(context.Background(), &emailgrpcservice.UpsertEmailGrpcRequest{
			Id:            utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant:        v.LinkedNodeId,
			RawEmail:      utils.GetStringPropOrEmpty(v.Node.Props, "rawEmail"),
			Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
			SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
			AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
			CreatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
		})
		if err != nil {
			failedRecords++
			logrus.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}

func (s *syncToEventStoreService) SyncPhoneNumbers(ctx context.Context, batchSize int) {
	logrus.Infof("start sync phone numbers to eventstore at %v", utils.Now())
	completedCount := 0
	failedCount := 0

	completedCount, failedCount, _ = s.upsertPhoneNumbersIntoEventStore(ctx, batchSize)

	logrus.Infof("completed %v and faled %v phone numbers upserting to eventstore at %v", completedCount, failedCount, utils.Now())
}

func (s *syncToEventStoreService) upsertPhoneNumbersIntoEventStore(ctx context.Context, batchSize int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.PhoneNumberRepository.GetAllCrossTenantsWithRawPhoneNumber(ctx, batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.PhoneNumberClient.UpsertPhoneNumber(context.Background(), &phonenumbergrpcservice.UpsertPhoneNumberGrpcRequest{
			Id:            utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant:        v.LinkedNodeId,
			PhoneNumber:   utils.GetStringPropOrEmpty(v.Node.Props, "rawPhoneNumber"),
			Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
			SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
			AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
			CreatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
		})
		if err != nil {
			failedRecords++
			logrus.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}

func (s *syncToEventStoreService) SyncLocations(ctx context.Context, batchSize int) {
	logrus.Infof("start sync locations to eventstore at %v", utils.Now())
	completedCount := 0
	failedCount := 0

	completedCount, failedCount, _ = s.upsertLocationsIntoEventStore(ctx, batchSize)

	logrus.Infof("completed %v and faled %v locations upserting to eventstore at %v", completedCount, failedCount, utils.Now())
}

func (s *syncToEventStoreService) upsertLocationsIntoEventStore(ctx context.Context, batchSize int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.LocationRepository.GetAllCrossTenants(ctx, batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.LocationClient.UpsertLocation(context.Background(), &location_grpc_service.UpsertLocationGrpcRequest{
			Id:            utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant:        v.LinkedNodeId,
			Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
			SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
			AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
			CreatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
			Name:          utils.GetStringPropOrEmpty(v.Node.Props, "name"),
			RawAddress:    utils.GetStringPropOrEmpty(v.Node.Props, "rawAddress"),
			Country:       utils.GetStringPropOrEmpty(v.Node.Props, "country"),
			Region:        utils.GetStringPropOrEmpty(v.Node.Props, "region"),
			Locality:      utils.GetStringPropOrEmpty(v.Node.Props, "locality"),
			AddressLine1:  utils.GetStringPropOrEmpty(v.Node.Props, "address"),
			AddressLine2:  utils.GetStringPropOrEmpty(v.Node.Props, "address2"),
			ZipCode:       utils.GetStringPropOrEmpty(v.Node.Props, "zip"),
			AddressType:   utils.GetStringPropOrEmpty(v.Node.Props, "addressType"),
			HouseNumber:   utils.GetStringPropOrEmpty(v.Node.Props, "houseNumber"),
			PostalCode:    utils.GetStringPropOrEmpty(v.Node.Props, "postalCode"),
			PlusFour:      utils.GetStringPropOrEmpty(v.Node.Props, "plusFour"),
			Commercial:    utils.GetBoolPropOrFalse(v.Node.Props, "commercial"),
			Predirection:  utils.GetStringPropOrEmpty(v.Node.Props, "predirection"),
			District:      utils.GetStringPropOrEmpty(v.Node.Props, "district"),
			Street:        utils.GetStringPropOrEmpty(v.Node.Props, "street"),
			Latitude:      local_utils.FloatToString(utils.GetFloatPropOrNil(v.Node.Props, "latitude")),
			Longitude:     local_utils.FloatToString(utils.GetFloatPropOrNil(v.Node.Props, "longitude")),
		})
		if err != nil {
			failedRecords++
			logrus.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}
