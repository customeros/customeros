package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
)

type WorkflowService interface {
	ExecuteWorkflows()
}

type workflowService struct {
	cfg            *config.Config
	log            logger.Logger
	repositories   *repository.Repositories
	commonServices *commonService.Services
}

func NewWorkflowService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, commonServices *commonService.Services) WorkflowService {
	return &workflowService{
		cfg:            cfg,
		log:            log,
		repositories:   repositories,
		commonServices: commonServices,
	}
}

func (s *workflowService) ExecuteWorkflows() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	liveWorkflows, err := s.repositories.PostgresRepositories.WorkflowRepository.GetAllTenantsLiveWorkflows(ctx)
	if err != nil {
		tracing.TraceErr(nil, err)
		s.log.Errorf("Error getting live workflows: %v", err)
		return
	}

	// execute all live workflows
	for _, workflow := range liveWorkflows {
		err = s.commonServices.WorkflowService.ExecuteWorkflow(ctx, workflow.Tenant, workflow.ID)
		if err != nil {
			tracing.TraceErr(nil, err)
			s.log.Errorf("Error executing workflow {%d}: %v", workflow.ID, err)
		}
	}
}
