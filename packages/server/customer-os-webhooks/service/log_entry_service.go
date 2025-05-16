package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/constants"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/errors"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/model"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
)

type LogEntryService interface {
	SyncLogEntries(ctx context.Context, logEntries []model.LogEntryData) (SyncResult, error)
}

type logEntryService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
	maxWorkers   int
}

func NewLogEntryService(log logger.Logger, repositories *repository.Repositories, services *Services) LogEntryService {
	return &logEntryService{
		log:          log,
		repositories: repositories,
		services:     services,
		maxWorkers:   services.Cfg.App.ConcurrencyConfig.LogEntrySyncConcurrency,
	}
}

func (s *logEntryService) SyncLogEntries(ctx context.Context, logEntries []model.LogEntryData) (SyncResult, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "LogEntryService.SyncLogEntries")
	defer spans.Finish()

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		spans.TraceError(errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	// pre-validate log entry input before syncing
	for _, logEntry := range logEntries {
		if logEntry.ExternalSystem == "" {
			spans.TraceError(errors.ErrMissingExternalSystem)
			return SyncResult{}, errors.ErrMissingExternalSystem
		}
		if !neo4jentity.IsValidDataSource(strings.ToLower(logEntry.ExternalSystem)) {
			spans.TraceError(errors.ErrExternalSystemNotAccepted)
			spans.LogKV("externalSystem", logEntry.ExternalSystem)
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "LogEntryService.syncLogEntry")
	defer spans.Finish()
	spans.TagString(telemetry.SpanTagExternalSystem, logEntryInput.ExternalSystem)
	spans.TagString(telemetry.SpanTagExternalId, logEntryInput.ExternalId)
	spans.LogObjectAsJson("syncDate", syncDate)
	spans.LogObjectAsJson("logEntryInput", logEntryInput)

	failedSync := false
	reason := ""
	logEntryInput.Normalize()

	err := s.services.ExternalSystemService.MergeExternalSystem(ctx, tenant, logEntryInput.ExternalSystem)
	if err != nil {
		spans.TraceError(err)
		reason = fmt.Sprintf("failed merging external system %s for tenant %s :%s", logEntryInput.ExternalSystem, tenant, err.Error())
		s.log.Error(reason)
		spans.LogKV("output", "failed")
		return NewFailedSyncStatus(reason)
	}

	// Check if log entry sync should be skipped
	if logEntryInput.Skip {
		spans.LogKV("output", "skipped")
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
			spans.LogKV("output", "failed")
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
		spans.TraceError(err)
		reason = fmt.Sprintf("failed finding existing matched log entru with external reference %s for tenant %s :%s", logEntryInput.ExternalId, tenant, err.Error())
		s.log.Error(reason)
	}

	if !failedSync {
		matchingLogEntryExists := logEntryId != ""
		spans.LogKV("found matching log entry", matchingLogEntryExists)
		spans.LogKV("logEntryId", logEntryId)

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
			failedSync, reason = s.saveLogEntryToDb(ctx, logEntryId, logEntryInput.ExternalId, "", logEntryFields, spans, matchingLogEntryExists)
		} else {
			for _, orgId := range loggedOrgIds {
				failedSync, reason = s.saveLogEntryToDb(ctx, logEntryId, logEntryInput.ExternalId, orgId, logEntryFields, spans, matchingLogEntryExists)
				if failedSync {
					break
				}
			}
		}
	}

	spans.LogKV("failedSync", failedSync)
	if failedSync {
		spans.LogKV("output", "failed")
		return NewFailedSyncStatus(reason)
	}
	spans.LogKV("output", "success")
	return NewSuccessfulSyncStatus()
}

func (s *logEntryService) saveLogEntryToDb(ctx context.Context, logEntryId, externalId, organizationId string, logEntryFields data_fields.LogEntryFields, spans *telemetry.Spans, matchingLogEntryExists bool) (bool, string) {
	if organizationId != "" {
		logEntryFields.OrganizationId = utils.StringPtr(organizationId)
	}
	failedSync := false
	reason := ""
	logEntryId, err := s.services.CommonServices.LogEntryService.Save(ctx, &logEntryId, logEntryFields)
	if err != nil {
		failedSync = true
		spans.TraceError(err)
		reason = fmt.Sprintf("error saving log entry with external reference %s for tenant %s :%s", externalId, common.GetTenantFromContext(ctx), err.Error())
		s.log.Error(reason)
	}

	return failedSync, reason
}
