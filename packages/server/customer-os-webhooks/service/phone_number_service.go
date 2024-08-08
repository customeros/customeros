package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"
)

type PhoneNumberService interface {
	CreatePhoneNumber(ctx context.Context, phoneNumber, source, appSource string) (string, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4jentity.PhoneNumberEntity, error)
	GetById(ctx context.Context, phoneNumberId string) (*neo4jentity.PhoneNumberEntity, error)
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

	var phoneNumberEntity *neo4jentity.PhoneNumberEntity
	phoneNumberEntity, _ = s.GetByPhoneNumber(ctx, phoneNumber)
	if phoneNumberEntity == nil {
		// phone number not exist, create new one
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		response, err := CallEventsPlatformGRPCWithRetry[*phonenumberpb.PhoneNumberIdGrpcResponse](func() (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
			return s.grpcClients.PhoneNumberClient.UpsertPhoneNumber(ctx, &phonenumberpb.UpsertPhoneNumberGrpcRequest{
				Tenant:      common.GetTenantFromContext(ctx),
				PhoneNumber: phoneNumber,
				SourceFields: &commonpb.SourceFields{
					Source:    source,
					AppSource: utils.StringFirstNonEmpty(appSource, constants.AppSourceCustomerOsWebhooks),
				},
			})
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
			time.Sleep(utils.BackOffExponentialDelay(i))
		}
		span.LogFields(log.String("createdPhoneNumberId", response.Id))
		return response.Id, nil
	} else {
		return phoneNumberEntity.Id, nil
	}
}

func (s *phoneNumberService) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4jentity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	phoneNumberNode, err := s.repositories.PhoneNumberRepository.GetByPhoneNumber(ctx, common.GetTenantFromContext(ctx), phoneNumber)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToPhoneNumberEntity(phoneNumberNode), nil
}

func (s *phoneNumberService) GetById(ctx context.Context, phoneNumberId string) (*neo4jentity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	phoneNumberNode, err := s.repositories.PhoneNumberRepository.GetById(ctx, phoneNumberId)
	if err != nil {
		return nil, err
	}

	return neo4jmapper.MapDbNodeToPhoneNumberEntity(phoneNumberNode), nil
}
