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
	logentrygrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"sync"
	"time"
)

const maxWorkersLogEntrySync = 10

type LogEntryService interface {
	SyncLogEntries(ctx context.Context, logEntries []model.LogEntryData) error
	mapDbNodeToLogEntryEntity(dbNode dbtype.Node) *entity.LogEntryEntity
}

type logEntryService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewLogEntryService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) LogEntryService {
	return &logEntryService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *logEntryService) SyncLogEntries(ctx context.Context, logEntries []model.LogEntryData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryService.SyncLogEntries")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		return errors.ErrTenantNotValid
	}

	// pre-validate log entry input before syncing
	for _, logEntry := range logEntries {
		if logEntry.ExternalSystem == "" {
			return errors.ErrMissingExternalSystem
		}
		if !entity.IsValidDataSource(strings.ToLower(logEntry.ExternalSystem)) {
			return errors.ErrExternalSystemNotAccepted
		}
	}

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup
	// Create a channel to control the number of concurrent workers
	workerLimit := make(chan struct{}, maxWorkersLogEntrySync)

	syncMutex := &sync.Mutex{}
	syncDate := utils.Now()
	var statuses []SyncStatus

	// Sync all log entries concurrently
	for _, logEntryData := range logEntries {
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

		go func(syncLogEntry model.LogEntryData) {
			defer wg.Done()
			defer func() {
				// Release the worker slot when done
				<-workerLimit
			}()

			result := s.syncLogEntry(ctx, syncMutex, syncLogEntry, syncDate, common.GetTenantFromContext(ctx))
			statuses = append(statuses, result)
		}(logEntryData)
	}
	// Wait for all workers to finish
	wg.Wait()

	s.services.SyncStatusService.SaveSyncResults(ctx, common.GetTenantFromContext(ctx), logEntries[0].ExternalSystem,
		logEntries[0].AppSource, "logEntry", syncDate, statuses)

	return nil
}

func (s *logEntryService) syncLogEntry(ctx context.Context, syncMutex *sync.Mutex, logEntryInput model.LogEntryData, syncDate time.Time, tenant string) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryService.syncLogEntry")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", logEntryInput.ExternalSystem), log.Object("logEntryInput", logEntryInput), log.String("tenant", tenant))

	var failedSync = false
	var reason = ""
	logEntryInput.Normalize()

	// TODO: Merge external system, should be cached and moved to external system service
	err := s.repositories.ExternalSystemRepository.MergeExternalSystem(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", logEntryInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		return NewFailedSyncStatus(reason)
	}

	// Check if log entry sync should be skipped
	if logEntryInput.Skip {
		span.LogFields(log.Bool("skippedSync", true))
		return NewSkippedSyncStatus(logEntryInput.SkipReason)
	}

	if logEntryInput.LoggedEntityRequired {
		orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.LoggedOrganization)
		if orgId == "" {
			failedSync = true
			reason = fmt.Sprintf("organization not found for log entry %s for tenant %s", logEntryInput.ExternalId, tenant)
			s.log.Error(reason)
			return NewFailedSyncStatus(reason)
		}
	}

	// Lock log entry creation
	syncMutex.Lock()
	// Check if log entry already exists
	logEntryId, err := s.repositories.LogEntryRepository.GetMatchedLogEntryId(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.ExternalId)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed finding existing matched log entru with external reference %s for tenant %s :%s", logEntryInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}

	if !failedSync {
		matchingLogEntryExists := logEntryId != ""
		span.LogFields(log.Bool("found matching log entry", matchingLogEntryExists))

		request := logentrygrpc.UpsertLogEntryGrpcRequest{
			Id:          logEntryId,
			Tenant:      tenant,
			Content:     logEntryInput.Content,
			ContentType: logEntryInput.ContentType,
			CreatedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.CreatedAt, utils.NowAsPtr()).(time.Time)),
			UpdatedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.UpdatedAt, utils.NowAsPtr()).(time.Time)),
			StartedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.StartedAt, utils.NowAsPtr()).(time.Time)),
			SourceFields: &commongrpc.SourceFields{
				Source:    logEntryInput.ExternalSystem,
				AppSource: utils.StringFirstNonEmpty(logEntryInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commongrpc.ExternalSystemFields{
				ExternalSystemId: logEntryInput.ExternalSystem,
				ExternalId:       logEntryInput.ExternalId,
				ExternalSource:   logEntryInput.ExternalSourceEntity,
				ExternalUrl:      logEntryInput.ExternalUrl,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		}
		orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.LoggedOrganization)
		if orgId != "" {
			request.LoggedOrganizationId = utils.StringPtr(orgId)
		}
		userAuthorId, _ := s.services.UserService.GetIdForReferencedUser(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.AuthorUser)
		if userAuthorId != "" {
			request.AuthorUserId = utils.StringPtr(userAuthorId)
		}
		response, err := s.grpcClients.LogEntryClient.UpsertLogEntry(ctx, &request)
		if err != nil {
			failedSync = true
			tracing.TraceErr(span, err)
			reason = fmt.Sprintf("failed sending event to upsert log entry with external reference %s for tenant %s :%s", logEntryInput.ExternalId, tenant, err.Error())
			s.log.Error(reason)
		} else {
			logEntryId = response.GetId()
		}
		span.LogFields(log.String("logEntryId", logEntryId))
		// Wait for log entry to be created in neo4j
		if !failedSync && !matchingLogEntryExists {
			for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
				logEntry, findErr := s.repositories.LogEntryRepository.GetById(ctx, tenant, logEntryId)
				if logEntry != nil && findErr == nil {
					break
				}
				time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
			}
		}
	}
	syncMutex.Unlock()

	span.LogFields(log.Bool("failedSync", failedSync))
	if failedSync {
		return NewFailedSyncStatus(reason)
	}
	return NewSuccessfulSyncStatus()
}

func (s *logEntryService) mapDbNodeToLogEntryEntity(dbNode dbtype.Node) *entity.LogEntryEntity {
	props := utils.GetPropsFromNode(dbNode)
	return &entity.LogEntryEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrEpochStart(props, "startedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
}
