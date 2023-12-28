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
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"
)

type PhoneNumberService interface {
	CreatePhoneNumber(ctx context.Context, phoneNumber, source, appSource string) (string, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.PhoneNumberEntity, error)
	GetById(ctx context.Context, phoneNumberId string) (*entity.PhoneNumberEntity, error)
}

type phoneNumberService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewPhoneNumberService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) PhoneNumberService {
	return &phoneNumberService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *phoneNumberService) CreatePhoneNumber(ctx context.Context, phoneNumber, source, appSource string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.CreatePhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber), log.String("source", source), log.String("appSource", appSource))

	phoneNumber = strings.ToLower(strings.TrimSpace(phoneNumber))

	var phoneNumberEntity *entity.PhoneNumberEntity
	phoneNumberEntity, _ = s.GetByPhoneNumber(ctx, phoneNumber)
	if phoneNumberEntity == nil {
		// phone number not exist, create new one
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		response, err := s.grpcClients.PhoneNumberClient.UpsertPhoneNumber(ctx, &phonenumberpb.UpsertPhoneNumberGrpcRequest{
			Tenant:      common.GetTenantFromContext(ctx),
			PhoneNumber: phoneNumber,
			SourceFields: &commonpb.SourceFields{
				Source:    source,
				AppSource: utils.StringFirstNonEmpty(appSource, constants.AppSourceCustomerOsWebhooks),
			},
		})
		if err != nil {
			tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertPhoneNumber"))
			s.log.Errorf("Error from events processing %s", err.Error())
			return "", err
		}
		// wait for neo4j to finish processing
		for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
			phoneNumberEntity, findPhoneNumErr := s.GetById(ctx, response.Id)
			if phoneNumberEntity != nil && findPhoneNumErr == nil {
				break
			}
			time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
		}
		span.LogFields(log.String("createdPhoneNumberId", response.Id))
		return response.Id, nil
	} else {
		return phoneNumberEntity.Id, nil
	}
}

func (s *phoneNumberService) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	phoneNumberNode, err := s.repositories.PhoneNumberRepository.GetByPhoneNumber(ctx, common.GetTenantFromContext(ctx), phoneNumber)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToPhoneNumberEntity(*phoneNumberNode), nil
}

func (s *phoneNumberService) GetById(ctx context.Context, phoneNumberId string) (*entity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	phoneNumberNode, err := s.repositories.PhoneNumberRepository.GetById(ctx, phoneNumberId)
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToPhoneNumberEntity(*phoneNumberNode), nil
}

func (s *phoneNumberService) mapDbNodeToPhoneNumberEntity(node dbtype.Node) *entity.PhoneNumberEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.PhoneNumberEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		E164:           utils.GetStringPropOrEmpty(props, "e164"),
		RawPhoneNumber: utils.GetStringPropOrEmpty(props, "rawPhoneNumber"),
		Source:         entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
}
