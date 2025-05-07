package clients

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

// Define the interface for the repository
func NewLoggingClient(repo postgres_repository.APICallLogRepository, vendor enum.APIVendor, timeout *time.Duration) *http.Client {
	client := &http.Client{
		Transport: &dbLoggingTransport{
			base:   http.DefaultTransport,
			vendor: vendor,
			repo:   repo,
		},
	}

	if timeout != nil {
		client.Timeout = *timeout
	}

	return client
}

type dbLoggingTransport struct {
	base   http.RoundTripper
	vendor enum.APIVendor
	repo   postgres_repository.APICallLogRepository
}

func (t *dbLoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	requestID := utils.GenerateNanoIdWithPrefix("api", 21)

	// Capture request body for logging
	var requestBodyBytes []byte
	if req.Body != nil && req.Header.Get("Content-Type") != "multipart/form-data" {
		var err error
		requestBodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}

		// Restore the body for the actual request
		req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	}

	// Record start time
	start := time.Now()

	// Make the actual request
	resp, err := t.base.RoundTrip(req)

	// Calculate duration
	duration := time.Since(start)
	durationMs := int(duration.Milliseconds())

	// Prepare log entry
	logEntry := &postgres_entity.APICallLog{
		ID:          utils.GenerateNanoIdWithPrefix("api", 21),
		Vendor:      t.vendor,
		Method:      req.Method,
		URL:         req.URL.String(),
		RequestID:   requestID,
		RequestBody: requestBodyBytes,
		Timestamp:   start,
		Duration:    durationMs,
	}

	// Handle response or error
	if err != nil {
		errMsg := err.Error()
		logEntry.ErrorMessage = &errMsg
	} else {
		logEntry.StatusCode = &resp.StatusCode

		// Capture response body for logging
		if resp.Body != nil {
			responseBodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				errMsg := fmt.Sprintf("failed to read response body: %v", err)
				logEntry.ErrorMessage = &errMsg
			} else {
				logEntry.ResponseBody = &responseBodyBytes

				// Restore the body for the caller
				resp.Body = io.NopCloser(bytes.NewBuffer(responseBodyBytes))
			}
		}
	}

	// Log asynchronously to not block the request
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if logErr := t.repo.Create(bgCtx, logEntry); logErr != nil {
			// If logging fails, log to stderr as a fallback
			fmt.Fprintf(os.Stderr, "Failed to log API call: %v\n", logErr)
		}
	}()

	return resp, err
}
