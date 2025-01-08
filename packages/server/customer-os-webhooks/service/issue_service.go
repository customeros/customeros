package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	pkgerrors "github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

type IssueService interface {
	SyncIssues(ctx context.Context, contacts []model.IssueData) (SyncResult, error)
	GetIdForReferencedIssue(ctx context.Context, tenant, externalSystemId string, issue model.ReferencedIssue) (string, error)
}

type issueService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewIssueService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) IssueService {
	return &issueService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.IssueSyncConcurrency,
	}
}

func (s *issueService) SyncIssues(ctx context.Context, issues []model.IssueData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.SyncIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("num of issues", len(issues)))

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate issues input before syncing
	for _, issue := range issues {
		if issue.ExternalSystem == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if !neo4jentity.IsValidDataSource(strings.ToLower(issue.ExternalSystem)) {
			tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", issue.ExternalSystem))
			return SyncResult{}, errors.ErrExternalSystemNotAccepted
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, s.maxWorkers)

	syncMutex := &sync.Mutex{}
	statusesMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	// Sync all issues concurrently
	for _, issueData := range issues {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
		default:
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(issueData model.IssueData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncIssue(ctx, syncMutex, issueData, syncDate)
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
		}(issueData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), issues[0].ExternalSystem,
		issues[0].AppSource, "issue", syncDate, statuses)

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *issueService) syncIssue(ctx context.Context, syncMutex *sync.Mutex, issueInput model.IssueData, syncDate time.Time) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.syncIssue")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagExternalSystem, issueInput.ExternalSystem)
	span.SetTag(tracing.SpanTagExternalId, issueInput.ExternalId)
	span.LogFields(log.Object("syncDate", syncDate))
	tracing.LogObjectAsJson(span, "issueInput", issueInput)

	tenant := common.GetTenantFromContext(ctx)
	var failedSync = false
	var reason = ""

	issueInput.Normalize()

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, issueInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", issueInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	// Check if contact sync should be skipped
	if issueInput.Skip {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(issueInput.SkipReason)
	} else if issueInput.ExternalId == "" {
		reason = fmt.Sprintf("external id is empty for issue, tenant %s", tenant)
		s.log.Warnf("Skip issue sync: %v", reason)
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(reason)
	}

	reporterId, reporterLabel, err := s.services.FinderService.FindReferencedEntityId(ctx, issueInput.ExternalSystem, &issueInput.Reporter)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding reporter for issue %s for tenant %s :%s", issueInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}
	submitterId, submitterLabel, err := s.services.FinderService.FindReferencedEntityId(ctx, issueInput.ExternalSystem, &issueInput.Submitter)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding submitter for issue %s for tenant %s :%s", issueInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	if issueInput.OrganizationRequired && reporterLabel != commonmodel.NodeLabelOrganization {
		reason = fmt.Sprintf("organization(s) not found for issue %s for tenant %s", issueInput.ExternalId, tenant)
		s.log.Warnf("Skip issue sync: %v", reason)
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(reason)
	}

	// Lock issue creation
	syncMutex.Lock()
	defer syncMutex.Unlock()
	// Check if issue already exists
	issueId, err := s.repositories.Neo4jRepositories.IssueReadRepository.GetMatchedIssueId(ctx, tenant, issueInput.ExternalSystem, issueInput.ExternalId)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched issue with external reference %s for tenant %s :%s", issueInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}
	if !failedSync {
		span.LogFields(log.Bool("found matching issue", issueId != ""))

		issueFields := data_fields.IssueFields{
			Source:    utils.StringPtr(issueInput.ExternalSystem),
			CreatedAt: issueInput.CreatedAt,
			ExternalSystem: &neo4jmodel.ExternalSystem{
				ExternalSystemId: issueInput.ExternalSystem,
				ExternalId:       issueInput.ExternalId,
				ExternalIdSecond: issueInput.ExternalIdSecond,
				ExternalSource:   issueInput.ExternalSourceEntity,
				ExternalUrl:      issueInput.ExternalUrl,
				SyncDate:         &syncDate,
			},
		}
		if issueInput.Subject != "" {
			issueFields.Subject = &issueInput.Subject
		}
		if issueInput.Status != "" {
			issueFields.Status = &issueInput.Status
		}
		if issueInput.Priority != "" {
			issueFields.Priority = &issueInput.Priority
		}
		if issueInput.Description != "" {
			issueFields.Description = &issueInput.Description
		}
		if issueInput.GroupId != "" {
			issueFields.GroupId = &issueInput.GroupId
		}
		if reporterId != "" && reporterLabel == commonmodel.NodeLabelOrganization {
			issueFields.ReportedByOrganizationId = &reporterId
		}
		if submitterId != "" {
			switch submitterLabel {
			case commonmodel.NodeLabelOrganization:
				issueFields.SubmittedByOrganizationId = &submitterId
			case commonmodel.NodeLabelUser:
				issueFields.SubmittedByUserId = &submitterId
			}
		}

		issueId, err = s.services.CommonServices.IssueService.Save(ctx, nil, &issueId, issueFields)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			s.log.Error(reason)
			reason = fmt.Sprintf("error saving issue with external reference %s for tenant %s :%s", issueInput.ExternalId, common.GetTenantFromContext(ctx), err.Error())
		}
		issueInput.Id = issueId
		tracing.TagEntity(span, issueId)
	}

	processedFollowerUserIds := make([]string, 0)
	// add user followers
	if !failedSync && issueInput.HasFollowers() {
		for _, follower := range issueInput.Followers {
			// find follower
			followerId, followerLabel, err := s.services.FinderService.FindReferencedEntityId(ctx, issueInput.ExternalSystem, &follower)
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed finding follower for issue %s for tenant %s :%s", issueInput.ExternalId, tenant, err.Error())
				s.log.Error(reason)
			}
			if followerId != "" && followerLabel == commonmodel.NodeLabelUser && !utils.Contains(processedFollowerUserIds, followerId) {
				err = s.services.CommonServices.IssueService.AddUserFollower(ctx, nil, issueId, followerId)
				processedFollowerUserIds = append(processedFollowerUserIds, followerId)
				if err != nil {
					tracing.TraceErr(span, pkgerrors.Wrap(err, "AddUserFollower"))
					reason = fmt.Sprintf("failed to add follower %s to issue %s for tenant %s :%s", followerId, issueId, tenant, err.Error())
					s.log.Error(reason)
				}
			}
		}
	}

	// add user collaborators as followers
	if !failedSync && issueInput.HasCollaborators() {
		for _, collaborator := range issueInput.Collaborators {
			// find collaborator
			collaboratorId, collaboratorLabel, err := s.services.FinderService.FindReferencedEntityId(ctx, issueInput.ExternalSystem, &collaborator)
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed finding collaborator for issue %s for tenant %s :%s", issueInput.ExternalId, tenant, err.Error())
				s.log.Error(reason)
			}
			if collaboratorId != "" && collaboratorLabel == commonmodel.NodeLabelUser && !utils.Contains(processedFollowerUserIds, collaboratorId) {
				err = s.services.CommonServices.IssueService.AddUserFollower(ctx, nil, issueId, collaboratorId)
				processedFollowerUserIds = append(processedFollowerUserIds, collaboratorId)
				if err != nil {
					tracing.TraceErr(span, pkgerrors.Wrap(err, "AddUserFollower"))
					reason = fmt.Sprintf("failed to add follower %s to issue %s for tenant %s :%s", collaboratorId, issueId, tenant, err.Error())
					s.log.Error(reason)
				}
			}
		}
	}

	// add assignee
	if !failedSync {
		// find assignee
		assigneeId, err := s.services.UserService.GetIdForReferencedUser(ctx, tenant, issueInput.ExternalSystem, issueInput.Assignee)
		if err != nil {
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed finding assignee for issue %s for tenant %s :%s", issueInput.ExternalId, tenant, err.Error())
			s.log.Error(reason)
		}
		if assigneeId != "" {
			err = s.services.CommonServices.IssueService.AddUserAssignee(ctx, nil, issueId, assigneeId)
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed to add assignee %s to issue %s for tenant %s :%s", assigneeId, issueId, tenant, err.Error())
				s.log.Error(reason)
			}
		}
	}

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}
	span.LogFields(log.String("output", "success"))
	return NewSuccessfulSyncStatus()
}

func (s *issueService) GetIdForReferencedIssue(ctx context.Context, tenant, externalSystemId string, issue model.ReferencedIssue) (string, error) {
	if !issue.Available() {
		return "", nil
	}

	if issue.ReferencedByExternalId() {
		return s.repositories.Neo4jRepositories.IssueReadRepository.GetIssueIdByExternalId(ctx, tenant, issue.ExternalId, externalSystemId)
	}
	return "", nil
}
