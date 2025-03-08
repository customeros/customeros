package crust_data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// CrustDataFilter represents a single filter in the Crust Data API request
type CrustDataFilter struct {
	FilterType string   `json:"filter_type"`
	Type       string   `json:"type"`
	Value      []string `json:"value"`
}

// CrustDataRequest represents the request structure for the Crust Data API
type CrustDataRequest struct {
	Filters []CrustDataFilter `json:"filters"`
	Page    int               `json:"page"`
}

// ErrorResponse represents the error response from the API
type ErrorResponse struct {
	Error string `json:"error"`
}

type crustDataService struct {
	log      logger.Logger
	config   *config.CrustDataConfig
	postgres *postgres_repository.Repositories
}

func NewCrustDataService(log logger.Logger,
	config *config.CrustDataConfig,
	postgres *postgres_repository.Repositories) *crustDataService {
	return &crustDataService{
		log:      log,
		config:   config,
		postgres: postgres,
	}
}

func (s *crustDataService) SearchPeople(ctx context.Context, companyDomain string, jobTitles []string) (*CrustDataResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CrustDataService.SearchPeople")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(span)
	span.LogFields(log.String("companyDomain", companyDomain))
	tracing.LogObjectAsJson(span, "jobTitles", jobTitles)

	if s.config.ApiKey == "" {
		err := errors.New("crust data api key is not set")
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	var profiles []interfaces.CrustDataProfile

	// call crust data to get people
	response, err := s.callCrustData(ctx, companyDomain, jobTitle, 1)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	return response, nil
}

func (s *crustDataService) callCrustDataFilterByCompanyAndJobTitle(ctx context.Context, companyDomain string, jobTitle string, page int) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CrustDataService.callCrustDataFilterByCompanyAndJobTitle")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(span)
	span.LogFields(
		log.String("companyDomain", companyDomain),
		log.String("jobTitle", jobTitle),
		log.Int("page", page),
	)

	// Prepare request body
	reqBody := CrustDataRequest{
		Filters: []CrustDataFilter{
			{
				FilterType: "CURRENT_COMPANY",
				Type:       "in",
				Value:      []string{companyDomain},
			},
			{
				FilterType: "CURRENT_TITLE",
				Type:       "in",
				Value:      []string{jobTitle},
			},
		},
		Page: page,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.ApiUrl+"/screener/person/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.config.ApiKey)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle response based on status code
	switch resp.StatusCode {
	case http.StatusOK:
		return string(body), nil
	case http.StatusBadRequest:
		// Parse error response
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			tracing.TraceErr(span, err)
			return "", fmt.Errorf("failed to parse error response: %w", err)
		}

		// Check if it's the specific "no data found" error
		if errorResp.Error == "Failed to retrieve profiles using provided filters" {
			return "", nil // Return empty string and no error for this specific case
		}

		// For other 400 errors, return the error
		span.LogFields(log.String("response.body", string(body)))
		tracing.TraceErr(span, fmt.Errorf("API error: %s", errorResp.Error))
		return "", fmt.Errorf("API error: %s", errorResp.Error)
	default:
		span.SetTag("response.error_code", resp.StatusCode)
		span.LogFields(log.String("response.body", string(body)))
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		tracing.TraceErr(span, err)
		return "", err
	}
}
