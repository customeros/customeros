package service

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
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

	spans, ctx := telemetry.StartCronSpan(ctx, "AgentService.RerunExecutions")
	defer spans.Finish()

	limit := 200

	agentExecutions, err := s.commonServices.PostgresRepositories.AgentExecutionRepository.GetExecutionsForRetry(ctx, limit)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error getting agent executions for retry"))
		s.log.Errorf("Error getting agent executions for retry: %s", err.Error())
		return
	}

	spans.LogKV("agentExecutions", len(agentExecutions))
	if len(agentExecutions) == 0 {
		return
	}

	for _, agentExecution := range agentExecutions {
		s.rerunAgentExecution(ctx, agentExecution)
	}
}

func (s *agentService) rerunAgentExecution(ctx context.Context, agentExecution postgres_entity.AgentExecution) {
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    agentExecution.Tenant,
		AppSource: constants.AppSourceDataUpkeeper,
	})
	spans, ctx := telemetry.StartCronSpan(ctx, "AgentService.RerunExecutions.Record")
	defer spans.Finish()
	spans.TagEntity(agentExecution.ID)

	spans.LogKV("retryCount", agentExecution.RetryCount)

	err := s.commonServices.AgentRunnerService.RerunExecution(ctx, agentExecution.ID)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error retrying agent execution"))
		s.log.Errorf("Error retrying agent execution: %s", err.Error())
	}
}
