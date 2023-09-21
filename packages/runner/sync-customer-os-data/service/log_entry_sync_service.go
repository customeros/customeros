package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	grpc_common "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	log_entry_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/log_entry"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

type logEntrySyncService struct {
	repositories *repository.Repositories
	services     *Services
	grpcClients  *grpc_client.Clients
	log          logger.Logger
}

func NewDefaultLogEntrySyncService(repositories *repository.Repositories, services *Services, grpcClients *grpc_client.Clients, log logger.Logger) SyncService {
	return &logEntrySyncService{
		repositories: repositories,
		services:     services,
		grpcClients:  grpcClients,
		log:          log,
	}
}

func (s *logEntrySyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0
	logEntrySyncMutex := &sync.Mutex{}

	for {

		logEntries := dataService.GetDataForSync(ctx, common.LOG_ENTRIES, batchSize, runId)

		if len(logEntries) == 0 {
			break
		}

		s.log.Infof("syncing %d log entries from %s for tenant %s", len(logEntries), dataService.SourceId(), tenant)

		var wg sync.WaitGroup
		wg.Add(len(logEntries))

		results := make(chan result, len(logEntries))
		done := make(chan struct{})

		for _, v := range logEntries {
			v := v

			go func(logEntry entity.LogEntryData) {
				defer wg.Done()

				var comp, fail, skip int
				s.syncLogEntry(ctx, logEntrySyncMutex, logEntry, dataService, syncDate, tenant, runId, &comp, &fail, &skip)

				results <- result{comp, fail, skip}
			}(v.(entity.LogEntryData))
		}
		// Wait for goroutines to finish
		go func() {
			wg.Wait()
			close(done)
		}()
		go func() {
			<-done
			close(results)
		}()

		for r := range results {
			completed += r.completed
			failed += r.failed
			skipped += r.skipped
		}

		if len(logEntries) < batchSize {
			break
		}
	}

	return completed, failed, skipped
}

func (s *logEntrySyncService) syncLogEntry(ctx context.Context, logEntrySyncMutex *sync.Mutex, logEntryInput entity.LogEntryData, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, completed, failed, skipped *int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntrySyncService.syncLogEntry")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	var failedSync = false
	var reason string
	logEntryInput.Normalize()

	if logEntryInput.ExternalSystem == "" {
		_ = dataService.MarkProcessed(ctx, logEntryInput.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
		*failed++
		return
	}

	if logEntryInput.Skip {
		if err := dataService.MarkProcessed(ctx, logEntryInput.SyncId, runId, true, true, logEntryInput.SkipReason); err != nil {
			*failed++
			span.LogFields(log.Bool("failedSync", true))
			return
		}
		*skipped++
		span.LogFields(log.Bool("skippedSync", true))
		return
	}

	if logEntryInput.LoggedEntityRequired {
		orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.LoggedOrganization)
		if orgId == "" {
			_ = dataService.MarkProcessed(ctx, logEntryInput.SyncId, runId, false, false, "Logged organization not found.")
			*failed++
			return
		}
	}

	logEntrySyncMutex.Lock()
	logEntryId, err := s.repositories.LogEntryRepository.GetMatchedLogEntryId(ctx, tenant, logEntryInput)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("failed searching existing matched log entry with external reference id %s for tenant %s :%s", logEntryInput.ExternalId, tenant, err.Error())
		s.log.Errorf(reason)
	}

	request := log_entry_grpc_service.UpsertLogEntryGrpcRequest{
		Id:          logEntryId,
		Tenant:      tenant,
		Content:     logEntryInput.Content,
		ContentType: logEntryInput.ContentType,
		CreatedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.CreatedAt, utils.NowAsPtr()).(time.Time)),
		UpdatedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.UpdatedAt, utils.NowAsPtr()).(time.Time)),
		StartedAt:   timestamppb.New(utils.TimePtrFirstNonNilNillableAsAny(logEntryInput.StartedAt, utils.NowAsPtr()).(time.Time)),
		SourceFields: &grpc_common.SourceFields{
			Source:        logEntryInput.ExternalSystem,
			SourceOfTruth: logEntryInput.ExternalSystem,
			AppSource:     constants.AppSourceSyncCustomerOsData,
		},
		ExternalSystemFields: &grpc_common.ExternalSystemFields{
			ExternalSystemId: logEntryInput.ExternalSystem,
			ExternalId:       logEntryInput.ExternalId,
			SyncDate:         timestamppb.Now(),
			ExternalSource:   utils.IfNotNilString(logEntryInput.ExternalSourceTable),
			ExternalUrl:      logEntryInput.ExternalUrl,
		},
	}
	orgId, _ := s.services.OrganizationService.GetIdForReferencedOrganization(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.LoggedOrganization)
	if orgId != "" {
		request.LoggedOrganizationId = utils.StringPtr(orgId)
	}
	userId, _ := s.services.UserService.GetIdForReferencedUser(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.AuthorUser)
	if userId != "" {
		request.AuthorUserId = utils.StringPtr(userId)
	}
	response, err := s.grpcClients.LogEntryClient.UpsertLogEntry(ctx, &request)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		s.log.Errorf(fmt.Sprintf("failed upsert log entry with external reference id %s for tenant %s :%s", logEntryInput.ExternalId, tenant, err.Error()))
	} else {
		if response.GetId() != "" {
			logEntryId = response.GetId()
		}
	}
	if logEntryId == "" {
		failedSync = true
	} else {
		maxRetries := 30
		for i := 0; i < maxRetries; i++ {
			id, err := s.repositories.LogEntryRepository.GetLogEntryIdById(ctx, tenant, logEntryId)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf(fmt.Sprintf("failed searching existing log entry with id %s for tenant %s :%s", logEntryId, tenant, err.Error()))
				failedSync = true
				break
			}
			if id != "" {
				break
			}
			if i == maxRetries {
				failedSync = true
				break
			}
		}
	}
	logEntrySyncMutex.Unlock()

	if failedSync == false {
		s.log.Debugf("successfully merged log entry with id %v for tenant %v from %v", logEntryId, tenant, dataService.SourceId())
	}
	if err := dataService.MarkProcessed(ctx, logEntryInput.SyncId, runId, failedSync == false, false, reason); err != nil {
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
