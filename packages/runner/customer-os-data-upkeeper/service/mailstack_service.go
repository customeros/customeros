package service

import (
	"context"

	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
)

type MailstackService interface {
	CheckMailstackDomainReputation()
}

type mailstackService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.Services
}

func NewMailstackService(cfg *config.Config, log logger.Logger, commonServices *commonservice.Services) DomainService {
	return &domainService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *mailstackService) NewMailstackService() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "MailstackService.CheckMailstackDomainReputation")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		// get mailstack domains from database

		// no record
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			_, err := s.commonServices.MailboxService.ReputationScore(ctx, domain, tenant)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error getting mailbox reputation score"))
				s.log.Errorf("Error getting mailbox reputation score for %s: %s", domain, err.Error())
			}
		}

		// force exit after single iteration
		return
	}
}
