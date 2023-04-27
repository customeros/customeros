package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type IssueService interface {
	GetIssueSummaryByStatusForOrganization(ctx context.Context, organizationId string) (map[string]int64, error)
	GetById(ctx context.Context, issueId string) (*entity.IssueEntity, error)

	mapDbNodeToIssue(node dbtype.Node) *entity.IssueEntity
}

type issueService struct {
	repositories *repository.Repositories
}

func NewIssueService(repositories *repository.Repositories) IssueService {
	return &issueService{
		repositories: repositories,
	}
}

func (s *issueService) GetIssueSummaryByStatusForOrganization(ctx context.Context, organizationId string) (map[string]int64, error) {
	return s.repositories.IssueRepository.GetIssueCountByStatusForOrganization(ctx, common.GetTenantFromContext(ctx), organizationId)
}

func (s *issueService) GetById(ctx context.Context, issueId string) (*entity.IssueEntity, error) {
	if issueDbNode, err := s.repositories.IssueRepository.GetById(ctx, common.GetTenantFromContext(ctx), issueId); err != nil {
		return nil, err
	} else {
		return s.mapDbNodeToIssue(*issueDbNode), nil
	}
}

func (s *issueService) mapDbNodeToIssue(node dbtype.Node) *entity.IssueEntity {
	props := utils.GetPropsFromNode(node)
	issue := entity.IssueEntity{
		Id:          utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:   utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt:   utils.GetTimePropOrNow(props, "updatedAt"),
		Subject:     utils.GetStringPropOrEmpty(props, "subject"),
		Status:      utils.GetStringPropOrEmpty(props, "status"),
		Priority:    utils.GetStringPropOrEmpty(props, "priority"),
		Description: utils.GetStringPropOrEmpty(props, "description"),
	}
	return &issue
}
