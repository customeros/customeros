package service

import (
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type IssueSyncService interface {
	SyncIssues(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int)
}

type issueSyncService struct {
	repositories *repository.Repositories
}

func NewIssueSyncService(repositories *repository.Repositories) IssueSyncService {
	return &issueSyncService{
		repositories: repositories,
	}
}

func (s *issueSyncService) SyncIssues(ctx context.Context, dataService common.SourceDataService, syncDate time.Time, tenant, runId string) (int, int) {
	completed, failed := 0, 0
	for {
		issues := dataService.GetIssuesForSync(batchSize, runId)
		if len(issues) == 0 {
			logrus.Debugf("no issues found for sync from %s for tenant %s", dataService.SourceId(), tenant)
			break
		}
		logrus.Infof("syncing %d issues from %s for tenant %s", len(issues), dataService.SourceId(), tenant)

		for _, v := range issues {
			var failedSync = false

			issueId, err := s.repositories.IssueRepository.GetMatchedIssueId(ctx, tenant, v)
			if err != nil {
				failedSync = true
				logrus.Errorf("failed finding existing matched issue with external reference id %v for tenant %v :%v", v.ExternalId, tenant, err)
			}

			// Create new issue id if not found
			if len(issueId) == 0 {
				issueUuid, _ := uuid.NewRandom()
				issueId = issueUuid.String()
			}
			v.Id = issueId

			if !failedSync {
				err = s.repositories.IssueRepository.MergeIssue(ctx, tenant, syncDate, v)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed merging issue with external reference id %v for tenant %v :%v", v.ExternalId, tenant, err)
				}
			}

			if v.HasReporterOrganization() && !failedSync {
				err = s.repositories.IssueRepository.LinkIssueWithReporterOrganizationByExternalId(ctx, tenant, issueId, v.ReporterOrganizationExternalId, v.ExternalSystem)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link issue %v with reporter organization for tenant %v :%v", issueId, tenant, err)
				}
			}

			if v.HasCollaboratorUsers() && !failedSync {
				for _, userExternalId := range v.CollaboratorUserExternalIds {
					err = s.repositories.IssueRepository.LinkIssueWithCollaboratorUserByExternalId(ctx, tenant, issueId, userExternalId, v.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link issue %v with collaborator user for tenant %v :%v", issueId, tenant, err)
					}
				}
			}

			if v.HasFollowerUsers() && !failedSync {
				for _, userExternalId := range v.FollowerUserExternalIds {
					err = s.repositories.IssueRepository.LinkIssueWithFollowerUserByExternalId(ctx, tenant, issueId, userExternalId, v.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed link issue %v with follower user for tenant %v :%v", issueId, tenant, err)
					}
				}
			}

			if v.HasAssignee() && !failedSync {
				err = s.repositories.IssueRepository.LinkIssueWithAssigneeUserByExternalId(ctx, tenant, issueId, v.AssigneeUserExternalId, v.ExternalSystem)
				if err != nil {
					failedSync = true
					logrus.Errorf("failed link issue %v with assignee user for tenant %v :%v", issueId, tenant, err)
				}
			}

			v.Tags = append(v.Tags, v.Subject+" - "+v.ExternalId)
			if v.HasTags() && !failedSync {
				for _, tag := range v.Tags {
					err = s.repositories.IssueRepository.MergeTagForIssue(ctx, tenant, issueId, tag, v.ExternalSystem)
					if err != nil {
						failedSync = true
						logrus.Errorf("failed to merge tag %v for issue %v, tenant %v :%v", tag, issueId, tenant, err)
					}
				}
			}

			logrus.Debugf("successfully merged issue with id %v for tenant %v from %v", issueId, tenant, dataService.SourceId())
			if err := dataService.MarkIssueProcessed(v.ExternalSyncId, runId, failedSync == false); err != nil {
				failed++
				continue
			}
			if failedSync == true {
				failed++
			} else {
				completed++
			}
		}
		if len(issues) < batchSize {
			break
		}
	}
	return completed, failed
}
