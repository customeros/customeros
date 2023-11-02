package service

import (
	"context"
	_e "errors"
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
	commentpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/comment"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"sync"
	"time"
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
		maxWorkers:   services.cfg.ConcurrencyConfig.CommentSyncConcurrency,
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
		if !entity.IsValidDataSource(strings.ToLower(comment.ExternalSystem)) {
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
	span.SetTag("externalSystem", commentInput.ExternalSystem)
	span.LogFields(log.Object("commentInput", commentInput), log.Object("syncDate", syncDate))

	var tenant = common.GetTenantFromContext(ctx)
	var failedSync = false
	var reason = ""
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
		reason = fmt.Sprintf("no commented parent entity identified for comment  %s , tenant %s :%s", commentInput.ExternalId, tenant, err.Error())
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
		matchingCommentFound := commentId != ""
		span.LogFields(log.Bool("found matching comment", matchingCommentFound))
		span.LogFields(log.String("commentId", commentId))

		request := commentpb.UpsertCommentGrpcRequest{
			Id:          commentId,
			Tenant:      tenant,
			Content:     commentInput.Content,
			ContentType: commentInput.ContentType,
			CreatedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(commentInput.CreatedAt, utils.NowAsPtr()).(time.Time)),
			UpdatedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(commentInput.UpdatedAt, utils.NowAsPtr()).(time.Time)),
			SourceFields: &commonpb.SourceFields{
				Source:    commentInput.ExternalSystem,
				AppSource: utils.StringFirstNonEmpty(commentInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commonpb.ExternalSystemFields{
				ExternalSystemId: commentInput.ExternalSystem,
				ExternalId:       commentInput.ExternalId,
				ExternalSource:   commentInput.ExternalSourceEntity,
				ExternalUrl:      commentInput.ExternalUrl,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		}
		userAuthorId, _ := s.services.UserService.GetIdForReferencedUser(ctx, tenant, commentInput.ExternalSystem, commentInput.AuthorUser)
		if userAuthorId != "" {
			request.AuthorUserId = utils.StringPtr(userAuthorId)
		}
		if commentedIssueId != "" {
			request.CommentedIssueId = utils.StringPtr(commentedIssueId)
		}
		response, err := s.grpcClients.CommentClient.UpsertComment(ctx, &request)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertComment"))
			reason = fmt.Sprintf("failed sending event to upsert comment with external reference %s for tenant %s :%s", commentInput.ExternalId, tenant, err.Error())
			s.log.Error(reason)
		} else {
			commentId = response.GetId()
		}
		// Wait for comment to be created in neo4j
		if !failedSync && !matchingCommentFound {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				comment, forErr := s.repositories.CommentRepository.GetById(ctx, commentId)
				if comment != nil && forErr == nil {
					break
				}
				time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
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

func (s *commentService) mapDbNodeToCommentEntity(dbNode dbtype.Node) *entity.CommentEntity {
	props := utils.GetPropsFromNode(dbNode)
	return &entity.CommentEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
}
