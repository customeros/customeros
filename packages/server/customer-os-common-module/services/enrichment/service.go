package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type enrichmentService struct {
	config *config.GlobalConfig
}

func NewEnrichmentService(config *config.GlobalConfig) interfaces.EnrichmentService {
	return &enrichmentService{
		config: config,
	}
}

const (
	DEFAULT_API_TIMEOUT = 30 // seconds
)

func (s *enrichmentService) Snitcher(ctx context.Context, ipAddress string) (*interfaces.SnitcherDataResponse, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "EnrichmentService.Snitcher")
	tracing.SetDefaultServiceSpanTags(ctx, span)
	defer span.Finish()

	headers := map[string]string{
		"Content-Type":       "application/json",
		"X-OPENLINE-API-KEY": s.config.InternalServices.EnrichmentApiConfig.ApiKey,
	}

	resp, err := s.makeGetRequest(ctx, s.config.InternalServices.EnrichmentApiConfig.Url+"/snitcher?ipAddress="+ipAddress, headers)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Parse response
	var data interfaces.SnitcherDataResponse
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
