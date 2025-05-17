package api_issue

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"golang.org/x/exp/slices"
)

type issueService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewIssueService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.IssueService {
	return &issueService{
		log:          log,
		repositories: repositories,
	}
}

func (s *issueService) GetIssueSummaryByStatusForOrganization(ctx context.Context, organizationId string) (map[string]int64, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueService.GetIssueSummaryByStatusForOrganization")
	defer spans.Finish()
	spans.LogKV("organizationId", organizationId)

	return s.repositories.Neo4jRepositories.IssueReadRepository.GetIssueCountByStatusForOrganization(ctx, common.GetTenantFromContext(ctx), organizationId)
}

func (s *issueService) GetById(ctx context.Context, issueId string) (*neo4jentity.IssueEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueService.GetById")
	defer spans.Finish()
	spans.LogKV("issueId", issueId)

	if issueDbNode, err := s.repositories.Neo4jRepositories.IssueReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), issueId); err != nil {
		return nil, err
	} else {
		return neo4jmapper.MapDbNodeToIssueEntity(issueDbNode), nil
	}
}

func (s *issueService) GetIssuesForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.IssueEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueService.GetIssuesForInteractionEvents")
	defer spans.Finish()
	spans.LogObjectAsJson("ids", ids)

	issues, err := s.repositories.IssueRepository.GetAllForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	issueEntities := make(neo4jentity.IssueEntities, 0, len(issues))
	for _, v := range issues {
		issueEntity := neo4jmapper.MapDbNodeToIssueEntity(v.Node)
		issueEntity.DataloaderKey = v.LinkedNodeId
		issueEntities = append(issueEntities, *issueEntity)
	}
	return &issueEntities, nil
}

func (s *issueService) GetSubmitterParticipantsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.IssueParticipants, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueService.GetSubmitterParticipantsForIssues")
	defer spans.Finish()
	spans.LogObjectAsJson("issueIds", issueIds)

	records, err := s.repositories.IssueRepository.GetSubmitterParticipantsForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}

	issueParticipants := s.convertDbNodesToIssueParticipants(records)

	spans.LogKV("result count", len(issueParticipants))

	return &issueParticipants, nil
}

func (s *issueService) GetReporterParticipantsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.IssueParticipants, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueService.GetReporterParticipantsForIssues")
	defer spans.Finish()
	spans.LogObjectAsJson("issueIds", issueIds)

	records, err := s.repositories.IssueRepository.GetReporterParticipantsForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}

	issueParticipants := s.convertDbNodesToIssueParticipants(records)

	spans.LogKV("result count", len(issueParticipants))

	return &issueParticipants, nil
}

func (s *issueService) GetAssigneeParticipantsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.IssueParticipants, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueService.GetAssigneeParticipantsForIssues")
	defer spans.Finish()
	spans.LogObjectAsJson("issueIds", issueIds)

	records, err := s.repositories.IssueRepository.GetAssigneeParticipantsForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}

	issueParticipants := s.convertDbNodesToIssueParticipants(records)

	spans.LogKV("result count", len(issueParticipants))

	return &issueParticipants, nil
}

func (s *issueService) GetFollowerParticipantsForIssues(ctx context.Context, issueIds []string) (*neo4jentity.IssueParticipants, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IssueService.GetFollowerParticipantsForIssues")
	defer spans.Finish()
	spans.LogObjectAsJson("issueIds", issueIds)

	records, err := s.repositories.IssueRepository.GetFollowerParticipantsForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}

	issueParticipants := s.convertDbNodesToIssueParticipants(records)

	spans.LogKV("result count", len(issueParticipants))

	return &issueParticipants, nil
}

func (s *issueService) convertDbNodesToIssueParticipants(records []*utils.DbNodeAndId) neo4jentity.IssueParticipants {
	issueParticipants := neo4jentity.IssueParticipants{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, model.NodeLabelUser) {
			participant := neo4jmapper.MapDbNodeToUserEntity(v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			issueParticipants = append(issueParticipants, participant)
		} else if slices.Contains(v.Node.Labels, model.NodeLabelContact) {
			participant := neo4jmapper.MapDbNodeToContactEntity(v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			issueParticipants = append(issueParticipants, participant)
		} else if slices.Contains(v.Node.Labels, model.NodeLabelOrganization) {
			participant := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			issueParticipants = append(issueParticipants, participant)
		}
	}
	return issueParticipants
}
