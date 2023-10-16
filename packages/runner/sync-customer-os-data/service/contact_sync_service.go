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

type contactSyncService struct {
	repositories *repository.Repositories
	cfg          *config.Config
	log          logger.Logger
}

func NewDefaultContactSyncService(repositories *repository.Repositories, cfg *config.Config, log logger.Logger) SyncService {
	return &contactSyncService{
		repositories: repositories,
		cfg:          cfg,
		log:          log,
	}
}

func (s *contactSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0

	for {
		contacts := dataService.GetDataForSync(ctx, common.CONTACTS, batchSize, runId)
		if len(contacts) == 0 {
			break
		}
		s.log.Infof("syncing %d contacts from %s for tenant %s", len(contacts), dataService.SourceId(), tenant)

		var contactsForWebhooks []entity.ContactData
		for _, contact := range contacts {
			inputContact := contact.(entity.ContactData)
			if inputContact.ExternalSystem == "" {
				_ = dataService.MarkProcessed(ctx, inputContact.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
				failed++
			} else {
				inputContact.AppSource = constants.AppSourceSyncCustomerOsData
				contactsForWebhooks = append(contactsForWebhooks, inputContact)
			}
		}

		if len(contactsForWebhooks) > 0 {
			err := s.postContacts(ctx, tenant, contactsForWebhooks)
			if err != nil {
				s.log.Errorf("error while posting log contacts to webhooks: %v", err.Error())
				for _, contactForWebhooks := range contactsForWebhooks {
					failed++
					_ = dataService.MarkProcessed(ctx, contactForWebhooks.SyncId, runId, false, false, "")
				}
			} else {
				s.log.Infof("successfully posted %d log contacts to webhooks", len(contactsForWebhooks))
				for _, contactForWebhooks := range contactsForWebhooks {
					completed++
					_ = dataService.MarkProcessed(ctx, contactForWebhooks.SyncId, runId, true, false, "")
				}
			}
		}

		if len(contacts) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *contactSyncService) postContacts(ctx context.Context, tenant string, contacts []entity.ContactData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactSyncService.postContacts")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	// Convert the log entries slice to JSON
	jsonData, err := json.Marshal(contacts)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", s.cfg.Service.CustomerOsWebhooksAPI+"/sync/contacts", bytes.NewBuffer(jsonData))
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
