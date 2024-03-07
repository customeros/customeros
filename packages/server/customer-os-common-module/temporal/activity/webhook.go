package activity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.temporal.io/sdk/activity"
)

func WebhookActivity(ctx context.Context, targetUrl, authHeaderName, authHeaderValue string, reqBody string, notifyFailure bool, notification string, apiKey string, retryLimit int32) error {
	// after 7 attempts, this wraps up notifying user
	defer func() {
		if notifyFailure && activity.GetInfo(ctx).Attempt >= retryLimit {
			err := NotifyUserActivity(notification, apiKey)
			if err != nil {
				fmt.Println("Error notifying user:", err)
			}
		}
	}()
	var data map[string]interface{}

	// Unmarshal the JSON string into a map
	if err := json.Unmarshal([]byte(reqBody), &data); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	// Marshal the map back to JSON
	requestBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// Create a POST request with headers and body
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authHeaderName, authHeaderValue)
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
