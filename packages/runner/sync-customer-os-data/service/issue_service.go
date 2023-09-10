package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
)

type IssueService interface {
	GetIdForReferencedIssue(ctx context.Context, tenant, externalSystem string, issue entity.ReferencedIssue) (string, error)
}

type issueService struct {
	repositories *repository.Repositories
}

func NewIssueService(repositories *repository.Repositories) IssueService {
	return &issueService{
		repositories: repositories,
	}
}

func (s *issueService) GetIdForReferencedIssue(ctx context.Context, tenant, externalSystemId string, issue entity.ReferencedIssue) (string, error) {
	if !issue.Available() {
		return "", nil
	}

	if issue.ReferencedByExternalId() {
		return s.repositories.IssueRepository.GetIssueIdByExternalId(ctx, tenant, issue.ExternalId, externalSystemId)
	}
	return "", nil
}
