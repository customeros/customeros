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
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"time"
)

type issueSyncService struct {
	repositories *repository.Repositories
	cfg          *config.Config
	log          logger.Logger
}

func NewDefaultIssueSyncService(repositories *repository.Repositories, cfg *config.Config, log logger.Logger) SyncService {
	return &issueSyncService{
		repositories: repositories,
		cfg:          cfg,
		log:          log,
	}
}

func (s *issueSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0

	for {
		issues := dataService.GetDataForSync(ctx, common.CONTACTS, batchSize, runId)
		if len(issues) == 0 {
			break
		}
		s.log.Infof("syncing %d issues from %s for tenant %s", len(issues), dataService.SourceId(), tenant)

		var issuesForWebhooks []entity.IssueData
		for _, issue := range issues {
			inputIssue := issue.(entity.IssueData)
			if inputIssue.ExternalSystem == "" {
				_ = dataService.MarkProcessed(ctx, inputIssue.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
				failed++
			} else {
				inputIssue.AppSource = constants.AppSourceSyncCustomerOsData
				issuesForWebhooks = append(issuesForWebhooks, inputIssue)
			}
		}

		if len(issuesForWebhooks) > 0 {
			err := s.postIssues(ctx, tenant, issuesForWebhooks)
			if err != nil {
				s.log.Errorf("error while posting issues to webhooks: %v", err.Error())
				for _, issueForWebhooks := range issuesForWebhooks {
					failed++
					_ = dataService.MarkProcessed(ctx, issueForWebhooks.SyncId, runId, false, false, "")
				}
			} else {
				s.log.Infof("successfully posted %d issues to webhooks", len(issuesForWebhooks))
				for _, issueForWebhooks := range issuesForWebhooks {
					completed++
					_ = dataService.MarkProcessed(ctx, issueForWebhooks.SyncId, runId, true, false, "")
				}
			}
		}

		if len(issues) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *issueSyncService) postIssues(ctx context.Context, tenant string, issues []entity.IssueData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueSyncService.postIssues")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	// Convert the log entries slice to JSON
	jsonData, err := json.Marshal(issues)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", s.cfg.Service.CustomerOsWebhooksAPI+"/sync/issues", bytes.NewBuffer(jsonData))
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
