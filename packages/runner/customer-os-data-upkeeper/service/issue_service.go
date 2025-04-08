package service

import (
	"context"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type IssueService interface {
	LinkUnthreadIssues()
}

type issueService struct {
	cfg          *config.Config
	log          logger.Logger
	repositories *repository.Repositories
}

func NewIssueService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories) IssueService {
	return &issueService{
		cfg:          cfg,
		log:          log,
		repositories: repositories,
	}
}

func (s *issueService) LinkUnthreadIssues() {
	ctx, cancel := utils.GetLongLivedContext(context.Background())
	defer cancel() // Cancel context on exit

	spans, ctx := telemetry.StartCronSpan(ctx, "IssueService.LinkUnthreadIssues")
	defer spans.Finish()

	err := s.repositories.Neo4jRepositories.IssueWriteRepository.LinkUnthreadIssuesToOrganizationByGroupId(ctx)
	if err != nil {
		spans.TraceError(err)
		s.log.Errorf("Error linking unthread issues to organization by group id: %s", err.Error())
		return
	}
}
