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

type userSyncService struct {
	repositories *repository.Repositories
	cfg          *config.Config
	log          logger.Logger
}

func NewDefaultUserSyncService(repositories *repository.Repositories, cfg *config.Config, log logger.Logger) SyncService {
	return &userSyncService{
		repositories: repositories,
		cfg:          cfg,
		log:          log,
	}
}

func (s *userSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0

	for {
		users := dataService.GetDataForSync(ctx, common.USERS, batchSize, runId)
		if len(users) == 0 {
			break
		}
		s.log.Infof("syncing %d users from %s for tenant %s", len(users), dataService.SourceId(), tenant)

		var usersForWebhooks []entity.UserData
		for _, user := range users {
			inputUser := user.(entity.UserData)
			if inputUser.ExternalSystem == "" {
				_ = dataService.MarkProcessed(ctx, inputUser.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
				failed++
			} else {
				inputUser.AppSource = constants.AppSourceSyncCustomerOsData
				usersForWebhooks = append(usersForWebhooks, inputUser)
			}
		}
		if len(usersForWebhooks) > 0 {
			err := s.postUsers(tenant, usersForWebhooks)
			if err != nil {
				s.log.Errorf("error while posting users to webhooks: %v", err.Error())
				for _, userForWebhooks := range usersForWebhooks {
					failed++
					_ = dataService.MarkProcessed(ctx, userForWebhooks.SyncId, runId, false, false, "")
				}
			} else {
				s.log.Infof("successfully posted %d users to webhooks", len(usersForWebhooks))
				for _, userForWebhooks := range usersForWebhooks {
					completed++
					_ = dataService.MarkProcessed(ctx, userForWebhooks.SyncId, runId, true, false, "")
				}
			}
		}

		if len(users) < batchSize {
			break
		}
	}
	return completed, failed, skipped
}

func (s *userSyncService) postUsers(tenant string, users []entity.UserData) error {
	// Convert the users slice to JSON
	jsonData, err := json.Marshal(users)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", s.cfg.Service.CustomerOsWebhooksAPI+"/sync/users", bytes.NewBuffer(jsonData))
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
