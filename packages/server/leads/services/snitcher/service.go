package snitcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/leads/internal/config"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
)

type SnitcherService struct {
	config       *config.SnitcherConfig
	natsConn     *nats_internal.NATSConnections
	repositories *repository.Repositories
}

func NewSnitcherService(config *config.SnitcherConfig, repos *repository.Repositories, natsConn *nats_internal.NATSConnections) *SnitcherService {
	return &SnitcherService{
		config:       config,
		natsConn:     natsConn,
		repositories: repos,
	}
}

const (
	HTTP_TIMEOUT      = 60 * time.Second
	MAX_RESPONSE_SIZE = 1
)

func (s *SnitcherService) AskSnitcher(ctx context.Context, ip string) (*SnitcherResponse, *string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "snitcherService.AskSnitcher")
	defer spans.Finish()

	spans.LogKV("ip", ip)

	// validate if snitcher is configured
	if s.config.ApiKey == "" || s.config.Url == "" {
		err := fmt.Errorf("snitcher is not configured")
		spans.TraceError(err)
		return nil, nil, err
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: HTTP_TIMEOUT,
	}

	// Create POST request with context
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/company/find?ip=%s", s.config.Url, ip), nil)
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.ApiKey)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("failed to perform POST request: %w", err)
	}
	defer resp.Body.Close()

	// Read response with size limit
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	spans.LogKV("response.statusCode", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		spans.LogKV("result.rawSnitcherResponse", string(responseBody))
		return nil, nil, fmt.Errorf("snitcher API returned non-200 status code: %d", resp.StatusCode)
	}

	// Validate and compact JSON
	validatedAndCompacted, err := validateAndCompactJSON(responseBody)
	if err != nil {
		spans.TraceError(err)
		spans.LogKV("json.response.invalid", string(responseBody))
		return nil, nil, fmt.Errorf("failed to process JSON response: %w", err)
	}

	// Parse the response
	var snitcherResponse SnitcherResponse
	if err := json.Unmarshal(responseBody, &snitcherResponse); err != nil {
		spans.TraceError(err)
		spans.LogKV("json.response.parsing", string(responseBody))
		return nil, nil, fmt.Errorf("failed to parse snitcher response: %w", err)
	}

	spans.LogObjectAsJson("result.snitcherData", snitcherResponse)
	return &snitcherResponse, &validatedAndCompacted, nil
}

func validateAndCompactJSON(data []byte) (string, error) {
	if !json.Valid(data) {
		return "", fmt.Errorf("invalid JSON")
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, data); err != nil {
		return "", fmt.Errorf("failed to compact JSON: %w", err)
	}

	return buf.String(), nil
}
