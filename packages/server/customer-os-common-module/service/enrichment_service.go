package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type EnrichmentService interface {
	Snitcher(ctx context.Context, ipAddress string) (*SnitcherDataResponse, error)
}

type enrichmentService struct {
	services *Services
	config   *config.GlobalConfig
}

const (
	DEFAULT_API_TIMEOUT = 30 // seconds
)

func NewEnrichmentService(services *Services, config *config.GlobalConfig) EnrichmentService {
	return &enrichmentService{
		services: services,
		config:   config,
	}
}

func (s *enrichmentService) Snitcher(ctx context.Context, ipAddress string) (*SnitcherDataResponse, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "EnrichmentService.Snitcher")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	headers := map[string]string{
		"Content-Type":       "application/json",
		"X-OPENLINE-API-KEY": s.config.InternalServices.EnrichmentApiConfig.ApiKey,
	}

	resp, err := s.makeGetRequest(ctx, s.config.InternalServices.EnrichmentApiConfig.Url, headers)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Parse response
	var data SnitcherDataResponse
	if err := json.Unmarshal(resp, &data); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return &data, nil
}

func (s *enrichmentService) makeGetRequest(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create client and execute request
	client := &http.Client{
		Timeout: DEFAULT_API_TIMEOUT * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
