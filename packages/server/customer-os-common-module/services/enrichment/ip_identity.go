package enrichment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"io"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

const (
	CACHE_LOOKBACK_DAYS = 30
	MAX_RESPONSE_SIZE   = 1 * 1024 * 1024 // 1MB
)

func (s *enrichmentService) IPIdentity(ctx context.Context, ip string) (*interfaces.SnitcherResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.IPIdentity")
	defer spans.Finish()

	spans.LogKV("ip", ip)

	// check cache if ip exists
	results, err := s.postgres.CacheIPIdentifyRepository.FindByIP(ctx, ip, CACHE_LOOKBACK_DAYS)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to check cache: %w", err)
	}

	if results != nil && results.SnitcherData != "" && results.Domain != "" {
		var snitcherResponse interfaces.SnitcherResponse
		err = json.Unmarshal([]byte(results.SnitcherData), &snitcherResponse)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling snitcher response: %w", err)
		}
		return &snitcherResponse, nil
	}

	// call snitcher
	snitcherResponse, respString, err := s.callSnitcher(ctx, ip)
	if err != nil {
		spans.TraceError(err)
		spans.LogKV("result.rawSnitcherResponse", utils.IfNotNilString(respString))
		return nil, fmt.Errorf("failed to call snitcher: %w", err)
	}

	// Store response
	snitcherData := postgres_entity.CacheIPIdentify{
		IPAddress:    ip,
		Domain:       snitcherResponse.CompanyDomain(),
		LinkedinSlug: snitcherResponse.LinkedinSlug(),
		SnitcherData: *respString,
	}

	err = s.postgres.CacheIPIdentifyRepository.Create(ctx, snitcherData)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to store response: %w", err)
	}

	return snitcherResponse, nil
}

func (s *enrichmentService) callSnitcher(ctx context.Context, ip string) (*interfaces.SnitcherResponse, *string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.callSnitcher")
	defer spans.Finish()
	spans.LogKV("ip", ip)

	// validate if snitcher is configured
	if s.config.SnitcherConfig.ApiKey == "" || s.config.SnitcherConfig.Url == "" {
		err := fmt.Errorf("snitcher is not configured")
		spans.TraceError(err)
		return nil, nil, err
	}
	spans.LogKV("snitcher.url", s.config.SnitcherConfig.Url)
	spans.LogKV("snitcher.apiKey", utils.Mask(s.config.SnitcherConfig.ApiKey))

	// Create POST request with context
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/company/find?ip=%s", s.config.SnitcherConfig.Url, ip), nil)
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.SnitcherConfig.ApiKey)

	// Create HTTP client with timeout
	clientTimeout := 30 * time.Second
	httpClient := clients.NewLoggingClient(s.warehouse.APICallLogRepository, enum.VendorSnitcher, &clientTimeout)

	// Perform the request
	resp, err := httpClient.Do(req)
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
	if err = json.Unmarshal(responseBody, &snitcherResponse); err != nil {
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
