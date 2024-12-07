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

func NewMailstackService(cfg *config.Config, log logger.Logger, commonServices *commonservice.Services) MailstackService {
	return &mailstackService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *mailstackService) CheckMailstackDomainReputation() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "MailstackService.CheckMailstackDomainReputation")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	// get all mailstack domains from database
	domains, err := s.commonServices.MailstackService.GetAllMailstackDomains(ctx)
	if err != nil {
		err = errors.Wrap(err, "Unable to get mailstack domains from db")
		tracing.TraceErr(span, err)
		return
	}

	// no record
	if len(domains) == 0 {
		return
	}

	for domain, tenant := range domains {
		_, err := s.commonServices.MailboxService.ReputationScore(ctx, domain, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error getting mailbox reputation score"))
			s.log.Errorf("Error getting mailbox reputation score for %s: %s", domain, err.Error())
		}
	}

}
