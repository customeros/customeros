package service

import (
	"context"
	_e "errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/errors"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/model"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
)

type CommentService interface {
	SyncComments(ctx context.Context, comments []model.CommentData) (SyncResult, error)
}

type commentService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewCommentService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) CommentService {
	return &commentService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.App.ConcurrencyConfig.CommentSyncConcurrency,
	}
}

func (s *commentService) SyncComments(ctx context.Context, comments []model.CommentData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentService.SyncComments")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate comment input before syncing
	for _, comment := range comments {
		if comment.ExternalSystem == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if !neo4jentity.IsValidDataSource(strings.ToLower(comment.ExternalSystem)) {
			tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", comment.ExternalSystem))
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

	// Sync all comments
	for _, commentData := range comments {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
		default:
		}

		// Acquire a worker slot
		workerLimit <- struct{}{}
		wg.Add(1)

		go func(syncComment model.CommentData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncComment(ctx, syncMutex, syncComment, syncDate)
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
		}(commentData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), comments[0].ExternalSystem,
		comments[0].AppSource, "comment", syncDate, statuses)

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *commentService) syncComment(ctx context.Context, syncMutex *sync.Mutex, commentInput model.CommentData, syncDate time.Time) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentService.syncComment")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagExternalSystem, commentInput.ExternalSystem)
	span.SetTag(tracing.SpanTagExternalId, commentInput.ExternalId)
	span.LogFields(log.Object("syncDate", syncDate))
	tracing.LogObjectAsJson(span, "commentInput", commentInput)

	tenant := common.GetTenantFromContext(ctx)
	failedSync := false
	reason := ""
	commentInput.Normalize()

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, commentInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err, log.String("externalSystem", commentInput.ExternalSystem))
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", commentInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	// Check if comment sync should be skipped
	if commentInput.Skip {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(commentInput.SkipReason)
	}

	commentedIssueId, err := s.services.IssueService.GetIdForReferencedIssue(ctx, tenant, commentInput.ExternalSystem, commentInput.CommentedIssue)
	if err != nil {
		tracing.TraceErr(span, err, log.String("commentedIssue", commentInput.CommentedIssue.ExternalId))
		s.log.Error(reason)
	}
	if commentedIssueId == "" {
		reason = fmt.Sprintf("no commented parent entity identified for comment %v , tenant %v", commentInput.ExternalId, tenant)
		tracing.TraceErr(span, _e.New(reason))
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	// Lock comment creation
	syncMutex.Lock()
	defer syncMutex.Unlock()
	// Check if comment already exists
	commentId, err := s.repositories.CommentRepository.GetMatchedCommentId(ctx, commentInput.ExternalSystem, commentInput.ExternalId)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched log entru with external reference %s for tenant %s :%s", commentInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}

	if !failedSync {
		span.LogFields(log.Bool("found matching comment", commentId != ""))

		commentFields := data_fields.CommentFields{
			Source:           utils.StringPtr(commentInput.ExternalSystem),
			CreatedAt:        commentInput.CreatedAt,
			CommentedIssueId: &commentedIssueId,
			ExternalSystem: &neo4jmodel.ExternalSystem{
				ExternalSystemId: commentInput.ExternalSystem,
				ExternalId:       commentInput.ExternalId,
				ExternalIdSecond: commentInput.ExternalIdSecond,
				ExternalSource:   commentInput.ExternalSourceEntity,
				ExternalUrl:      commentInput.ExternalUrl,
				SyncDate:         &syncDate,
			},
		}

		if commentInput.Content != "" {
			commentFields.Content = &commentInput.Content
		}
		if commentInput.ContentType != "" {
			commentFields.ContentType = &commentInput.ContentType
		}
		userAuthorId, _ := s.services.UserService.GetIdForReferencedUser(ctx, tenant, commentInput.ExternalSystem, commentInput.AuthorUser)
		if userAuthorId != "" {
			commentFields.AuthorUserId = &userAuthorId
		}

		commentId, err = s.services.CommonServices.CommentService.Save(ctx, nil, &commentId, commentFields)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertComment"))
			reason = fmt.Sprintf("error saving comment with external reference %s for tenant %s :%s", commentInput.ExternalId, tenant, err.Error())
			s.log.Error(reason)
		}
		commentInput.Id = commentId
		tracing.TagEntity(span, commentId)
	}

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}
	span.LogFields(log.String("output", "success"))
	return NewSuccessfulSyncStatus()
}
