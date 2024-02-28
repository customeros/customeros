package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	temporal_client "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/workflows"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func DispatchWebhook(tenant string, event WebhookEvent, payload *InvoicePayload, db *repository.Repositories, cfg config.Config) error {
	// fetch webhook data from db
	webhookResult := db.CommonRepositories.TenantWebhookRepository.GetWebhook(tenant, event.String())
	if webhookResult.Error != nil {
		return fmt.Errorf("error fetching webhook data: %v", webhookResult.Error)
	}

	// if webhook data is not found, return
	if webhookResult.Result == nil {
		return nil
	}

	wh := mapResultToWebhook(webhookResult)

	if wh == nil {
		return nil
	}

	requestBodyJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("(webhook.DispatchWebhook) error marshalling request body: %v", err)
	}
	// TODO: use Temporal here
	tClient, err := temporal_client.TemporalClient(cfg.Temporal.HostPort, cfg.Temporal.Namespace)
	if err != nil {
		return fmt.Errorf("error creating Temporal client: %v", err)
	}
	defer tClient.Close()

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        time.Second * 100, // 100 * InitialInterval
		MaximumAttempts:        3,                 // if set to 0 means Unlimited attempts; not inclusive eg. n < 3.
		NonRetryableErrorTypes: []string{},        // empty
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                       "webhook-calls_" + uuid.New().String(),
		WorkflowExecutionTimeout: time.Hour * 24 * 365 * 10,
		TaskQueue:                workflows.WEBHOOK_CALLS_TASK_QUEUE, // "webhook-calls",
	}

	workflowParams := workflows.WHWorkflowParam{
		TargetUrl:       wh.WebhookUrl,
		RequestBody:     bytes.NewBuffer(requestBodyJSON),
		AuthHeaderName:  wh.AuthHeaderName,
		AuthHeaderValue: wh.AuthHeaderValue,
		RetryPolicy:     retryPolicy,
	}

	_, err = tClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.WebhookWorkflow, workflowParams)
	if err != nil {
		return fmt.Errorf("error executing Temporal workflow: %v", err)
	}
	// Create a POST request with headers and body
	// req, err := http.NewRequest("POST", wh.WebhookUrl, bytes.NewBuffer(requestBodyJSON))
	// if err != nil {
	// 	return fmt.Errorf("error creating request: %v", err)
	// }

	// req.Header.Set("Content-Type", "application/json")
	// authHeader := fmt.Sprintf("%s %s", wh.AuthHeaderValue, wh.ApiKey) // e.g. "Bearer <api_key>"
	// headerName := "Authorization"
	// if wh.AuthHeaderName != "" {
	// 	headerName = wh.AuthHeaderName
	// }
	// req.Header.Set(headerName, authHeader)

	// // Send the request
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return fmt.Errorf("error sending request: %v", err)
	// }
	// defer resp.Body.Close()
	// TODO: Up to here with Temporal

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return nil
}

func mapResultToWebhook(result helper.QueryResult) *entity.TenantWebhook {
	if result.Error != nil {
		return nil
	}
	webhook, ok := result.Result.(*entity.TenantWebhook)
	if !ok {
		return nil
	}
	return webhook
}
