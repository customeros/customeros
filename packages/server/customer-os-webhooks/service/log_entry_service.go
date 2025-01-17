package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
)

type LogEntryService interface {
	SyncLogEntries(ctx context.Context, logEntries []model.LogEntryData) (SyncResult, error)
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
		maxWorkers:   services.cfg.App.ConcurrencyConfig.LogEntrySyncConcurrency,
	}
}

func (s *logEntryService) SyncLogEntries(ctx context.Context, logEntries []model.LogEntryData) (SyncResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryService.SyncLogEntries")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate log entry input before syncing
	for _, logEntry := range logEntries {
		if logEntry.ExternalSystem == "" {
			tracing.TraceErr(span, errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if !neo4jentity.IsValidDataSource(strings.ToLower(logEntry.ExternalSystem)) {
			tracing.TraceErr(span, errors.ErrExternalSystemNotAccepted, log.String("externalSystem", logEntry.ExternalSystem))
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

	// Sync all log entries concurrently
	for _, logEntryData := range logEntries {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return SyncResult{}, ctx.Err()
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

	return s.services.SyncStatusService.PrepareSyncResult(statuses), nil
}

func (s *logEntryService) syncLogEntry(ctx context.Context, syncMutex *sync.Mutex, logEntryInput model.LogEntryData, syncDate time.Time, tenant string) SyncStatus {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryService.syncLogEntry")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagExternalSystem, logEntryInput.ExternalSystem)
	span.SetTag(tracing.SpanTagExternalId, logEntryInput.ExternalId)
	span.LogFields(log.Object("syncDate", syncDate))
	tracing.LogObjectAsJson(span, "logEntryInput", logEntryInput)

	failedSync := false
	reason := ""
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

		logEntryFields := data_fields.LogEntryFields{
			Content:     utils.StringPtr(logEntryInput.Content),
			ContentType: utils.StringPtr(logEntryInput.ContentType),
			ExternalSystem: &neo4jmodel.ExternalSystem{
				ExternalSystemId: logEntryInput.ExternalSystem,
				ExternalId:       logEntryInput.ExternalId,
				ExternalSource:   logEntryInput.ExternalSourceEntity,
				ExternalUrl:      logEntryInput.ExternalUrl,
				SyncDate:         &syncDate,
			},
			Source:    utils.StringPtr(logEntryInput.ExternalSystem),
			AppSource: utils.StringPtr(utils.StringFirstNonEmpty(logEntryInput.AppSource, constants.AppSourceCustomerOsWebhooks)),
			CreatedAt: logEntryInput.CreatedAt,
			StartedAt: logEntryInput.StartedAt,
		}

		userAuthorId, _ := s.services.UserService.GetIdForReferencedUser(ctx, tenant, logEntryInput.ExternalSystem, logEntryInput.AuthorUser)
		if userAuthorId != "" {
			logEntryFields.AuthorUserId = utils.StringPtr(userAuthorId)
		}
		if len(loggedOrgIds) == 0 {
			failedSync, reason = s.saveLogEntryToDb(ctx, logEntryId, logEntryInput.ExternalId, "", logEntryFields, span, matchingLogEntryExists)
		} else {
			for _, orgId := range loggedOrgIds {
				failedSync, reason = s.saveLogEntryToDb(ctx, logEntryId, logEntryInput.ExternalId, orgId, logEntryFields, span, matchingLogEntryExists)
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

func (s *logEntryService) saveLogEntryToDb(ctx context.Context, logEntryId, externalId, organizationId string, logEntryFields data_fields.LogEntryFields, span opentracing.Span, matchingLogEntryExists bool) (bool, string) {
	if organizationId != "" {
		logEntryFields.OrganizationId = utils.StringPtr(organizationId)
	}
	failedSync := false
	reason := ""
	logEntryId, err := s.services.CommonServices.LogEntryService.Save(ctx, &logEntryId, logEntryFields)
	if err != nil {
		failedSync = true
		tracing.TraceErr(span, err)
		reason = fmt.Sprintf("error saving log entry with external reference %s for tenant %s :%s", externalId, common.GetTenantFromContext(ctx), err.Error())
		s.log.Error(reason)
	}

	return failedSync, reason
}
