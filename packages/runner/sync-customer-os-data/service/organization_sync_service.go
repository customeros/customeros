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

type organizationSyncService struct {
	repositories *repository.Repositories
	cfg          *config.Config
	log          logger.Logger
}

func NewDefaultOrganizationSyncService(repositories *repository.Repositories, cfg *config.Config, log logger.Logger) SyncService {
	return &organizationSyncService{
		repositories: repositories,
		cfg:          cfg,
		log:          log,
	}
}

func (s *organizationSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0

	for {
		organizations := dataService.GetDataForSync(ctx, common.ORGANIZATIONS, batchSize, runId)
		if len(organizations) == 0 {
			break
		}
		s.log.Infof("syncing %d organizations from %s for tenant %s", len(organizations), dataService.SourceId(), tenant)

		var organizationsForWebhooks []entity.OrganizationData
		for _, organization := range organizations {
			inputOrganization := organization.(entity.OrganizationData)
			if inputOrganization.ExternalSystem == "" {
				_ = dataService.MarkProcessed(ctx, inputOrganization.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
				failed++
			} else {
				inputOrganization.AppSource = constants.AppSourceSyncCustomerOsData
				organizationsForWebhooks = append(organizationsForWebhooks, inputOrganization)
			}
		}

		if len(organizationsForWebhooks) > 0 {
			err := s.postOrganizations(ctx, tenant, organizationsForWebhooks)
			if err != nil {
				s.log.Errorf("error while posting log organizations to webhooks: %v", err.Error())
				for _, organizationForWebhooks := range organizationsForWebhooks {
					failed++
					_ = dataService.MarkProcessed(ctx, organizationForWebhooks.SyncId, runId, false, false, "")
				}
			} else {
				s.log.Infof("successfully posted %d log organizations to webhooks", len(organizationsForWebhooks))
				for _, organizationForWebhooks := range organizationsForWebhooks {
					completed++
					_ = dataService.MarkProcessed(ctx, organizationForWebhooks.SyncId, runId, true, false, "")
				}
			}
		}

		if len(organizations) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *organizationSyncService) postOrganizations(ctx context.Context, tenant string, organizations []entity.OrganizationData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationSyncService.postOrganizations")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	// Convert the log entries slice to JSON
	jsonData, err := json.Marshal(organizations)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", s.cfg.Service.CustomerOsWebhooksAPI+"/sync/organizations", bytes.NewBuffer(jsonData))
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
