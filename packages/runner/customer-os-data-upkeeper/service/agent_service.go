package service

import (
	"context"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
)

type AgentService interface {
	RerunExecutions()
}

type agentService struct {
	cfg            *config.Config
	log            logger.Logger
	commonServices *commonService.CommonServices
}

func NewAgentService(cfg *config.Config, log logger.Logger, commonServices *commonService.CommonServices) AgentService {
	return &agentService{
		cfg:            cfg,
		log:            log,
		commonServices: commonServices,
	}
}

func (s *agentService) RerunExecutions() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.RerunExecutions")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	limit := 200

	agentExecutions, err := s.commonServices.PostgresRepositories.AgentExecutionRepository.GetExecutionsForRetry(ctx, limit)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting agent executions for retry"))
		s.log.Errorf("Error getting agent executions for retry: %s", err.Error())
		return
	}

	if len(agentExecutions) == 0 {
		return
	}

	for _, agentExecution := range agentExecutions {
		func(agentExecution postgres_entity.AgentExecution) {
			recordCtx := common.WithCustomContext(ctx, &common.CustomContext{
				Tenant:    agentExecution.Tenant,
				AppSource: constants.AppSourceDataUpkeeper,
			})
			recordSpan, recordCtx := tracing.StartTracerSpan(recordCtx, "AgentService.RerunExecutions.Record")
			defer recordSpan.Finish()
			tracing.TagTenant(recordSpan, agentExecution.Tenant)
			tracing.TagEntity(recordSpan, agentExecution.ID)
			span.LogFields(log.Int("retryCount", agentExecution.RetryCount))

			err = s.commonServices.AgentRunnerService.RerunExecution(recordCtx, agentExecution.ID)
			if err != nil {
				tracing.TraceErr(recordSpan, errors.Wrap(err, "error retrying agent execution"))
				s.log.Errorf("Error retrying agent execution: %s", err.Error())
			}
		}(agentExecution)
	}
}
