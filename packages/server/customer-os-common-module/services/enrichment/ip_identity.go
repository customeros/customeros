package enrichment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

const CACHE_LOOKBACK_DAYS = 30

func (s *enrichmentService) IPIdentity(ctx context.Context, ip string) (*interfaces.SnitcherResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichmentService.GetSnitcherData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("ip", ip)

	// check caching
	results, err := s.postgres.CacheIPIdentifyRepository.FindByIP(ctx, ip, CACHE_LOOKBACK_DAYS)
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store response: %v", err)
	}

	return snitcherResponse, nil
}

func (s *enrichmentService) callSnitcher(ctx context.Context, ip string) (*interfaces.SnitcherResponse, *string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichmentService.callSnitcher")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("ip", ip)

	// validate if snitcher is configured
	if s.config.SnitcherConfig.ApiKey == "" || s.config.SnitcherConfig.Url == "" {
		err := fmt.Errorf("snitcher is not configured")
		tracing.TraceErr(span, err)
		return nil, nil, err
	}

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/company/find?ip=%s", s.config.SnitcherConfig.Url, ip), nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.SnitcherConfig.ApiKey)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check status code
	span.LogFields(log.Int("response.statusCode", resp.StatusCode))
	if resp.StatusCode != http.StatusOK {
		span.LogKV("result.rawSnitcherResponse", string(responseBody))
		return nil, nil, fmt.Errorf("snitcher API returned non-200 status code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, responseBody); err != nil {
		return nil, nil, fmt.Errorf("error compacting json: %w", err)
	}
	responseString := buf.String()

	snitcherData, err := buildSnitcherResponse(ctx, responseBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, fmt.Errorf("faild to parse snitcher response: %v", err)
	}

	span.LogKV("result.rawSnitcherResponse", string(responseBody))
	tracing.LogObjectAsJson(span, "result.snitcherData", snitcherData)
	return snitcherData, &responseString, nil
}

func buildSnitcherResponse(ctx context.Context, responseBody []byte) (*interfaces.SnitcherResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichmentService.buildSnitcherResponse")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	// Parse the JSON request body
	var snitcherResponse interfaces.SnitcherResponse
	if err := json.Unmarshal(responseBody, &snitcherResponse); err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.String("snitcherResponseBody", string(responseBody)))
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &snitcherResponse, nil
}
