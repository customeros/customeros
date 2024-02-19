package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

func DispatchWebhook(tenant string, event WebhookEvent, payload *InvoicePayload, db *repository.Repositories) error {
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

	// Create a POST request with headers and body
	req, err := http.NewRequest("POST", wh.WebhookUrl, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	authHeader := fmt.Sprintf("%s %s", wh.AuthHeader, wh.ApiKey) // e.g. "Bearer <api_key>"
	req.Header.Set("Authorization", authHeader)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

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
