package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type issueSyncService struct {
	repositories *repository.Repositories
	services     *Services
	log          logger.Logger
}

func NewDefaultIssueSyncService(repositories *repository.Repositories, services *Services, log logger.Logger) SyncService {
	return &issueSyncService{
		repositories: repositories,
		services:     services,
		log:          log,
	}
}

func (s *issueSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	for {
		issues := dataService.GetDataForSync(ctx, common.ISSUES, batchSize, runId)
		if len(issues) == 0 {
			s.log.Debugf("no issues found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		s.log.Infof("syncing %d issues from %s for tenant %s", len(issues), dataService.SourceId(), tenant)

		for _, v := range issues {
			s.syncIssue(ctx, v.(entity.IssueData), dataService, syncDate, tenant, runId, &completed, &failed, &skipped)
		}
		if len(issues) < batchSize {
			break
		}
	}
	return completed, failed, skipped
}

func (s *issueSyncService) syncIssue(ctx context.Context, issueInput entity.IssueData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueSyncService.syncIssue")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason = ""
	issueInput.Normalize()

	if issueInput.Skip {
		if err := dataService.MarkProcessed(ctx, issueInput.SyncId, runId, true, true, issueInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	dbNode, err := s.repositories.IssueRepository.GetMatchedIssue(ctx, tenant, issueInput.ExternalSystem, issueInput.ExternalId)
	var issueId string
	if dbNode != nil {
		issueId = utils.GetStringPropOrEmpty(dbNode.Props, "id")
	}
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched issue with external reference id %v for tenant %v :%v", issueInput.ExternalId, tenant, err)
		s.log.Errorf(reason)
	}

	// Create new issue id if not found
	if len(issueId) == 0 {
		issueUuid, _ := uuid.NewRandom()
		issueId = issueUuid.String()
	}
	issueInput.Id = issueId
	span.LogFields(log.String("issueId", issueId))

	if !failedSync {
		err = s.repositories.IssueRepository.MergeIssue(ctx, tenant, syncDate, issueInput)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed merging issue with external reference id %v for tenant %v :%v", issueInput.ExternalId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	if issueInput.HasReporterOrganization() && !failedSync {
		err = s.repositories.IssueRepository.LinkIssueWithReporterOrganizationByExternalId(ctx, tenant, issueId, issueInput.ReporterOrganizationExternalId, issueInput.ExternalSystem)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed link issue %v with reporter organization for tenant %v :%v", issueId, tenant, err)
			s.log.Errorf(reason)
		}
		s.services.OrganizationService.UpdateLastTouchpointByOrganizationExternalId(ctx, tenant, issueInput.ReporterOrganizationExternalId, issueInput.ExternalSystem)
	}

	if issueInput.HasCollaboratorUsers() && !failedSync {
		for _, userExternalId := range issueInput.CollaboratorUserExternalIds {
			err = s.repositories.IssueRepository.LinkIssueWithCollaboratorUserByExternalId(ctx, tenant, issueId, userExternalId, issueInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link issue %v with collaborator user for tenant %v :%v", issueId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if issueInput.HasFollowerUsers() && !failedSync {
		for _, userExternalId := range issueInput.FollowerUserExternalIds {
			err = s.repositories.IssueRepository.LinkIssueWithFollowerUserByExternalId(ctx, tenant, issueId, userExternalId, issueInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed link issue %v with follower user for tenant %v :%v", issueId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	if issueInput.HasAssignee() && !failedSync {
		err = s.repositories.IssueRepository.LinkIssueWithAssigneeUserByExternalId(ctx, tenant, issueId, issueInput.AssigneeUserExternalId, issueInput.ExternalSystem)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed link issue %v with assignee user for tenant %v :%v", issueId, tenant, err)
			s.log.Errorf(reason)
		}
	}

	issueInput.Tags = append(issueInput.Tags, issueInput.Subject+" - "+issueInput.ExternalId)
	if issueInput.HasTags() && !failedSync {
		for _, tag := range issueInput.Tags {
			err = s.repositories.IssueRepository.MergeTagForIssue(ctx, tenant, issueId, tag, issueInput.ExternalSystem)
			if err != nil {
				failedSync = true
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed to merge tag %v for issue %v, tenant %v :%v", tag, issueId, tenant, err)
				s.log.Errorf(reason)
				break
			}
		}
	}

	s.log.Debugf("successfully merged issue with id %v for tenant %v from %v", issueId, tenant, dataService.SourceId())
	if err := dataService.MarkProcessed(ctx, issueInput.SyncId, runId, failedSync == false, false, reason); err != nil {
		*failed++
		span.LogFields(log.Bool("failedSync", true))
		return
	}
	if failedSync == true {
		*failed++
	} else {
		*completed++
	}
	span.LogFields(log.Bool("failedSync", failedSync))
}
