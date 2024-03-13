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

type orderSyncService struct {
	repositories *repository.Repositories
	cfg          *config.Config
	log          logger.Logger
}

func NewDefaultOrderSyncService(repositories *repository.Repositories, cfg *config.Config, log logger.Logger) SyncService {
	return &orderSyncService{
		repositories: repositories,
		log:          log,
		cfg:          cfg,
	}
}

func (s *orderSyncService) Sync(ctx context.Context, dataService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int) {
	completed, failed, skipped := 0, 0, 0

	for {
		orders := dataService.GetDataForSync(ctx, common.ORDERS, batchSize, runId)
		if len(orders) == 0 {
			break
		}
		s.log.Infof("syncing %d orders from %s for tenant %s", len(orders), dataService.SourceId(), tenant)

		var ordersForWebhooks []entity.OrderData
		for _, event := range orders {
			order := event.(entity.OrderData)
			if order.ExternalSystem == "" {
				_ = dataService.MarkProcessed(ctx, order.SyncId, runId, false, false, "External system is empty. Error during reading data from source")
				failed++
			} else {
				order.AppSource = constants.AppSourceSyncCustomerOsData
				ordersForWebhooks = append(ordersForWebhooks, order)
			}
		}

		if len(ordersForWebhooks) > 0 {
			err := s.postOrders(ctx, tenant, ordersForWebhooks)
			if err != nil {
				s.log.Errorf("error while posting orders to webhooks: %v", err.Error())
				for _, orderForWebhooks := range ordersForWebhooks {
					failed++
					_ = dataService.MarkProcessed(ctx, orderForWebhooks.SyncId, runId, false, false, "")
				}
			} else {
				s.log.Infof("successfully posted %d orders to webhooks", len(ordersForWebhooks))
				for _, orderForWebhooks := range ordersForWebhooks {
					completed++
					_ = dataService.MarkProcessed(ctx, orderForWebhooks.SyncId, runId, true, false, "")
				}
			}
		}

		if len(orders) < batchSize {
			break
		}

	}

	return completed, failed, skipped
}

func (s *orderSyncService) postOrders(ctx context.Context, tenant string, orders []entity.OrderData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrderSyncService.postOrders")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)

	// Convert the log entries slice to JSON
	jsonData, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", s.cfg.Service.CustomerOsWebhooksAPI+"/sync/orders", bytes.NewBuffer(jsonData))
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
