package clients

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

// Define the interface for the repository
func NewLoggingClient(repo repository.APICallLogRepository, vendor enum.APIVendor) *http.Client {
	return &http.Client{
		Transport: &dbLoggingTransport{
			base:   http.DefaultTransport,
			vendor: vendor,
			repo:   repo,
		},
	}
}

type dbLoggingTransport struct {
	base   http.RoundTripper
	vendor enum.APIVendor
	repo   repository.APICallLogRepository
}

func (t *dbLoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	requestID := utils.GenerateNanoIDWithPrefix("api", 21)

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
	logEntry := &models.APICallLog{
		ID:          utils.GenerateNanoIDWithPrefix("api", 21),
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
