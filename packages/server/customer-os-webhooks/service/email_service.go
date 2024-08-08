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
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"
)

type EmailService interface {
	CreateEmail(ctx context.Context, email, source, appSource string) (string, error)
	GetByEmailAddress(ctx context.Context, email string) (*neo4jentity.EmailEntity, error)
	GetById(ctx context.Context, emailId string) (*neo4jentity.EmailEntity, error)
}

type emailService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewEmailService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) EmailService {
	return &emailService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *emailService) CreateEmail(ctx context.Context, email, source, appSource string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.CreateEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email), log.String("source", source), log.String("appSource", appSource))

	email = strings.ToLower(strings.TrimSpace(email))

	var emailEntity *neo4jentity.EmailEntity
	emailEntity, _ = s.GetByEmailAddress(ctx, strings.TrimSpace(email))
	if emailEntity == nil {
		// email address not exist, create new one
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		response, err := CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
			return s.grpcClients.EmailClient.UpsertEmail(ctx, &emailpb.UpsertEmailGrpcRequest{
				Tenant:   common.GetTenantFromContext(ctx),
				RawEmail: email,
				SourceFields: &commonpb.SourceFields{
					Source:    source,
					AppSource: utils.StringFirstNonEmpty(appSource, constants.AppSourceCustomerOsWebhooks),
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertEmail"))
			s.log.Errorf("Error from events processing %s", err.Error())
			return "", err
		}
		// wait for neo4j to finish processing
		for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
			emailEntity, findEmailErr := s.GetById(ctx, response.Id)
			if emailEntity != nil && findEmailErr == nil {
				break
			}
			time.Sleep(utils.BackOffExponentialDelay(i))
		}
		span.LogFields(log.String("createdEmailId", response.Id))
		return response.Id, nil
	} else {
		return emailEntity.Id, nil
	}
}

func (s *emailService) GetByEmailAddress(ctx context.Context, email string) (*neo4jentity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetByEmailAddress")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	emailNode, err := s.repositories.EmailRepository.GetByEmail(ctx, common.GetTenantFromContext(ctx), email)
	if err != nil {
		return nil, err
	}

	return neo4jmapper.MapDbNodeToEmailEntity(emailNode), nil
}

func (s *emailService) GetById(ctx context.Context, emailId string) (*neo4jentity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", emailId))

	emailNode, err := s.repositories.EmailRepository.GetById(ctx, emailId)
	if err != nil {
		return nil, err
	}

	return neo4jmapper.MapDbNodeToEmailEntity(emailNode), nil
}
