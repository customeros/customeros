package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/pkg/errors"
)

type EmailService interface {
	ValidateEmails()
}

type emailService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.Services
}

func (s *emailService) ValidateEmails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.ValidateEmails")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := s.cfg.Limits.EmailsValidationLimit
	delayFromLastUpdateInMinutes := 2
	delayFromLastValidationAttemptInMinutes := 24 * 60

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.EmailReadRepository.GetEmailsForValidation(ctx, delayFromLastUpdateInMinutes, delayFromLastValidationAttemptInMinutes, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
				return s.commonServices.GrpcClients.EmailClient.RequestEmailValidation(ctx, &emailpb.RequestEmailValidationGrpcRequest{
					Tenant:    record.Tenant,
					Id:        record.EmailId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error requesting email validation"))
				s.log.Errorf("Error validating email {%s}: %s", record.EmailId, err.Error())
			}

			err = s.commonServices.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, record.Tenant, model.NodeLabelEmail, record.EmailId, string(neo4jentity.EmailPropertyValidationRequestedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, err)
			}
		}
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}

}

func NewEmailService(cfg *config.Config, log logger.Logger, commonServices *commonservice.Services) EmailService {
	return &emailService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}
