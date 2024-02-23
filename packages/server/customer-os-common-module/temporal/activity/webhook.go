package activity

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

func WebhookActivity(ctx context.Context, targetUrl, authHeaderName, authHeaderValue string, reqBody *bytes.Buffer) error {
	// Create a POST request with headers and body
	req, err := http.NewRequest("POST", targetUrl, reqBody)
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
