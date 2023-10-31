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
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	logentrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"sync"
	"time"
)

type LogEntryService interface {
	SyncLogEntries(ctx context.Context, logEntries []model.LogEntryData) error
}

type logEntryService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewLogEntryService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) LogEntryService {
	return &logEntryService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.LogEntrySyncConcurrency,
	}
}

func (s *logEntryService) SyncLogEntries(ctx context.Context, logEntries []model.LogEntryData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryService.SyncLogEntries")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return errors.ErrTenantNotValid
	}

	// pre-validate log entry input before syncing
	for _, logEntry := range logEntries {
		if logEntry.ExternalSystem == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return errors.ErrMissingExternalSystem
		}
		if !entity.IsValidDataSource(strings.ToLower(logEntry.ExternalSystem)) {
			tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", logEntry.ExternalSystem))
			return errors.ErrExternalSystemNotAccepted
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

	// Sync all log entries concurrently
	for _, logEntryData := range logEntries {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
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
			statusesMutex.Lock()
			statuses = append(statuses, result)
			statusesMutex.Unlock()
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
	span.SetTag("externalSystem", logEntryInput.ExternalSystem)
	span.LogFields(log.Object("syncDate", syncDate), log.Object("logEntryInput", logEntryInput))

	var failedSync = false
	var reason = ""
	logEntryInput.Normalize()

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, logEntryInput.ExternalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", logEntryInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		span.LogFields(log.String("output", "failed"))
		return NewFailedSyncStatus(reason)
	}

	// Check if log entry sync should be skipped
	if logEntryInput.Skip {
		span.LogFields(log.String("output", "skipped"))
		return NewSkippedSyncStatus(logEntryInput.SkipReason)
	}

	loggedOrgIds := make([]string, 0)
	if logEntryInput.LoggedEntityRequired {
		found := false
		orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.LoggedOrganization)
		if orgId != "" {
			loggedOrgIds = append(loggedOrgIds, orgId)
			found = true
		}
		for _, loggedOrganization := range logEntryInput.LoggedOrganizations {
			orgId, _ = s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, logEntryInput.ExternalSystem, loggedOrganization)
			if orgId != "" {
				loggedOrgIds = append(loggedOrgIds, orgId)
				found = true
			}
		}
		if !found {
			failedSync = true
			reason = fmt.Sprintf("organization not found for log entry %s for tenant %s", logEntryInput.ExternalId, tenant)
			s.log.Error(reason)
			span.LogFields(log.String("output", "failed"))
			return NewFailedSyncStatus(reason)
		}
		loggedOrgIds = utils.RemoveDuplicates(loggedOrgIds)
	}

	// Lock log entry creation
	syncMutex.Lock()
	defer syncMutex.Unlock()
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
		span.LogFields(log.String("logEntryId", logEntryId))

		request := logentrypb.UpsertLogEntryGrpcRequest{
			Id:          logEntryId,
			Tenant:      tenant,
			Content:     logEntryInput.Content,
			ContentType: logEntryInput.ContentType,
			CreatedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.CreatedAt, utils.NowAsPtr()).(time.Time)),
			UpdatedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.UpdatedAt, utils.NowAsPtr()).(time.Time)),
			StartedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.StartedAt, utils.NowAsPtr()).(time.Time)),
			SourceFields: &commonpb.SourceFields{
				Source:    logEntryInput.ExternalSystem,
				AppSource: utils.StringFirstNonEmpty(logEntryInput.AppSource, constants.AppSourceCustomerOsWebhooks),
			},
			ExternalSystemFields: &commonpb.ExternalSystemFields{
				ExternalSystemId: logEntryInput.ExternalSystem,
				ExternalId:       logEntryInput.ExternalId,
				ExternalSource:   logEntryInput.ExternalSourceEntity,
				ExternalUrl:      logEntryInput.ExternalUrl,
				SyncDate:         utils.ConvertTimeToTimestampPtr(&syncDate),
			},
		}
		userAuthorId, _ := s.services.UserService.GetIdForReferencedUser(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.AuthorUser)
		if userAuthorId != "" {
			request.AuthorUserId = utils.StringPtr(userAuthorId)
		}
		if len(loggedOrgIds) == 0 {
			failedSync, reason = s.sendLogEntryToEventStoreForLoggedOrganization(ctx, logEntryId, logEntryInput.ExternalId, "", &request, span, matchingLogEntryExists)
		} else {
			for _, orgId := range loggedOrgIds {
				failedSync, reason = s.sendLogEntryToEventStoreForLoggedOrganization(ctx, logEntryId, logEntryInput.ExternalId, orgId, &request, span, matchingLogEntryExists)
				if failedSync {
					break
				}
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

func (s *logEntryService) sendLogEntryToEventStoreForLoggedOrganization(ctx context.Context, logEntryId, externalId, organizationId string, request *logentrypb.UpsertLogEntryGrpcRequest, span opentracing.Span, matchingLogEntryExists bool) (bool, string) {
	if organizationId != "" {
		request.LoggedOrganizationId = utils.StringPtr(organizationId)
	}
	failedSync := false
	reason := ""
	response, err := s.grpcClients.LogEntryClient.UpsertLogEntry(ctx, request)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err, log.String("grpcMethod", "UpsertLogEntry"))
		reason = fmt.Sprintf("failed sending event to upsert log entry with external reference %s for tenant %s :%s", externalId, common.GetTenantFromContext(ctx), err.Error())
		s.log.Error(reason)
	} else {
		logEntryId = response.GetId()
	}
	// Wait for log entry to be created in neo4j
	if !failedSync && !matchingLogEntryExists {
		for i := 1; i <= constants.MaxRetryCheckDataInNeo4jAfterEventRequest; i++ {
			logEntry, findErr := s.repositories.LogEntryRepository.GetById(ctx, common.GetTenantFromContext(ctx), logEntryId)
			if logEntry != nil && findErr == nil {
				break
			}
			time.Sleep(time.Duration(i*constants.TimeoutIntervalMs) * time.Millisecond)
		}
	}
	return failedSync, reason
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
