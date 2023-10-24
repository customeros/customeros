package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	commongrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	issuegrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/issue"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"sync"
	"time"
)

const maxWorkersIssueSync = 5

type IssueService interface {
	SyncIssues(ctx context.Context, contacts []model.IssueData) error
}

type issueService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewIssueService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) IssueService {
	return &issueService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *issueService) SyncIssues(ctx context.Context, issues []model.IssueData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.SyncIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("num of issues", len(issues)))

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return errors.ErrTenantNotValid
	}

	// pre-validate issues input before syncing
	for _, issue := range issues {
		if issue.ExternalSystem == "" {
			return errors.ErrMissingExternalSystem
		}
		if !entity.IsValidDataSource(strings.ToLower(issue.ExternalSystem)) {
			return errors.ErrExternalSystemNotAccepted
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, maxWorkersIssueSync)

	syncMutex := &sync.Mutex{}
	statusesMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	// Sync all issues concurrently
	for _, issueData := range issues {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue with Slack sync
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

	return nil
}

func (s *issueService) syncIssue(ctx context.Context, syncMutex *sync.Mutex, issueInput model.IssueData, syncDate time.Time) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueService.syncIssue")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", issueInput.ExternalSystem), log.Object("issueInput", issueInput))

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

	if issueInput.OrganizationRequired && reporterLabel != entity.NodeLabel_Organization {
		reason = fmt.Sprintf("organization(s) not found for issue %s for tenant %s", issueInput.ExternalId, tenant)
		s.log.Warnf("Skip issue sync: %v", reason)
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(reason)
	}

	// Lock issue creation
	syncMutex.Lock()
	// Check if issue already exists
	issueId, err := s.repositories.IssueRepository.GetMatchedIssueId(ctx, tenant, issueInput.ExternalSystem, issueInput.ExternalId)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched issue with external reference %s for tenant %s :%s", issueInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}
	if !failedSync {
		matchingIssueExists := issueId != ""
		span.LogFields(log.Bool("found matching issue", matchingIssueExists))

		// Create new issue id if not found
		issueId = utils.NewUUIDIfEmpty(issueId)
		issueInput.Id = issueId
		span.LogFields(log.String("issueId", issueId))

		// Create or update issue
		issueGrpcRequest := issuegrpc.UpsertIssueGrpcRequest{
			Tenant:      tenant,
			Id:          issueId,
			Subject:     issueInput.Subject,
			Status:      issueInput.Status,
			Priority:    issueInput.Priority,
			Description: issueInput.Description,
			CreatedAt:   utils.ConvertTimeToTimestampPtr(issueInput.CreatedAt),
			UpdatedAt:   utils.ConvertTimeToTimestampPtr(issueInput.UpdatedAt),
			SourceFields: &commongrpc.SourceFields{
				Source:    issueInput.ExternalSystem,
				AppSource: utils.StringFirstNonEmpty(issueInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commongrpc.ExternalSystemFields{
				ExternalSystemId: issueInput.ExternalSystem,
				ExternalId:       issueInput.ExternalId,
				ExternalUrl:      issueInput.ExternalUrl,
				ExternalIdSecond: issueInput.ExternalIdSecond,
				ExternalSource:   issueInput.ExternalSourceEntity,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		}
		if reporterId != "" && reporterLabel == entity.NodeLabel_Organization {
			issueGrpcRequest.ReportedByOrganizationId = &reporterId
		}
		if submitterId != "" {
			switch submitterLabel {
			case entity.NodeLabel_Organization:
				issueGrpcRequest.SubmittedByOrganizationId = &submitterId
			case entity.NodeLabel_User:
				issueGrpcRequest.SubmittedByUserId = &submitterId
			}
		}
		_, err = s.grpcClients.IssueClient.UpsertIssue(ctx, &issueGrpcRequest)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed sending event to upsert issue with external reference %s for tenant %s :%s", issueInput.ExternalId, tenant, err)
			s.log.Error(reason)
		}
		// Wait for issue to be created in neo4j
		if !failedSync && !matchingIssueExists {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				issue, findErr := s.repositories.IssueRepository.GetById(ctx, tenant, issueId)
				if issue != nil && findErr == nil {
					break
				}
				time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
			}
		}
	}
	syncMutex.Unlock()

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
			if followerId != "" && followerLabel == entity.NodeLabel_User && !utils.Contains(processedFollowerUserIds, followerId) {
				_, err = s.grpcClients.IssueClient.AddUserFollower(ctx, &issuegrpc.AddUserFollowerToIssueGrpcRequest{
					Tenant:    common.GetTenantFromContext(ctx),
					IssueId:   issueId,
					UserId:    followerId,
					AppSource: utils.StringFirstNonEmpty(issueInput.AppSource, constants.AppSourceCustomerOsWebhooks),
				})
				processedFollowerUserIds = append(processedFollowerUserIds, followerId)
				if err != nil {
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed sending event to add follower %s to issue %s for tenant %s :%s", followerId, issueId, tenant, err.Error())
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
			if collaboratorId != "" && collaboratorLabel == entity.NodeLabel_User && !utils.Contains(processedFollowerUserIds, collaboratorId) {
				_, err = s.grpcClients.IssueClient.AddUserFollower(ctx, &issuegrpc.AddUserFollowerToIssueGrpcRequest{
					Tenant:    common.GetTenantFromContext(ctx),
					IssueId:   issueId,
					UserId:    collaboratorId,
					AppSource: utils.StringFirstNonEmpty(issueInput.AppSource, constants.AppSourceCustomerOsWebhooks),
				})
				processedFollowerUserIds = append(processedFollowerUserIds, collaboratorId)
				if err != nil {
					tracing.TraceErr(span, err)
					reason = fmt.Sprintf("failed sending event to add follower %s to issue %s for tenant %s :%s", collaboratorId, issueId, tenant, err.Error())
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
			_, err = s.grpcClients.IssueClient.AddUserAssignee(ctx, &issuegrpc.AddUserAssigneeToIssueGrpcRequest{
				Tenant:    common.GetTenantFromContext(ctx),
				IssueId:   issueId,
				UserId:    assigneeId,
				AppSource: utils.StringFirstNonEmpty(issueInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			})
			if err != nil {
				tracing.TraceErr(span, err)
				reason = fmt.Sprintf("failed sending event to add assignee %s to issue %s for tenant %s :%s", assigneeId, issueId, tenant, err.Error())
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

func (s *issueService) mapDbNodeToIssueEntity(dbNode dbtype.Node) *entity.IssueEntity {
	props := utils.GetPropsFromNode(dbNode)
	output := entity.IssueEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Subject:       utils.GetStringPropOrEmpty(props, "subject"),
		Status:        utils.GetStringPropOrEmpty(props, "status"),
		Priority:      utils.GetStringPropOrEmpty(props, "priority"),
		Description:   utils.GetStringPropOrEmpty(props, "description"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &output
}
