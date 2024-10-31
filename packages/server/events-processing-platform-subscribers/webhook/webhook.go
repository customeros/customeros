package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
	temporal_client "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/temporal/workflows"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func DispatchWebhook(ctx context.Context, tenant string, event WebhookEvent, payload *InvoicePayload, postgresRepositories *postgresRepository.Repositories, cfg config.Config) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "DispatchWebhook")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("webhookEvent", event.String()))

	if !cfg.Temporal.RunWorker {
		err := errors.New("temporal worker is not running")
		tracing.TraceErr(span, err)
		return err
	}

	// fetch webhook data from db
	webhookResult := postgresRepositories.TenantWebhookRepository.GetWebhook(ctx, tenant, event.String())
	if webhookResult.Error != nil {
		err := fmt.Errorf("error fetching webhook data: %v", webhookResult.Error)
		tracing.TraceErr(span, err)
		return err
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
		tracing.TraceErr(span, err)
		return fmt.Errorf("(webhook.DispatchWebhook) error marshalling request body: %v", err)
	}
	// Start Temporal Client to queue webhook workflow
	tClient, err := temporal_client.TemporalClient(cfg.Temporal.HostPort, cfg.Temporal.Namespace)
	if err != nil {
		tracing.TraceErr(span, err)
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

	var notification *commonService.NovuNotification
	if cfg.Temporal.NotifyOnFailure {
		notification = populateNotification(tenant, event.String(), wh)
	}
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("(webhook.DispatchWebhook) error marshalling notification obj: %v", err)
	}

	workflowParams := workflows.WHWorkflowParam{
		TargetUrl:                  wh.WebhookUrl,
		RequestBody:                string(requestBodyJSON),
		AuthHeaderName:             wh.AuthHeaderName,
		AuthHeaderValue:            wh.AuthHeaderValue,
		RetryPolicy:                retryPolicy,
		Notification:               string(notificationJSON),
		NotificationProviderApiKey: cfg.Services.Novu.ApiKey,
		NotifyFailure:              cfg.Temporal.NotifyOnFailure,
		NotifyAfterAttempts:        cfg.Temporal.NotifyAfterAttempts,
	}

	// the workflow will run async, so we don't need to wait for it to finish
	_, err = tClient.ExecuteWorkflow(context.Background(), workflowOptions, workflows.WebhookWorkflow, workflowParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("error executing Temporal workflow: %v", err)
	}

	return nil
}

func mapResultToWebhook(result helper.QueryResult) *postgresEntity.TenantWebhook {
	if result.Error != nil {
		return nil
	}
	webhook, ok := result.Result.(*postgresEntity.TenantWebhook)
	if !ok {
		return nil
	}
	return webhook
}

func populateNotification(tenant, webhookName string, wh *postgresEntity.TenantWebhook) *commonService.NovuNotification {
	subject := fmt.Sprintf(commonService.WorkflowFailedWebhookSubject, webhookName)
	payload := map[string]interface{}{
		"subject":       subject,
		"email":         wh.UserEmail,
		"tenant":        tenant,
		"userFirstName": wh.UserFirstName,
		"webhookUrl":    wh.WebhookUrl,
	}

	notification := &commonService.NovuNotification{
		WorkflowId: commonService.WorkflowFailedWebhook,
		TemplateData: map[string]string{
			"{{userFirstName}}": wh.UserFirstName,
			"{{webhookName}}":   webhookName,
			"{{webhookUrl}}":    wh.WebhookUrl,
		},
		To: &commonService.NotifiableUser{
			FirstName:    wh.UserFirstName,
			LastName:     wh.UserLastName,
			Email:        wh.UserEmail,
			SubscriberID: wh.UserId,
		},
		Subject: subject,
		Payload: payload,
	}
	return notification
}
