package enrichment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/pkg/errors"
)

const CacheCrustDataTTL = 90 * 24 * time.Hour

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

func (s *crustDataService) SearchPeople(ctx context.Context, companyDomain string, jobTitles []string) (*interfaces.CrustDataResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CrustDataService.SearchPeople")
	defer spans.Finish()

	spans.LogKV("companyDomain", companyDomain)
	spans.LogObjectAsJson("jobTitles", jobTitles)

	if s.config.ApiKey == "" {
		err := errors.New("crust data api key is not set")
		spans.TraceError(err)
		s.log.Error(err)
		return nil, err
	}

	var profiles []interfaces.CrustDataProfile

	for _, jobTitle := range jobTitles {
		// Check cache if record exists by companyDomain and jobTitle and ttl, call cache crust data repository
		cachedRecords, err := s.postgres.CacheCrustDataRepository.GetByCompanyAndTitle(ctx, companyDomain, jobTitle, CacheCrustDataTTL)
		if err != nil {
			spans.TraceError(err)
			s.log.Error(err)
			return nil, err
		}
		if len(cachedRecords) > 0 {
			for _, cachedRecord := range cachedRecords {
				var crustDataResponse interfaces.CrustDataResponse
				if err := json.Unmarshal([]byte(cachedRecord.Response), &crustDataResponse); err != nil {
					spans.TraceError(err)
					s.log.Error(err)
					return nil, err
				}
				profiles = append(profiles, crustDataResponse.Profiles...)
			}
		} else {
			// call crust data api
			response, err := s.callCrustDataFilterByCompanyAndJobTitle(ctx, companyDomain, jobTitle, 1)
			if err != nil {
				spans.TraceError(err)
				s.log.Error(err)
				return nil, err
			}
			if response != "" {
				var crustDataResponse interfaces.CrustDataResponse
				if err := json.Unmarshal([]byte(response), &crustDataResponse); err != nil {
					spans.TraceError(err)
					s.log.Error(err)
					return nil, err
				}
				profiles = append(profiles, crustDataResponse.Profiles...)
				// save to cache
				_, err = s.postgres.CacheCrustDataRepository.Create(ctx, postgres_entity.CacheCrustData{
					RequestCompanyDomain: companyDomain,
					RequestJobTitle:      jobTitle,
					Response:             response,
				})
				if err != nil {
					spans.TraceError(err)
					s.log.Error(err)
					return nil, err
				}
			}
		}
	}

	response := interfaces.CrustDataResponse{
		Profiles: profiles,
	}

	return &response, nil
}

func (s *crustDataService) callCrustDataFilterByCompanyAndJobTitle(ctx context.Context, companyDomain string, jobTitle string, page int) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "CrustDataService.callCrustDataFilterByCompanyAndJobTitle")
	defer spans.Finish()

	spans.LogFields(
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
		spans.TraceError(err)
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.ApiUrl+"/screener/person/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		spans.TraceError(err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.config.ApiKey)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		spans.TraceError(err)
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(err)
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
			spans.TraceError(err)
			return "", fmt.Errorf("failed to parse error response: %w", err)
		}

		// Check if it's the specific "no data found" error
		if errorResp.Error == "Failed to retrieve profiles using provided filters" {
			return "", nil // Return empty string and no error for this specific case
		}

		// For other 400 errors, return the error
		spans.LogKV("response.body", string(body))
		spans.TraceError(fmt.Errorf("API error: %s", errorResp.Error))
		return "", fmt.Errorf("API error: %s", errorResp.Error)
	default:
		spans.LogKV("response.error_code", resp.StatusCode)
		spans.LogKV("response.body", string(body))
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		spans.TraceError(err)
		return "", err
	}
}
