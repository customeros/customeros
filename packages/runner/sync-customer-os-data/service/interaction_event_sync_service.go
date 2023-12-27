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

type interactionEventSyncService struct {
	repositories *repository.Repositories
	cfg          *config.Config
	log          logger.Logger
}

func NewDefaultInteractionEventSyncService(repositories *repository.Repositories, cfg *config.Config, log logger.Logger) SyncService {
	return &interactionEventSyncService{
		repositories: repositories,
		log:          log,
		cfg:          cfg,
	}
}

func (s *interactionEventSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0

	for {
		interactionEvents := dataService.GetDataForSync(ctx, common.INTERACTION_EVENTS, batchSize, runId)
		if len(interactionEvents) == 0 {
			break
		}
		s.log.Infof("syncing %d interaction events from %s for tenant %s", len(interactionEvents), dataService.SourceId(), tenant)

		var interactionEventsForWebhooks []entity.InteractionEventData
		for _, event := range interactionEvents {
			inputInteractionEvent := event.(entity.InteractionEventData)
			if inputInteractionEvent.ExternalSystem == "" {
				_ = dataService.MarkProcessed(ctx, inputInteractionEvent.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
				failed++
			} else {
				inputInteractionEvent.AppSource = constants.AppSourceSyncCustomerOsData
				interactionEventsForWebhooks = append(interactionEventsForWebhooks, inputInteractionEvent)
			}
		}

		if len(interactionEventsForWebhooks) > 0 {
			err := s.postInteractionEvents(ctx, tenant, interactionEventsForWebhooks)
			if err != nil {
				s.log.Errorf("error while posting interaction events to webhooks: %v", err.Error())
				for _, interactionEventForWebhooks := range interactionEventsForWebhooks {
					failed++
					_ = dataService.MarkProcessed(ctx, interactionEventForWebhooks.SyncId, runId, false, false, "")
				}
			} else {
				s.log.Infof("successfully posted %d interaction events to webhooks", len(interactionEventsForWebhooks))
				for _, interactionEventForWebhooks := range interactionEventsForWebhooks {
					completed++
					_ = dataService.MarkProcessed(ctx, interactionEventForWebhooks.SyncId, runId, true, false, "")
				}
			}
		}

		if len(interactionEvents) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *interactionEventSyncService) postInteractionEvents(ctx context.Context, tenant string, interactionEvents []entity.InteractionEventData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventSyncService.postInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	// Convert the log entries slice to JSON
	jsonData, err := json.Marshal(interactionEvents)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", s.cfg.Service.CustomerOsWebhooksAPI+"/sync/interaction-events", bytes.NewBuffer(jsonData))
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
