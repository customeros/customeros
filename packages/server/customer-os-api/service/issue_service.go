package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/exp/slices"
)

type IssueService interface {
	GetIssueSummaryByStatusForOrganization(ctx context.Context, organizationId string) (map[string]int64, error)
	GetById(ctx context.Context, issueId string) (*entity.IssueEntity, error)
	GetIssuesForInteractionEvents(ctx context.Context, ids []string) (*entity.IssueEntities, error)
	GetSubmitterParticipantsForIssues(ctx context.Context, ids []string) (*entity.IssueParticipants, error)
	GetReporterParticipantsForIssues(ctx context.Context, ids []string) (*entity.IssueParticipants, error)
	GetAssigneeParticipantsForIssues(ctx context.Context, ids []string) (*entity.IssueParticipants, error)
	GetFollowerParticipantsForIssues(ctx context.Context, ids []string) (*entity.IssueParticipants, error)

	mapDbNodeToIssue(node dbtype.Node) *entity.IssueEntity
}

type issueService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewIssueService(log logger.Logger, repositories *repository.Repositories, services *Services) IssueService {
	return &issueService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *issueService) GetIssueSummaryByStatusForOrganization(ctx context.Context, organizationId string) (map[string]int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.GetIssueSummaryByStatusForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	return s.repositories.IssueRepository.GetIssueCountByStatusForOrganization(ctx, common.GetTenantFromContext(ctx), organizationId)
}

func (s *issueService) GetById(ctx context.Context, issueId string) (*entity.IssueEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("issueId", issueId))

	if issueDbNode, err := s.repositories.IssueRepository.GetById(ctx, common.GetTenantFromContext(ctx), issueId); err != nil {
		return nil, err
	} else {
		return s.mapDbNodeToIssue(*issueDbNode), nil
	}
}

func (s *issueService) GetIssuesForInteractionEvents(ctx context.Context, ids []string) (*entity.IssueEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.GetIssuesForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	issues, err := s.repositories.IssueRepository.GetAllForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	issueEntities := make(entity.IssueEntities, 0, len(issues))
	for _, v := range issues {
		issueEntity := s.mapDbNodeToIssue(*v.Node)
		issueEntity.DataloaderKey = v.LinkedNodeId
		issueEntities = append(issueEntities, *issueEntity)
	}
	return &issueEntities, nil
}

func (s *issueService) GetSubmitterParticipantsForIssues(ctx context.Context, issueIds []string) (*entity.IssueParticipants, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.GetSubmitterParticipantsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	records, err := s.repositories.IssueRepository.GetSubmitterParticipantsForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}

	issueParticipants := s.convertDbNodesToIssueParticipants(records)

	span.LogFields(log.Int("result count", len(issueParticipants)))

	return &issueParticipants, nil
}

func (s *issueService) GetReporterParticipantsForIssues(ctx context.Context, issueIds []string) (*entity.IssueParticipants, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.GetReporterParticipantsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	records, err := s.repositories.IssueRepository.GetReporterParticipantsForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}

	issueParticipants := s.convertDbNodesToIssueParticipants(records)

	span.LogFields(log.Int("result count", len(issueParticipants)))

	return &issueParticipants, nil
}

func (s *issueService) GetAssigneeParticipantsForIssues(ctx context.Context, issueIds []string) (*entity.IssueParticipants, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.GetAssigneeParticipantsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	records, err := s.repositories.IssueRepository.GetAssigneeParticipantsForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}

	issueParticipants := s.convertDbNodesToIssueParticipants(records)

	span.LogFields(log.Int("result count", len(issueParticipants)))

	return &issueParticipants, nil
}

func (s *issueService) GetFollowerParticipantsForIssues(ctx context.Context, issueIds []string) (*entity.IssueParticipants, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.GetFollowerParticipantsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	records, err := s.repositories.IssueRepository.GetFollowerParticipantsForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		return nil, err
	}

	issueParticipants := s.convertDbNodesToIssueParticipants(records)

	span.LogFields(log.Int("result count", len(issueParticipants)))

	return &issueParticipants, nil
}

func (s *issueService) mapDbNodeToIssue(node dbtype.Node) *entity.IssueEntity {
	props := utils.GetPropsFromNode(node)
	issue := entity.IssueEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrNow(props, "updatedAt"),
		Subject:       utils.GetStringPropOrEmpty(props, "subject"),
		Status:        utils.GetStringPropOrEmpty(props, "status"),
		Priority:      utils.GetStringPropOrEmpty(props, "priority"),
		Description:   utils.GetStringPropOrEmpty(props, "description"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &issue
}

func (s *issueService) convertDbNodesToIssueParticipants(records []*utils.DbNodeAndId) entity.IssueParticipants {
	issueParticipants := entity.IssueParticipants{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, entity.NodeLabel_User) {
			participant := s.services.UserService.mapDbNodeToUserEntity(*v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			issueParticipants = append(issueParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Contact) {
			participant := s.services.ContactService.mapDbNodeToContactEntity(*v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			issueParticipants = append(issueParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Organization) {
			participant := s.services.OrganizationService.mapDbNodeToOrganizationEntity(*v.Node)
			participant.DataloaderKey = v.LinkedNodeId
			issueParticipants = append(issueParticipants, participant)
		}
	}
	return issueParticipants
}
