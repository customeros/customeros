package enrichment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"io"
	"net/http"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

const CACHE_LOOKBACK_DAYS = 30

func (s *enrichmentService) IPIdentity(ctx context.Context, ip string) (*interfaces.SnitcherResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.GetSnitcherData")
	defer spans.Finish()

	spans.LogKV("ip", ip)

	// check caching
	results, err := s.postgres.CacheIPIdentifyRepository.FindByIP(ctx, ip, CACHE_LOOKBACK_DAYS)
	if err != nil {
		spans.TraceError(err)
	}

	if results != nil && results.SnitcherData != "" && results.Domain != "" {
		var snitcherResponse interfaces.SnitcherResponse
		err := json.Unmarshal([]byte(results.SnitcherData), &snitcherResponse)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling snitcher response: %w", err)
		}
		return &snitcherResponse, nil
	}

	// call snitcher
	snitcherResponse, respString, err := s.callSnitcher(ctx, ip)
	if err != nil {
		spans.TraceError(err)
		return nil, err
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
		return nil, fmt.Errorf("failed to store response: %v", err)
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

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/company/find?ip=%s", s.config.SnitcherConfig.Url, ip), nil)
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.SnitcherConfig.ApiKey)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check status code
	spans.LogKV("response.statusCode", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		spans.LogKV("result.rawSnitcherResponse", string(responseBody))
		return nil, nil, fmt.Errorf("snitcher API returned non-200 status code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, responseBody); err != nil {
		return nil, nil, fmt.Errorf("error compacting json: %w", err)
	}
	responseString := buf.String()

	snitcherData, err := buildSnitcherResponse(ctx, responseBody)
	if err != nil {
		spans.TraceError(err)
		return nil, nil, fmt.Errorf("faild to parse snitcher response: %v", err)
	}

	spans.LogKV("result.rawSnitcherResponse", string(responseBody))
	spans.LogObjectAsJson("result.snitcherData", snitcherData)
	return snitcherData, &responseString, nil
}

func buildSnitcherResponse(ctx context.Context, responseBody []byte) (*interfaces.SnitcherResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EnrichmentService.buildSnitcherResponse")
	defer spans.Finish()

	// Parse the JSON request body
	var snitcherResponse interfaces.SnitcherResponse
	if err := json.Unmarshal(responseBody, &snitcherResponse); err != nil {
		spans.TraceError(err)
		spans.LogKV("snitcherResponseBody", string(responseBody))
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &snitcherResponse, nil
}
