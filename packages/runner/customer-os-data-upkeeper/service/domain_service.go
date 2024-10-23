package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"
)

type DomainService interface {
	CheckDomains()
}

type domainService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonservice.Services
}

func NewDomainService(cfg *config.Config, log logger.Logger, commonServices *commonservice.Services) DomainService {
	return &domainService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *domainService) CheckDomains() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "DomainService.CheckDomains")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 500
	delayFromLastUpdateInDays := 90

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.commonServices.Neo4jRepositories.DomainReadRepository.GetDomainsForPrimaryCheck(ctx, delayFromLastUpdateInDays, limit)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error getting domains for primary check"))
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			_, err = s.commonServices.DomainService.UpdateDomainPrimaryDetails(ctx, record)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error updating domain primary details"))
				s.log.Errorf("Error updating domain primary details: %s", err.Error())
			}
		}
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}
