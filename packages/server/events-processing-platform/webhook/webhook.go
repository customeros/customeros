package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/notifications"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	temporal_client "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/workflows"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func DispatchWebhook(tenant string, event WebhookEvent, payload *InvoicePayload, db *repository.Repositories, cfg config.Config, notificationProvider notifications.NotificationProvider, failureNotify bool) error {
	if !cfg.Temporal.RunWorker {
		return fmt.Errorf("temporal worker is not running")
	}

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
	// Start Temporal Client to queue webhook workflow
	tClient, err := temporal_client.TemporalClient(cfg.Temporal.HostPort, cfg.Temporal.Namespace)
	if err != nil {
		return fmt.Errorf("error creating Temporal client: %v", err)
	}
	defer tClient.Close()

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        time.Second * 259200, // 259200 seconds = 3 days
		MaximumAttempts:        0,                    // if set to 0 means Unlimited attempts; this number is not inclusive eg. n < 3.
		NonRetryableErrorTypes: []string{},           // empty
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                       "webhook-calls_" + uuid.New().String(),
		WorkflowExecutionTimeout: time.Hour * 24 * 3,                 // timeout after 3 days
		TaskQueue:                workflows.WEBHOOK_CALLS_TASK_QUEUE, // "webhook-calls",
	}

	var notification *notifications.NovuNotification
	if failureNotify {
		notification = populateNotification(tenant, event.String(), wh)
	}

	workflowParams := workflows.WHWorkflowParam{
		TargetUrl:            wh.WebhookUrl,
		RequestBody:          string(requestBodyJSON),
		AuthHeaderName:       wh.AuthHeaderName,
		AuthHeaderValue:      wh.AuthHeaderValue,
		RetryPolicy:          retryPolicy,
		Notification:         notification,
		NotificationProvider: notificationProvider,
		NotifyFailure:        failureNotify,
	}

	// the workflow will run async, so we don't need to wait for it to finish
	_, err = tClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.WebhookWorkflow, workflowParams)
	if err != nil {
		return fmt.Errorf("error executing Temporal workflow: %v", err)
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

func populateNotification(tenant, webhookName string, wh *entity.TenantWebhook) *notifications.NovuNotification {
	payload := map[string]interface{}{
		"subject":       fmt.Sprintf("Webhook %s is currently offline", webhookName),
		"email":         wh.UserEmail,
		"tenant":        tenant,
		"userFirstName": wh.UserFirstName,
		"webhookUrl":    wh.WebhookUrl,
	}

	notification := &notifications.NovuNotification{
		WorkflowId: notifications.WorkflowFailedWebhook,
		TemplateData: map[string]string{
			"{{userFirstName}}": wh.UserFirstName,
			"{{webhookName}}":   webhookName,
			"{{webhookUrl}}":    wh.WebhookUrl,
		},
		To: &notifications.NotifiableUser{
			FirstName:    wh.UserFirstName,
			LastName:     wh.UserLastName,
			Email:        wh.UserEmail,
			SubscriberID: wh.UserId,
		},
		Subject: fmt.Sprintf("Webhook %s is currently offline", webhookName),
		Payload: payload,
	}
	return notification
}
