package service

import (
	"context"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"
)

type DomainService interface {
	CheckDomains()
}

type domainService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *service.CommonServices
}

func NewDomainService(cfg *config.Config, log logger.Logger, commonServices *service.CommonServices) DomainService {
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

	limit := 50
	delayFromLastUpdateInDays := 30

	records, err := s.commonServices.Neo4jRepositories.DomainReadRepository.GetDomainsForPrimaryCheck(ctx, delayFromLastUpdateInDays, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting domains for primary check"))
		return
	}

	// no record
	if len(records) == 0 {
		return
	}

	for _, domain := range records {
		err = s.commonServices.DomainService.UpdateDomainPrimaryDetails(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error updating domain primary details"))
			s.log.Errorf("Error updating domain primary details: %s", err.Error())
		}
	}
}
