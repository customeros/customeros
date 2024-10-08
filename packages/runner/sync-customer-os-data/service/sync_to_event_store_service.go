package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
)

type SyncToEventStoreService interface {
	SyncPhoneNumbers(ctx context.Context, batchSize int)
	SyncLocations(ctx context.Context, batchSize int)
	SyncContacts(ctx context.Context, batchSize int)
	SyncOrganizations(ctx context.Context, batchSize int)
	SyncOrganizationsLinksWithDomains(ctx context.Context, batchSize int)
}

type syncToEventStoreService struct {
	repositories *repository.Repositories
	services     *Services
	grpcClients  *grpc_client.Clients
	log          logger.Logger
}

func NewSyncToEventStoreService(repositories *repository.Repositories, services *Services, grpcClients *grpc_client.Clients, log logger.Logger) SyncToEventStoreService {
	return &syncToEventStoreService{
		repositories: repositories,
		services:     services,
		grpcClients:  grpcClients,
		log:          log,
	}
}

func (s *syncToEventStoreService) upsertEmailsIntoEventStore(ctx context.Context, batchSize int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.EmailRepository.GetAllCrossTenantsWithRawEmail(ctx, batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.EmailClient.UpsertEmail(context.Background(), &emailpb.UpsertEmailGrpcRequest{
			Id:       utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant:   v.LinkedNodeId,
			RawEmail: utils.GetStringPropOrEmpty(v.Node.Props, "rawEmail"),
			SourceFields: &commonpb.SourceFields{
				Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
				SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
				AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
			},
			CreatedAt: utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt: utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
		})
		if err != nil {
			failedRecords++
			s.log.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}

func (s *syncToEventStoreService) SyncPhoneNumbers(ctx context.Context, batchSize int) {
	s.log.Info("start sync phone numbers to eventstore")

	completed, failed, _ := s.upsertPhoneNumbersIntoEventStore(ctx, batchSize)

	s.log.Infof("completed {%d} and failed {%d} phone numbers upserting to eventstore", completed, failed)
}

func (s *syncToEventStoreService) upsertPhoneNumbersIntoEventStore(ctx context.Context, batchSize int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.PhoneNumberRepository.GetAllCrossTenantsWithRawPhoneNumber(ctx, batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.PhoneNumberClient.UpsertPhoneNumber(context.Background(), &phonenumberpb.UpsertPhoneNumberGrpcRequest{
			Id:          utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant:      v.LinkedNodeId,
			PhoneNumber: utils.GetStringPropOrEmpty(v.Node.Props, "rawPhoneNumber"),
			SourceFields: &commonpb.SourceFields{
				Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
				SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
				AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
			},
			CreatedAt: utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt: utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
		})
		if err != nil {
			failedRecords++
			s.log.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}

func (s *syncToEventStoreService) SyncLocations(ctx context.Context, batchSize int) {
	s.log.Info("start sync locations to eventstore")

	completed, failed, _ := s.upsertLocationsIntoEventStore(ctx, batchSize)

	s.log.Infof("completed {%d} and failed {%d} locations upserting to eventstore", completed, failed)
}

func (s *syncToEventStoreService) upsertLocationsIntoEventStore(ctx context.Context, batchSize int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.LocationRepository.GetAllCrossTenants(ctx, batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.LocationClient.UpsertLocation(context.Background(), &locationpb.UpsertLocationGrpcRequest{
			Id:     utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant: v.LinkedNodeId,
			SourceFields: &commonpb.SourceFields{
				Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
				SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
				AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
			},
			CreatedAt:    utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt:    utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
			Name:         utils.GetStringPropOrEmpty(v.Node.Props, "name"),
			RawAddress:   utils.GetStringPropOrEmpty(v.Node.Props, "rawAddress"),
			Country:      utils.GetStringPropOrEmpty(v.Node.Props, "country"),
			Region:       utils.GetStringPropOrEmpty(v.Node.Props, "region"),
			Locality:     utils.GetStringPropOrEmpty(v.Node.Props, "locality"),
			AddressLine1: utils.GetStringPropOrEmpty(v.Node.Props, "address"),
			AddressLine2: utils.GetStringPropOrEmpty(v.Node.Props, "address2"),
			ZipCode:      utils.GetStringPropOrEmpty(v.Node.Props, "zip"),
			AddressType:  utils.GetStringPropOrEmpty(v.Node.Props, "addressType"),
			HouseNumber:  utils.GetStringPropOrEmpty(v.Node.Props, "houseNumber"),
			PostalCode:   utils.GetStringPropOrEmpty(v.Node.Props, "postalCode"),
			PlusFour:     utils.GetStringPropOrEmpty(v.Node.Props, "plusFour"),
			Commercial:   utils.GetBoolPropOrFalse(v.Node.Props, "commercial"),
			Predirection: utils.GetStringPropOrEmpty(v.Node.Props, "predirection"),
			District:     utils.GetStringPropOrEmpty(v.Node.Props, "district"),
			Street:       utils.GetStringPropOrEmpty(v.Node.Props, "street"),
			Latitude:     utils.FloatToString(utils.GetFloatPropOrNil(v.Node.Props, "latitude")),
			Longitude:    utils.FloatToString(utils.GetFloatPropOrNil(v.Node.Props, "longitude")),
		})
		if err != nil {
			failedRecords++
			s.log.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}

func (s *syncToEventStoreService) SyncContacts(ctx context.Context, batchSize int) {
	s.log.Info("start sync contacts to eventstore")

	completed, failed, _ := s.upsertContactsIntoEventStore(ctx, batchSize)

	s.log.Infof("completed {%d} and failed {%d} contacts upserting to eventstore", completed, failed)
}

func (s *syncToEventStoreService) upsertContactsIntoEventStore(ctx context.Context, batchSize int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.ContactRepository.GetAllCrossTenantsNotSynced(ctx, batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.ContactClient.UpsertContact(context.Background(), &contactpb.UpsertContactGrpcRequest{
			Id:              utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant:          v.LinkedNodeId,
			FirstName:       utils.GetStringPropOrEmpty(v.Node.Props, "firstName"),
			LastName:        utils.GetStringPropOrEmpty(v.Node.Props, "lastName"),
			Name:            utils.GetStringPropOrEmpty(v.Node.Props, "name"),
			Description:     utils.GetStringPropOrEmpty(v.Node.Props, "description"),
			Timezone:        utils.GetStringPropOrEmpty(v.Node.Props, "timezone"),
			ProfilePhotoUrl: utils.GetStringPropOrEmpty(v.Node.Props, "profilePhotoUrl"),
			Prefix:          utils.GetStringPropOrEmpty(v.Node.Props, "prefix"),
			AppSource:       utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
			Source:          utils.GetStringPropOrEmpty(v.Node.Props, "source"),
			SourceOfTruth:   utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
			CreatedAt:       utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt:       utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
		})
		if err != nil {
			failedRecords++
			s.log.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}

func (s *syncToEventStoreService) SyncOrganizations(ctx context.Context, batchSize int) {
	s.log.Info("start sync organizations to eventstore")

	completed, failed, _ := s.upsertOrganizationsIntoEventStore(ctx, batchSize)

	s.log.Infof("completed {%d} and failed {%d} organizations upserting to eventstore", completed, failed)
}

func (s *syncToEventStoreService) upsertOrganizationsIntoEventStore(ctx context.Context, batchSize int) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.OrganizationRepository.GetAllCrossTenantsNotSynced(ctx, batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.OrganizationClient.UpsertOrganization(context.Background(), &organizationpb.UpsertOrganizationGrpcRequest{
			Id:                utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant:            v.LinkedNodeId,
			Name:              utils.GetStringPropOrEmpty(v.Node.Props, "name"),
			Hide:              utils.GetBoolPropOrFalse(v.Node.Props, "hide"),
			Description:       utils.GetStringPropOrEmpty(v.Node.Props, "description"),
			Website:           utils.GetStringPropOrEmpty(v.Node.Props, "website"),
			Industry:          utils.GetStringPropOrEmpty(v.Node.Props, "industry"),
			IsPublic:          utils.GetBoolPropOrFalse(v.Node.Props, "isPublic"),
			Employees:         utils.GetInt64PropOrZero(v.Node.Props, "employees"),
			Market:            utils.GetStringPropOrEmpty(v.Node.Props, "market"),
			ValueProposition:  utils.GetStringPropOrEmpty(v.Node.Props, "valueProposition"),
			TargetAudience:    utils.GetStringPropOrEmpty(v.Node.Props, "targetAudience"),
			SubIndustry:       utils.GetStringPropOrEmpty(v.Node.Props, "subIndustry"),
			IndustryGroup:     utils.GetStringPropOrEmpty(v.Node.Props, "industryGroup"),
			LastFundingRound:  utils.GetStringPropOrEmpty(v.Node.Props, "lastFundingRound"),
			LastFundingAmount: utils.GetStringPropOrEmpty(v.Node.Props, "lastFundingAmount"),
			Note:              utils.GetStringPropOrEmpty(v.Node.Props, "note"),
			ReferenceId:       utils.GetStringPropOrEmpty(v.Node.Props, "referenceId"),
			Relationship:      utils.GetStringPropOrEmpty(v.Node.Props, "relationship"),
			Stage:             utils.GetStringPropOrEmpty(v.Node.Props, "stage"),
			SourceFields: &commonpb.SourceFields{
				AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
				Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
				SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
			},
			CreatedAt: utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt: utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
			FieldsMask: []organizationpb.OrganizationMaskField{
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NAME,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_HIDE,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_DESCRIPTION,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_WEBSITE,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_INDUSTRY,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_IS_PUBLIC,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_EMPLOYEES,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_MARKET,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_VALUE_PROPOSITION,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_TARGET_AUDIENCE,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_SUB_INDUSTRY,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_INDUSTRY_GROUP,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_LAST_FUNDING_ROUND,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_LAST_FUNDING_AMOUNT,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_NOTE,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_REFERENCE_ID,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_RELATIONSHIP,
				organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_STAGE,
			},
		})
		if err != nil {
			failedRecords++
			s.log.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}

func (s *syncToEventStoreService) SyncOrganizationsLinksWithDomains(ctx context.Context, batchSize int) {
	s.log.Info("start sync organization-domain links to eventstore")
	completed := 0
	failed := 0

	records, _ := s.repositories.OrganizationRepository.GetAllDomainLinksCrossTenantsNotSynced(ctx, batchSize)
	for _, v := range records {
		_, err := s.grpcClients.OrganizationClient.LinkDomainToOrganization(context.Background(), &organizationpb.LinkDomainToOrganizationGrpcRequest{
			OrganizationId: v.Values[0].(string),
			Tenant:         v.Values[1].(string),
			Domain:         v.Values[2].(string),
		})
		if err != nil {
			failed++
			s.log.Errorf("Failed to call method: %v", err)
		} else {
			completed++
		}
	}

	s.log.Infof("completed {%d} and failed {%d} organization-domain syncing to eventstore at %v", completed, failed, utils.Now())
}
