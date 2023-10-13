package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"net/http"
	"time"
)

type logEntrySyncService struct {
	repositories *repository.Repositories
	cfg          *config.Config
	log          logger.Logger
}

func NewDefaultLogEntrySyncService(repositories *repository.Repositories, cfg *config.Config, log logger.Logger) SyncService {
	return &logEntrySyncService{
		repositories: repositories,
		cfg:          cfg,
		log:          log,
	}
}

func (s *logEntrySyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0

	for {
		logEntries := dataService.GetDataForSync(ctx, common.LOG_ENTRIES, batchSize, runId)
		if len(logEntries) == 0 {
			break
		}
		s.log.Infof("syncing %d log entries from %s for tenant %s", len(logEntries), dataService.SourceId(), tenant)

		var logEntriesForWebhooks []entity.LogEntryData
		for _, logEntry := range logEntries {
			inputLogEntry := logEntry.(entity.LogEntryData)
			if inputLogEntry.ExternalSystem == "" {
				_ = dataService.MarkProcessed(ctx, inputLogEntry.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
				failed++
			} else {
				inputLogEntry.AppSource = constants.AppSourceSyncCustomerOsData
				logEntriesForWebhooks = append(logEntriesForWebhooks, inputLogEntry)
			}
		}

		if len(logEntriesForWebhooks) > 0 {
			err := s.postLogEntries(tenant, logEntriesForWebhooks)
			if err != nil {
				s.log.Errorf("error while posting log entries to webhooks: %v", err.Error())
				for _, logEntryForWebhooks := range logEntriesForWebhooks {
					failed++
					_ = dataService.MarkProcessed(ctx, logEntryForWebhooks.SyncId, runId, false, false, "")
				}
			} else {
				s.log.Infof("successfully posted %d log entries to webhooks", len(logEntriesForWebhooks))
				for _, logEntryForWebhooks := range logEntriesForWebhooks {
					completed++
					_ = dataService.MarkProcessed(ctx, logEntryForWebhooks.SyncId, runId, true, false, "")
				}
			}
		}

		if len(logEntries) < batchSize {
			break
		}
	}

	return completed, failed, skipped
}

func (s *logEntrySyncService) postLogEntries(tenant string, logEntries []entity.LogEntryData) error {
	// Convert the log entries slice to JSON
	jsonData, err := json.Marshal(logEntries)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", s.cfg.Service.CustomerOsWebhooksAPI+"/sync/log-entries", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OPENLINE-API-KEY", s.cfg.Service.CustomerOsWebhooksAPIKey)
	req.Header.Set("tenant", tenant)

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status code to determine if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	return nil
}
