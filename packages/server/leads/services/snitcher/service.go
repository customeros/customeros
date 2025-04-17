package snitcher

import (
	"context"

	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
)

type SnitcherService interface {
	AskSnitcher(ctx context.Context, ipAddress string) *SnitcherResponse
}

type snitcherService struct {
	config       *config.SnitcherConfig
	repositories *repository.Repositories
}

func NewSnitcherService(config *config.SnitcherConfig, repos *repository.Repositories) SnitcherService {
	return &snitcherService{
		config:       config,
		repositories: repos,
	}
}

func (s *snitcherService) AskSnitcher(ctx context.Context, ip string) (*SnitcherResponse, *string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "snitcherService.AskSnitcher")
	defer spans.Finish()

	spans.LogKV("ip", ip)

	// validate if snitcher is configured
	if s.config.SnitcherConfig.ApiKey == "" || s.config.SnitcherConfig.Url == "" {
		err := fmt.Errorf("snitcher is not configured")
		spans.TraceError(err)
		return nil, nil, err
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: HTTP_TIMEOUT,
	}

	// Create POST request with context
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/company/find?ip=%s", s.config.SnitcherConfig.Url, ip), nil)
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.SnitcherConfig.ApiKey)

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
	var snitcherResponse interfaces.SnitcherResponse
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
